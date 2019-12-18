package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/google/logger"
	"github.com/spf13/cobra"

	"github.com/jamf/groupsync/services"
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

var syncCmd = &cobra.Command{
	Use:   "sync <source>... <target>",
	Args:  cobra.MinimumNArgs(0),
	Short: "List the members of a group (or groups)",
	Long:  `List the members of a group (or groups).`,
	Run: func(cmd *cobra.Command, args []string) {
		var mappings []services.Mapping
		var err error

		if MappingFile != "" {
			mappings, err = parseFileMappings(MappingFile)
			if err != nil {
				logger.Fatal(err)
			}
		} else {
			mapping, err := parseCLIMapping(args)
			if err != nil {
				logger.Fatal(err)
			}

			mappings = append(mappings, mapping)
		}

		for _, mapping := range mappings {
			_, err := mapping.Diff()
			if err != nil {
				logger.Fatal(err)
			}

			fmt.Println(mapping.String())

			if DryRun {
				fmt.Println("This is a dry run. No changes committed.")
			} else {
				err = mapping.CommitChanges()
				if err != nil {
					logger.Fatalf("Cannot commit changes! Cause: %s\n", err)
				}
			}
		}
	},
}

func parseFileMappings(filename string) ([]services.Mapping, error) {
	var mappings []services.Mapping

	data, err := ioutil.ReadFile(MappingFile)
	if err != nil {
		return nil, err
	}

	var mappingData []services.YAMLMapping

	yaml.Unmarshal(data, &mappingData)

	for _, mapping := range mappingData {
		mappings = append(mappings, mapping.IntoMapping())
	}

	return mappings, nil
}

func parseCLIMapping(args []string) (services.Mapping, error) {
	if len(args) < 2 {
		return services.Mapping{}, fmt.Errorf(
			"sync requires at least one source and a target",
		)
	}

	var sources []services.GroupIdent
	for _, srcString := range args[:len(args)-1] {
		src, err := services.ParseGroupIdent(srcString)
		if err != nil {
			return services.Mapping{}, err
		}
		sources = append(sources, src)
	}

	target, err := services.ParseGroupIdent(args[len(args)-1])
	if err != nil {
		return services.Mapping{}, err
	}

	return services.NewMapping(sources, target), nil
}
