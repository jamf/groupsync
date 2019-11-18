package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jamf/groupsync/services"
	. "github.com/logrusorgru/aurora"
)

var DryRun bool

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolVarP(
		&DryRun,
		"dry-run",
		"d",
		false,
		"don't commit any changes, only print what would be added/removed",
	)
}

type groupIdent struct {
	name  string
	group []services.User
	svc   string
}
type mapping struct {
	src groupIdent
	tar groupIdent
}

var syncCmd = &cobra.Command{
	Use:   "sync <source> <target>",
	Args:  cobra.MinimumNArgs(2),
	Short: "List the members of a group (or groups)",
	Long:  `List the members of a group (or groups).`,
	Run: func(cmd *cobra.Command, args []string) {
		if !DryRun {
			// Commit changes
			panic("sync is not implemented yet except with the dry-run flag")
		}

		mappings, err := parseMappings(args)
		if err != nil {
			panic(err)
		}

		for _, mapping := range mappings {
			rem, add := services.Diff(
				mapping.src.group,
				mapping.tar.group,
				mapping.tar.svc,
			)
			fmt.Printf(
				"Results for %s:%s -> %s:%s:\n",
				Cyan(mapping.src.svc),
				Blue(mapping.src.name),
				Cyan(mapping.tar.svc),
				Blue(mapping.tar.name),
			)
			fmt.Printf(
				"Rem: %+v\nAdd: %+v\n\n",
				rem,
				add,
			)
		}
	},
}

func parseMappings(args []string) ([]mapping, error) {
	if len(args)%2 != 0 {
		return nil, fmt.Errorf("uneven number of arguments")
	}

	var result []mapping

	for i := 0; i < len(args); i += 2 {
		mapping, err := parseMapping(args[i], args[i+1])
		if err != nil {
			return nil, err
		}

		result = append(result, mapping)
	}

	return result, nil
}

func parseMapping(ident1, ident2 string) (mapping, error) {
	src, err := parseGroupIdent(ident1)
	if err != nil {
		return mapping{}, err
	}
	tar, err := parseGroupIdent(ident2)
	if err != nil {
		return mapping{}, err
	}

	result := mapping{
		src: src,
		tar: tar,
	}

	return result, nil
}

func parseGroupIdent(str string) (groupIdent, error) {
	splitStr := strings.SplitN(str, ":", 2)
	if len(splitStr) != 2 {
		return groupIdent{}, fmt.Errorf(
			"string `%s` should follow the `service:group` format",
			str,
		)
	}

	svc, err := services.SvcFromString(splitStr[0])
	if err != nil {
		return groupIdent{}, err
	}

	grp, err := svc.GroupMembers(splitStr[1])
	if err != nil {
		return groupIdent{}, err
	}

	result := groupIdent{
		name:  splitStr[1],
		group: grp,
		svc:   splitStr[0],
	}

	return result, nil
}
