package services

import (
	"bytes"
	"fmt"
)

// User is used to identify users by their unique data acquired from
// services.
type User struct {
	identities map[string]Identity
}

func (u User) String() string {
	buf := bytes.Buffer{}
	for _, id := range u.identities {
		buf.WriteString(fmt.Sprintf("%s ", id))
	}
	return buf.String()
}

func newUser() User {
	return User{identities: make(map[string]Identity)}
}

func (u *User) addIdentity(svc string, i Identity) {
	u.identities[svc] = i
}

type Identity interface {
	uniqueID() string
	String() string
}

type NoneIdentity struct{}

func (_ NoneIdentity) uniqueID() string {
	panic("identity doesn't exist")
}

func (_ NoneIdentity) String() string {
	return ""
}

func IdentityExists(i Identity) bool {
	_, ok := i.(NoneIdentity)

	return !ok
}

func (u *User) getIdentity(svc_name string) (Identity, error) {
	// Check if the identity is already stored in this instance of User
	id, ok := u.identities[svc_name]
	if ok {
		return id, nil
	}

	// Attempt to acquire the identity
	svc, err := SvcFromString(svc_name)
	if err != nil {
		return nil, err
	}

	id, err = svc.getSvcIdentity(u.identities)
	if err != nil {
		return nil, err
	}

	// Both store the identity and return it
	u.identities[svc_name] = id
	return id, nil
}
