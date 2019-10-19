package services

// Service represents a service that holds information about groups and
// group memberships.
type Service interface {
	GroupMembers(group string) []string
}
