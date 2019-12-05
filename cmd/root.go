// Package cmd contains all the CLI logic, backed by Cobra. `root.go`
// provides the root command, which is the entry point to all the Cobra stuff.
// Other files define subcommands.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var version = "undefined"

var rootCmd = &cobra.Command{
	Use:   "groupsync",
	Short: "Groupsync is a tool for syncing LDAP groups with GitHub teams.",
	Long:  `Groupsync is a tool for syncing LDAP groups with GitHub teams.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute runs the root CLI command handler, backed by Cobra.
// It parses parameters, flags, etc. and calls subcommands where appropriate.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
