package services

// Tools for parsing mappings.

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/google/logger"
	"github.com/logrusorgru/aurora"
)

//A Mapping is a single Mapping of source group(s) onto a target group.
type Mapping struct {
	src  []GroupIdent
	tar  GroupIdent
	diff *DiffResult
}

func NewMapping(src []GroupIdent, tar GroupIdent) Mapping {
	return Mapping{
		src: src,
		tar: tar,
	}
}

// TODO: package the two []User values into a DiffResult object
// Diff() should probably both return the DiffResult for inspection
// AND store it inside the mapping for later use with commit changes!

func (m *Mapping) Diff() (DiffResult, error) {
	if m.diff != nil {
		return *m.diff, nil
	}

	var flattenedSrc []User

	for _, src := range m.src {
		for _, user := range src.group {
			flattenedSrc = append(flattenedSrc, user)
		}
	}

	diff, err := Diff(flattenedSrc, m.tar.group, m.tar.svc)
	// Yes, the equality is intended here! Only cache the DiffResult
	// if there was no error calculating it.
	if err == nil {
		m.diff = &diff
	}

	return diff, err
}

func (m *Mapping) CommitChanges() error {
	diff, err := m.Diff()
	if err != nil {
		return err
	}

	svc, err := TargetFromString(m.tar.svc)
	if err != nil {
		return err
	}

	err = svc.AddMembers(m.tar.name, diff.Add)
	if err != nil {
		return err
	}

	err = svc.RemoveMembers(m.tar.name, diff.Rem)
	if err != nil {
		return err
	}

	return nil
}

func (m Mapping) String() string {
	var b bytes.Buffer

	b.WriteString("Sources:\n")
	for _, src := range m.src {
		b.WriteString(
			fmt.Sprintf(
				"- %s:%s\n",
				aurora.Cyan(src.svc),
				aurora.Blue(src.name),
			),
		)
	}

	b.WriteString("Target:\n")
	b.WriteString(
		fmt.Sprintf(
			"- %s:%s\n",
			aurora.Cyan(m.tar.svc),
			aurora.Blue(m.tar.name),
		),
	)

	if m.diff != nil {
		b.WriteString("Rem:\n")
		for _, u := range m.diff.Rem {
			b.WriteString(
				fmt.Sprintf("- %v\n", u.String()),
			)
		}

		b.WriteString("Add:\n")
		for _, u := range m.diff.Add {
			b.WriteString(
				fmt.Sprintf("- %v\n", u.String()),
			)
		}
	}

	return b.String()
}

type GroupIdent struct {
	name  string
	group []User
	svc   string
}

func (y YAMLGroupIdent) intoGroupIdent() (GroupIdent, error) {
	return ParseGroupIdent(
		fmt.Sprintf("%s:%s", y.Service, y.Group),
	)
}

func ParseGroupIdent(str string) (GroupIdent, error) {
	splitStr := strings.SplitN(str, ":", 2)
	if len(splitStr) != 2 {
		return GroupIdent{}, fmt.Errorf(
			"string `%s` should follow the `service:group` format",
			str,
		)
	}

	svc, err := SvcFromString(splitStr[0])
	if err != nil {
		return GroupIdent{}, err
	}

	grp, err := svc.GroupMembers(splitStr[1])
	if err != nil {
		return GroupIdent{}, err
	}

	result := GroupIdent{
		name:  splitStr[1],
		group: grp,
		svc:   splitStr[0],
	}

	return result, nil
}

// YAMLMapping is a mapping parsed from a YAML mappings file.
type YAMLMapping struct {
	Sources []YAMLGroupIdent
	Target  YAMLGroupIdent
}

// YAML
type YAMLGroupIdent struct {
	Service string
	Group   string
}

func (y YAMLMapping) IntoMapping() Mapping {
	sources := make([]GroupIdent, 0)

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

	return Mapping{
		src: sources,
		tar: target,
	}
}
