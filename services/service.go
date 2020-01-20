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

var initializedServices map[string]Service = make(map[string]Service)

// SvcFromString produces a Service object with config taken from the global
// cfg variable.
func SvcFromString(name string) (Service, error) {
	// attempt to lookup an already initialized service
	svc, ok := lookUpServiceInCache(name)
	if ok {
		return svc, nil
	}

	logger.Infof("Service %v not in cache; initializing.", name)
	// if service wasn't found in the cache, attempt to initialize it
	svc, err := newSvcFromName(name)
	if err != nil {
		return nil, err
	}

	saveSvcInCache(name, svc)

	return svc, nil
}

func saveSvcInCache(name string, svc Service) {
	initializedServices[name] = svc
}

func lookUpServiceInCache(name string) (svc Service, ok bool) {
	svc, ok = initializedServices[name]
	return
}

func newSvcFromName(name string) (Service, error) {
	cfg, err := getConfig()
	if err != nil {
		return nil, err
	}

	switch name {
	case "ldap":
		return NewLDAP(cfg.LDAP), nil
	case "github":
		return NewGitHub(cfg.GitHub), nil
	case "appstoreconnect":
		return NewAppStoreConnect(cfg.AppStoreConnect), nil
	case "mockservice":
		return newMockService(), nil
	default:
		return nil, newServiceNotDefined(name)
	}
}

type ServiceNotDefined struct {
	serviceName string
}

func newServiceNotDefined(serviceName string) ServiceNotDefined {
	return ServiceNotDefined{
		serviceName: serviceName,
	}
}

func (e ServiceNotDefined) Error() string {
	return fmt.Sprintf("service `%s` not defined", e.serviceName)
}

type DiffResult struct {
	Rem []User
	Add []User
}

func newDiffResult(rem, add []User) DiffResult {
	return DiffResult{
		Rem: rem,
		Add: add,
	}
}

func Diff(srcGrp, tarGrp []User, tar string) (DiffResult, error) {
	// Build hashmaps of identities for faster lookup.
	// This approach also takes care of duplicates for free.
	srcMap := make(map[string]User)
	tarMap := make(map[string]User)

	if len(srcGrp) < 1 {
		return DiffResult{}, newSourceGroupEmptyError()
	}

	for _, u := range srcGrp {
		i, e := u.getIdentity(tar)
		if e != nil {
			switch e.(type) {
			case FatalError:
				return DiffResult{}, e
			default:
				logger.Warningf(
					"error acquiring identity for a user - skipping\n"+
						"user: %v\nerror: %v\n",
					u,
					e,
				)
			}

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

	var add []User
	var rem []User

	// What's left in srcIdentities and tarIdentities is what we have to
	// add/remove.
	for _, identity := range srcMap {
		add = append(add, identity)
	}

	for _, identity := range tarMap {
		rem = append(rem, identity)
	}

	return newDiffResult(rem, add), nil
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

// A wrapper used to let upper layers know the error isn't recoverable.
type FatalError struct {
	source  error
	context string
}

func newFatalError(context string, source error) FatalError {
	return FatalError{
		source:  source,
		context: context,
	}
}

func (e FatalError) Error() string {
	if e.context == "" {
		return e.source.Error()
	}

	return fmt.Sprintf(
		"error when %v: %v",
		e.context,
		e.source,
	)
}
