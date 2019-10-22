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

func NewGitHub() *GitHub {
	return &GitHub{
		cfg: getConfig().GitHub,
	}
}

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

func (g *GitHub) initClient() {
	if g.client == nil {
		src := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: g.cfg.Token},
		)
		httpClient := oauth2.NewClient(context.Background(), src)

		g.client = githubv4.NewClient(httpClient)
	}
}
