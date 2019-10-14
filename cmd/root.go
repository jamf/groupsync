package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "groupsync",
	Short: "Groupsync is a tool for syncing LDAP groups with GitHub teams.",
	Long:  `Groupsync is a tool for syncing LDAP groups with GitHub teams.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
