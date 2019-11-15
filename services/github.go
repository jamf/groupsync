package services

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type GitHub struct {
	client *githubv4.Client
	cfg    GitHubConfig
}

type GitHubConfig struct {
	Token string
	Org   string
}

type GitHubIdentity struct {
	ID    string
	Login string
}

// Implement Identity for GitHubIdentity
func (i GitHubIdentity) uniqueID() string {
	return i.ID
}

func (i GitHubIdentity) String() string {
	return fmt.Sprintf("github{uid: %s, login: %s}", i.ID, i.Login)
}

// GitHubSAMLMapping represents a mapping of a GitHub identity to a SAML
// identity.
type GitHubSAMLMapping struct {
	User struct {
		ID    string
		Name  string
		Email string
		Login string
	}
	SamlIdentity struct {
		NameID string `graphql:"nameId"`
	} `graphql:"samlIdentity"`
}

func NewGitHub() *GitHub {
	return &GitHub{
		cfg: getConfig().GitHub,
	}
}

// Implement Service for GitHub

func (g *GitHub) GroupMembers(group string) ([]User, error) {
	g.initClient()

	var membersQuery struct {
		Viewer struct {
			Organization struct {
				Team struct {
					Name    string
					Members struct {
						Edges []struct {
							Node GitHubIdentity
						}
					}
				} `graphql:"team(slug: $grp)"`
			} `graphql:"organization(login: $org)"`
		}
	}

	vars := map[string]interface{}{
		"org": githubv4.String(g.cfg.Org),
		"grp": githubv4.String(group),
	}

	err := g.client.Query(
		context.Background(),
		&membersQuery,
		vars,
	)
	if err != nil {
		return nil, err
	}

	if membersQuery.Viewer.Organization.Team.Name == "" {
		return nil, fmt.Errorf("Cannot find GitHub team called \"%s\"", group)
	}

	var result []User

	for _, entry := range membersQuery.Viewer.Organization.Team.Members.Edges {
		user := newUser()
		user.addIdentity("github", entry.Node)
		result = append(result, user)
	}

	return result, nil
}

func (g *GitHub) acquireIdentity(user *User) (Identity, error) {
	_, err := user.getIdentity("ldap")
	if err == nil {
		return NoneIdentity{}, nil
	}

	return nil, fmt.Errorf("couldn't get the GitHub ID for user\n")
}

func (g *GitHub) initClient() {
	if g.client == nil {
		src := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: g.cfg.Token},
		)
		httpClient := oauth2.NewClient(context.Background(), src)

		g.client = githubv4.NewClient(httpClient)
	}
}

// getAllGitHubMappings fetches all the mappings of GitHub identities to SAML
// identities within the given org.
func (g *GitHub) getAllGitHubMappings() (map[string]GitHubSAMLMapping, error) {
	g.initClient()

	result := make(map[string]GitHubSAMLMapping)

	var firstQuery struct {
		Viewer struct {
			Organization struct {
				SamlIdentityProvider struct {
					ExternalIdentities struct {
						Edges []struct {
							Node GitHubSAMLMapping
						}
						PageInfo struct {
							EndCursor   string
							HasNextPage bool
						}
					} `graphql:"externalIdentities(first:20)"`
				}
			} `graphql:"organization(login: $org)"`
		}
	}

	var nextQuery struct {
		Viewer struct {
			Organization struct {
				SamlIdentityProvider struct {
					ExternalIdentities struct {
						Edges []struct {
							Node GitHubSAMLMapping
						}
						PageInfo struct {
							EndCursor   string
							HasNextPage bool
						}
					} `graphql:"externalIdentities(first:20 after:$page_cursor)"`
				}
			} `graphql:"organization(login: $org)"`
		}
	}

	vars := map[string]interface{}{
		"org": githubv4.String(g.cfg.Org),
	}

	err := g.client.Query(
		context.Background(),
		&firstQuery,
		vars,
	)
	if err != nil {
		return nil, err
	}

	for _, e := range firstQuery.Viewer.Organization.
		SamlIdentityProvider.ExternalIdentities.Edges {
		result[e.Node.User.ID] = e.Node
	}

	hasNextPage := firstQuery.Viewer.Organization.SamlIdentityProvider.
		ExternalIdentities.PageInfo.HasNextPage
	cursor := firstQuery.Viewer.Organization.SamlIdentityProvider.
		ExternalIdentities.PageInfo.EndCursor

	for hasNextPage {
		vars = map[string]interface{}{
			"org":         githubv4.String(g.cfg.Org),
			"page_cursor": githubv4.String(cursor),
		}

		err := g.client.Query(
			context.Background(),
			&nextQuery,
			vars,
		)
		if err != nil {
			return nil, err
		}

		for _, e := range nextQuery.Viewer.Organization.
			SamlIdentityProvider.ExternalIdentities.Edges {
			result[e.Node.User.ID] = e.Node
		}

		hasNextPage = nextQuery.Viewer.Organization.SamlIdentityProvider.
			ExternalIdentities.PageInfo.HasNextPage
		cursor = nextQuery.Viewer.Organization.SamlIdentityProvider.
			ExternalIdentities.PageInfo.EndCursor
	}

	if result == nil {
		return nil, fmt.Errorf(
			"no SAML identities found in the GitHub org `%s` at all",
			g.cfg.Org,
		)
	}

	return result, nil
}
