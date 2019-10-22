package services

type GitHub struct {
	cfg GitHubConfig
}

type GitHubConfig struct {
	Token string
}

func NewGitHub() *GitHub {
	return &GitHub{
		cfg: getConfig().GitHub,
	}
}

func (g *GitHub) GroupMembers(group string) ([]User, error) {
	return nil, nil
}
