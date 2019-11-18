package services

// Target represents a service whose group memberships can be mutated.
type Target interface {
	AddMembers(users []User) error
	RemoveMembers(users []User) error

	// Target implementors should also implement Service.
	GroupMembers(group string) ([]User, error)
	acquireIdentity(user *User) (Identity, error)
}
