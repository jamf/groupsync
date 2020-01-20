package services

import (
	"fmt"

	"github.com/spf13/viper"
)

type config struct {
	LDAP            LDAPConfig
	GitHub          GitHubConfig
	AppStoreConnect AppStoreConnectConfig
}

var cfg *config = nil

func initConfig() error {
	if cfg != nil {
		return newConfigError(
			fmt.Errorf("config already initialized"),
		)
	}

	viper.SetConfigName("groupsync")
	viper.SetDefault("LDAP", LDAPConfig{})
	viper.SetDefault("GitHub", GitHubConfig{})
	viper.SetDefault("AppStoreConnect", AppStoreConnectConfig{})
	viper.AddConfigPath("/etc/groupsync/")
	viper.AddConfigPath("$HOME/.groupsync/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil && configFileMustExist {
		return newConfigError(err)
	}

	var c config
	err = viper.Unmarshal(&c)
	if err != nil {
		return newConfigError(err)
	}

	cfg = &c
	return nil
}

func getConfig() (config, error) {
	if cfg == nil {
		err := initConfig()
		if err != nil {
			return config{}, err
		}
	}

	return *cfg, nil
}

type ConfigError struct {
	originalError error
}

func newConfigError(original error) ConfigError {
	return ConfigError{
		originalError: original,
	}
}

func (e ConfigError) Error() string {
	return fmt.Sprintf("Error initializing config: %v", e.originalError)
}
