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
			panic(err)
		}

		for _, grp := range args[1:] {
			members, err := svc.GroupMembers(grp)
			if err != nil {
				logger.Errorf(
					"Error looking up members of group %s!\nError: %s\n",
					grp,
					err,
				)
				continue
			}
			for _, m := range members {
				fmt.Println(m)
			}
		}
	},
}
