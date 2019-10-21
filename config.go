package main

import (
	"fmt"

	"github.com/spf13/viper"

	"stash.jamf.build/devops/groupsync/services"
)

type config struct {
	Test string
	LDAP services.LDAPConfig
}

func getConfig() config {
	viper.SetConfigName("groupsync")
	viper.AddConfigPath("/etc/groupsync/")
	viper.AddConfigPath("$HOME/.groupsync/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	var c config
	err = viper.Unmarshal(&c)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshalling config: %s", err))
	}

	return c
}
