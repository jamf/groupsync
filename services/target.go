package services

import "fmt"

// Target represents a service whose group memberships can be mutated.
type Target interface {
	AddMembers(team string, users []User) error
	RemoveMembers(team string, users []User) error
	acquireIdentity(user *User) (Identity, error)

	// Target implementors should also implement Service.
	GroupMembers(group string) ([]User, error)
}

func TargetFromString(name string) (Target, error) {
	switch name {
	case "github":
		return githubSvc, nil
	default:
		return nil, fmt.Errorf(
			"no target %s defined",
			name,
		)
	}
}
