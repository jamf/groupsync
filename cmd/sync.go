package cmd

import (
	"bytes"
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
	src []groupIdent
	tar groupIdent
}

func newMapping(src []groupIdent, tar groupIdent) mapping {
	return mapping{
		src: src,
		tar: tar,
	}
}

func (m mapping) diff() ([]services.User, []services.User) {
	var flattenedSrc []services.User

	for _, src := range m.src {
		for _, user := range src.group {
			flattenedSrc = append(flattenedSrc, user)
		}
	}

	return services.Diff(flattenedSrc, m.tar.group, m.tar.svc)
}

var syncCmd = &cobra.Command{
	Use:   "sync <source>... <target>",
	Args:  cobra.MinimumNArgs(2),
	Short: "List the members of a group (or groups)",
	Long:  `List the members of a group (or groups).`,
	Run: func(cmd *cobra.Command, args []string) {
		if !DryRun {
			// Commit changes
			panic("sync is not implemented yet except with the dry-run flag")
		}

		mapping, err := parseCLIMapping(args)
		if err != nil {
			panic(err)
		}

		rem, add := mapping.diff()

		var b bytes.Buffer

		b.WriteString("Sources:\n")
		for _, src := range mapping.src {
			b.WriteString(
				fmt.Sprintf(
					"- %s:%s\n",
					Cyan(src.svc),
					Blue(src.name),
				),
			)
		}

		b.WriteString("Target:\n")
		b.WriteString(
			fmt.Sprintf(
				"- %s:%s\n",
				Cyan(mapping.tar.svc),
				Blue(mapping.tar.name),
			),
		)

		b.WriteString("Rem:\n")
		for _, u := range rem {
			b.WriteString(
				fmt.Sprintf("- %v\n", u.String()),
			)
		}

		b.WriteString("Add:\n")
		for _, u := range add {
			b.WriteString(
				fmt.Sprintf("- %v\n", u.String()),
			)
		}

		fmt.Println(b.String())
	},
}

func parseCLIMapping(args []string) (mapping, error) {
	var sources []groupIdent
	for _, srcString := range args[:len(args)-1] {
		src, err := parseGroupIdent(srcString)
		if err != nil {
			return mapping{}, err
		}
		sources = append(sources, src)
	}

	target, err := parseGroupIdent(args[len(args)-1])
	if err != nil {
		return mapping{}, err
	}

	return newMapping(sources, target), nil
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
