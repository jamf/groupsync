package services

import (
	"fmt"
)

// Service represents a service that holds information about groups and
// group memberships.
type Service interface {
	// Get the members of group `group` as a slice of User instances.
	GroupMembers(group string) ([]User, error)
	acquireIdentity(user *User) (Identity, error)
}

func SvcFromString(name string) (Service, error) {
	switch name {
	case "ldap":
		return NewLDAP(), nil
	case "github":
		return NewGitHub(), nil
	case "mockservice":
		return newMockService(), nil
	default:
		return nil, fmt.Errorf(
			"no service %s defined",
			name,
		)
	}
}
