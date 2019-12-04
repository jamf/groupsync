package services

func Init() error {
	err := initConfig()
	if err != nil {
		return err
	}

	ldapSvc = NewLDAP()
	githubSvc = NewGitHub()

	return nil
}
