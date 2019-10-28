package services

import (
	"fmt"
	"strings"
)

// User is used to identify users by their unique data acquired from
// services.
type User interface {
	samlUsername() (string, error)
	email() (string, error)
}

func userEqual(u1, u2 User) bool {
	un1, e1 := u1.samlUsername()
	un2, e2 := u2.samlUsername()
	if e1 != nil {
		panic(e1)
	}
	if e2 != nil {
		panic(e2)
	}

	return un1 == un2
}

func SprintUser(u User) string {
	result, err := u.samlUsername()
	if err != nil {
		panic(err)
	}

	mail, err := u.email()
	if err == nil && mail != "" {
		result = fmt.Sprintf("%s (%s)", result, mail)
	}

	return result
}

func SprintUsers(users []User) string {
	var userStrings []string

	for _, u := range users {
		userStrings = append(userStrings, SprintUser(u))
	}

	return strings.Join(userStrings, "\n")
}
