package services

func Init() error {
	err := initConfig()
	if err != nil {
		return err
	}

	return nil
}
