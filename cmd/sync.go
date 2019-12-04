package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/google/logger"
	"github.com/spf13/cobra"

	"github.com/jamf/groupsync/services"
	. "github.com/logrusorgru/aurora"
	"gopkg.in/yaml.v3"
)

var DryRun bool
var MappingFile string

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolVarP(
		&DryRun,
		"dry-run",
		"d",
		false,
		"don't commit any changes, only print what would be added/removed",
	)
	syncCmd.Flags().StringVarP(
		&MappingFile,
		"mapping-file",
		"m",
		"",
		"the file to use for sync mappings",
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

type YAMLMapping struct {
	Sources []YAMLGroupIdent
	Target  YAMLGroupIdent
}

func (y YAMLMapping) intoMapping() mapping {
	sources := make([]groupIdent, 0)

	for _, yamlSrc := range y.Sources {
		src, err := yamlSrc.intoGroupIdent()
		if err != nil {
			logger.Fatal(err)
		}

		sources = append(sources, src)
	}

	target, err := y.Target.intoGroupIdent()
	if err != nil {
		logger.Fatal(err)
	}

	return mapping{
		src: sources,
		tar: target,
	}
}

type YAMLGroupIdent struct {
	Service string
	Group   string
}

func (y YAMLGroupIdent) intoGroupIdent() (groupIdent, error) {
	return parseGroupIdent(
		fmt.Sprintf("%s:%s", y.Service, y.Group),
	)
}

func newMapping(src []groupIdent, tar groupIdent) mapping {
	return mapping{
		src: src,
		tar: tar,
	}
}

func (m mapping) diff() ([]services.User, []services.User, error) {
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
	Args:  cobra.MinimumNArgs(0),
	Short: "List the members of a group (or groups)",
	Long:  `List the members of a group (or groups).`,
	Run: func(cmd *cobra.Command, args []string) {
		var mappings []mapping

		if MappingFile != "" {
			data, err := ioutil.ReadFile(MappingFile)
			if err != nil {
				logger.Fatal(err)
			}

			var mappingData []YAMLMapping

			yaml.Unmarshal(data, &mappingData)

			for _, mapping := range mappingData {
				mappings = append(mappings, mapping.intoMapping())
			}

		} else {
			mapping, err := parseCLIMapping(args)
			if err != nil {
				logger.Fatal(err)
			}

			mappings = append(mappings, mapping)
		}

		for _, mapping := range mappings {
			rem, add, err := mapping.diff()
			if err != nil {
				logger.Fatal(err)
			}

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

			if DryRun {
				fmt.Println("This is a dry run. No changes committed.")
			} else {
				err = commitChanges(mapping.tar.svc, mapping.tar.name, add, rem)
				if err != nil {
					logger.Fatalf("Cannot commit changes! Cause: %s\n", err)
				}
			}
		}
	},
}

func commitChanges(tar, team string, add, rem []services.User) error {
	svc, err := services.TargetFromString(tar)
	if err != nil {
		return err
	}

	err = svc.AddMembers(team, add)
	if err != nil {
		return err
	}

	err = svc.RemoveMembers(team, rem)
	if err != nil {
		return err
	}

	return nil
}

func parseCLIMapping(args []string) (mapping, error) {
	if len(args) < 2 {
		return mapping{}, fmt.Errorf(
			"sync requires at least one source and a target",
		)
	}

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
