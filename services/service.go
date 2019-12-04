package services

import (
	"fmt"

	"github.com/google/logger"
)

// Service represents a service that holds information about groups and
// group memberships.
type Service interface {
	// Get the members of group `group` as a slice of User instances.
	GroupMembers(group string) ([]User, error)
}

var ldapSvc *LDAP
var githubSvc *GitHub

func SvcFromString(name string) (Service, error) {
	switch name {
	case "ldap":
		return ldapSvc, nil
	case "github":
		return githubSvc, nil
	case "mockservice":
		return newMockService(), nil
	default:
		return nil, fmt.Errorf(
			"no service %s defined",
			name,
		)
	}
}

func Diff(srcGrp, tarGrp []User, tar string) (rem, add []User, err error) {
	// Build hashmaps of identities for faster lookup.
	// This approach also takes care of duplicates for free.
	srcMap := make(map[string]User)
	tarMap := make(map[string]User)

	if len(srcGrp) < 1 {
		return nil, nil, newSourceGroupEmptyError()
	}

	for _, u := range srcGrp {
		i, e := u.getIdentity(tar)
		if e != nil {
			logger.Warningf(
				"error acquiring identity for a user - skipping\n"+
					"user: %v\nerror: %v\n",
				u,
				e,
			)
		} else if IdentityExists(i) {
			srcMap[i.uniqueID()] = u
		}
	}

	for _, u := range tarGrp {
		i, e := u.getIdentity(tar)
		if e != nil {
			logger.Warningf(
				"error acquiring identity for a user - skipping\n"+
					"user: %v\nerror: %v\n",
				u,
				e,
			)
		} else if IdentityExists(i) {
			tarMap[i.uniqueID()] = u
		}
	}

	// Remove elements that exist in both the source and the target.
	for id := range srcMap {
		_, ok := tarMap[id]
		if ok {
			delete(srcMap, id)
			delete(tarMap, id)
		}
	}

	// What's left in srcIdentities and tarIdentities is what we have to
	// add/remove.
	for _, identity := range srcMap {
		add = append(add, identity)
	}

	for _, identity := range tarMap {
		rem = append(rem, identity)
	}

	return
}

type SourceGroupEmptyError struct {
}

func newSourceGroupEmptyError() SourceGroupEmptyError {
	return SourceGroupEmptyError{}
}

func (e SourceGroupEmptyError) Error() string {
	return "sanity check failed: the source group is empty"
}

func createIDMap(ids []Identity) map[string]Identity {
	result := make(map[string]Identity)

	for _, id := range ids {
		result[id.uniqueID()] = id
	}

	return result
}
