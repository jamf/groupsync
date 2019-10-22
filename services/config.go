package services

import (
	"fmt"

	"github.com/spf13/viper"
)

type config struct {
	LDAP   LDAPConfig
	GitHub GitHubConfig
}

var cfg *config = nil

func getConfig() config {
	if cfg == nil {
		viper.SetConfigName("groupsync")
		viper.AddConfigPath("/etc/groupsync/")
		viper.AddConfigPath("$HOME/.groupsync/")
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("fatal error reading config file: %s", err))
		}

		var c config
		err = viper.Unmarshal(&c)
		if err != nil {
			panic(fmt.Errorf("fatal error unmarshalling config: %s", err))
		}

		cfg = &c
		return c
	}
	return *cfg
}
