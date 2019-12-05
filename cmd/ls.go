package cmd

import (
	"fmt"

	"github.com/google/logger"
	"github.com/spf13/cobra"

	"github.com/jamf/groupsync/services"
)

func init() {
	rootCmd.AddCommand(lsCmd)
}

var lsCmd = &cobra.Command{
	Use:   "ls <service> <group>...",
	Args:  cobra.MinimumNArgs(2),
	Short: "List the members of a group (or groups)",
	Long:  `List the members of a group (or groups).`,
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := services.SvcFromString(args[0])
		if err != nil {
			logger.Fatal(err)
		}

		for i, grp := range args[1:] {
			members, err := svc.GroupMembers(grp)
			if err != nil {
				msg := fmt.Sprintf(
					"Error looking up members of group %s!\nError: %s\n",
					grp,
					err,
				)

				// If we're only looking up one group, exit and return a proper
				// exit code for its success/failure.
				if len(args[1:]) == 1 {
					logger.Fatal(msg)
				}

				logger.Error(msg)
				continue
			}

			fmt.Printf("- Group `%s`\n", grp)
			for _, m := range members {
				fmt.Println(m)
			}
			if i < len(args[1:])-1 {
				fmt.Println("")
			}
		}
	},
}
