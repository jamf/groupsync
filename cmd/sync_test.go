package cmd

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/jamf/groupsync/services"
)

var fileMappings = `
- sources:
  - service: ldap
    group: my-group1
  target:
    service: github
    group: my-target-team

- sources:
  - service: ldap
    group: my-group1
  - service: ldap
    group: my-group2
  - service: github
    group: my-source-team
  target:
    service: github
    group: my-target-team
`

func TestParseFileMappings(t *testing.T) {
	file, err := ioutil.TempFile("", "groupsync-mappings-*.yaml")
	if err != nil {
		panic(err)
	}

	_, err = file.WriteString(fileMappings)
	if err != nil {
		panic(err)
	}

	file.Sync()
	file.Close()

	t.Log("trying to open: " + file.Name())

	mappings, err := parseFileMappings(file.Name())
	if err != nil {
		panic(err)
	}

	if len(mappings) != 2 {
		panic("should have found 2 Mapping objects in the test .yaml data")
	}

	ident1, err := services.ParseGroupIdent("ldap:my-group1")
	if err != nil {
		panic(err)
	}

	ident2, err := services.ParseGroupIdent("ldap:my-group2")
	if err != nil {
		panic(err)
	}

	ident3, err := services.ParseGroupIdent("github:my-source-team")
	if err != nil {
		panic(err)
	}

	ident4, err := services.ParseGroupIdent("github:my-target-team")
	if err != nil {
		panic(err)
	}

	var expectedMappings = []services.Mapping{
		services.NewMapping(
			[]services.GroupIdent{
				ident1,
			},
			ident4,
		),
		services.NewMapping(
			[]services.GroupIdent{
				ident1,
				ident2,
				ident3,
			},
			ident4,
		),
	}

	if !reflect.DeepEqual(mappings, expectedMappings) {
		panic(fmt.Sprintf(
			"parsed mappings different than expected\nparsed: %v\nexpected: %v\n",
			mappings,
			expectedMappings,
		))
	}
}

func TestParseCLIMapping(t *testing.T) {
	var err error
	var mapping services.Mapping

	_, err = parseCLIMapping([]string{
		"ldap:bleh",
	})
	if err == nil {
		panic("should have raised an error")
	}

	_, err = parseCLIMapping([]string{})
	if err == nil {
		panic("should have raised an error")
	}

	mapping, err = parseCLIMapping([]string{
		"foo:bar",
		"hoo:baz",
	})
	if err != nil {
		panic(err)
	}

	ok := reflect.DeepEqual(
		mapping,
		mappingConstructor(
			[]string{
				"foo:bar",
			},
			"hoo:baz",
		),
	)

	if !ok {
		panic("parsed mapping not as expected")
	}

	mapping, err = parseCLIMapping([]string{
		"foo:bar",
		"hoo:baz",
		"loo:laz",
	})
	if err != nil {
		panic(err)
	}

	ok = reflect.DeepEqual(
		mapping,
		mappingConstructor(
			[]string{
				"foo:bar",
				"hoo:baz",
			},
			"loo:laz",
		),
	)

	if !ok {
		panic("parsed mapping not as expected")
	}
}

func mappingConstructor(src []string, tar string) services.Mapping {
	var sources []services.GroupIdent

	for _, str := range src {
		s, err := services.ParseGroupIdent(str)
		if err != nil {
			panic(err)
		}

		sources = append(sources, s)
	}

	t, err := services.ParseGroupIdent(tar)
	if err != nil {
		panic(err)
	}

	return services.NewMapping(sources, t)
}
