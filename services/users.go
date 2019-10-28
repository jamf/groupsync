package services

// User is used to identify users by their unique data acquired from
// services.

type User interface {
	samlUsername() (string, error)
	email() (string, error)
}
