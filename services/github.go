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
	User struct {
		Name  string
		Email string
		Login string
	}
	SamlIdentity struct {
		NameID string `graphql:"nameId"`
	} `graphql:"samlIdentity"`
}

// Implement Identity for GitHubIdentity

func (i GitHubIdentity) uniqueID() string {
	return i.User.Login
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
							Node struct {
								ID   string
								Name string
							}
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

	fmt.Println(membersQuery.Viewer.Organization.Team)

	return nil, nil
}

func (g *GitHub) getSvcIdentity(identities map[string]Identity) (Identity, error) {
	_, ok := identities["ldap"]
	if ok {
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

func (g *GitHub) getAllGitHubIdentities() ([]GitHubIdentity, error) {
	g.initClient()

	var samlQuery struct {
		Viewer struct {
			Organization struct {
				SamlIdentityProvider struct {
					ExternalIdentities struct {
						Edges []struct {
							Node GitHubIdentity
						}
					} `graphql:"externalIdentities(first:20 after:null)"`
				}
			} `graphql:"organization(login: $org)"`
		}
	}

	vars := map[string]interface{}{
		"org": githubv4.String(g.cfg.Org),
	}

	err := g.client.Query(
		context.Background(),
		&samlQuery,
		vars,
	)
	if err != nil {
		return nil, err
	}

	var result []GitHubIdentity

	for _, e := range samlQuery.Viewer.Organization.
		SamlIdentityProvider.ExternalIdentities.Edges {
		result = append(result, e.Node)
	}

	if result == nil {
		return nil, fmt.Errorf(
			"no SAML identities found in the GitHub org `%s` at all",
			g.cfg.Org,
		)
	}

	return result, nil
}
