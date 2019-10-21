package services

import (
	"crypto/tls"
	"errors"
	"fmt"

	"gopkg.in/ldap.v3"
)

// LDAP contains all the info needed to connect to (and authenticate with)
// an LDAP instance, as well as how to fetch group membership data from the
// particular scheme used.
type LDAP struct {
	// Connection
	port       int32
	server     string
	ssl        bool
	conn       *ldap.Conn
	skipVerify bool

	// Auth
	bindUser     string
	bindPassword string

	// Schema
	userBaseDN        string
	groupBaseDN       string
	userClass         string
	searchAttribute   string
	usernameAttribute string
	emailAttribute    string
}

func (l *LDAP) connect() {
	var c *ldap.Conn
	var err error

	if l.ssl {
		c, err = ldap.DialTLS(
			"tcp",
			fmt.Sprintf("%s:%d", l.server, l.port),
			&tls.Config{
				InsecureSkipVerify: l.skipVerify,
			},
		)
	} else {
		c, err = ldap.Dial(
			"tcp",
			fmt.Sprintf("%s:%d", l.server, l.port),
		)
	}

	if err != nil {
		panic(
			fmt.Sprintf("Error when connecting!\n%s", err),
		)
	}

	l.conn = c
}

func (l LDAP) GroupMembers(group string) ([]User, error) {
	var attrs []string
	for _, attr := range []string{l.usernameAttribute, l.emailAttribute} {
		if attr != "" {
			attrs = append(attrs, attr)
		}
	}

	if len(attrs) < 1 {
		return nil,
			errors.New("LDAP config didn't provide any attributes to look up")
	}

	if l.conn == nil {
		return nil,
			errors.New("no LDAP connection")
	}

	filter := fmt.Sprintf(
		"(&(objectClass=%s)(%s=cn=%s,%s))",
		l.userClass,
		l.searchAttribute,
		group,
		l.groupBaseDN,
	)

	result, err := l.conn.Search(&ldap.SearchRequest{
		BaseDN:     l.userBaseDN,
		Filter:     filter,
		Scope:      1,
		Attributes: attrs,
	})
	if err != nil {
		return nil, err
	}

	var members []User

	for _, e := range result.Entries {
		member := User{}
		if l.usernameAttribute != "" {
			member.username = e.GetAttributeValue(l.usernameAttribute)
			if member.username == "" {
				panic(fmt.Sprintf(
					"Failed to get username (%s) for %s",
					l.usernameAttribute,
					e.DN,
				))
			}
		}
		if l.emailAttribute != "" {
			member.email = e.GetAttributeValue(l.emailAttribute)
			if member.email == "" {
				panic(fmt.Sprintf(
					"Failed to get e-mail (%s) for %s",
					l.emailAttribute,
					e.DN,
				))
			}
		}

		members = append(
			members,
			User{
				username: e.GetAttributeValue(l.usernameAttribute),
				email:    e.GetAttributeValue(l.emailAttribute),
			},
		)
	}

	return members, nil
}

func (l *LDAP) close() {
	if l.conn != nil {
		l.conn.Close()
		l.conn = nil
	}
}
