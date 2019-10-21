package services

import (
	"errors"
)

// User is used to identify users by their unique data acquired from
// services.
type User struct {
	username string
	email    string
}

func (u User) getUsername() (string, error) {
	if u.username == "" {
		return "", errors.New("username missing")
	}
	return u.username, nil
}

func (u User) getEmail() (string, error) {
	if u.email == "" {
		return "", errors.New("e-mail missing")
	}
	return u.email, nil
}
