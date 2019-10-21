package services

import (
	"crypto/tls"
	"errors"
	"fmt"

	"gopkg.in/ldap.v3"
)

// LDAP contains the LDAP config and (once established) the active connection
// to an LDAP server.
type LDAP struct {
	conn *ldap.Conn
	cfg  LDAPConfig
}

// LDAPConfig contains all the info needed to connect to (and authenticate with)
// an LDAP instance, as well as how to fetch group membership data from the
// particular scheme used.
type LDAPConfig struct {
	// Connection
	port       int32
	server     string
	ssl        bool
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

func (l *LDAPConfig) connect() *ldap.Conn {
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

	return c
}

func (l *LDAP) connect() {
	l.conn = l.cfg.connect()
}

// GroupMembers returns the members of group `group` as a slice of User
// instances. Implements the Service interface.
func (l LDAP) GroupMembers(group string) ([]User, error) {
	var attrs []string
	for _, attr := range []string{l.cfg.usernameAttribute, l.cfg.emailAttribute} {
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
		l.cfg.userClass,
		l.cfg.searchAttribute,
		group,
		l.cfg.groupBaseDN,
	)

	result, err := l.conn.Search(&ldap.SearchRequest{
		BaseDN:     l.cfg.userBaseDN,
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
		if l.cfg.usernameAttribute != "" {
			member.username = e.GetAttributeValue(l.cfg.usernameAttribute)
			if member.username == "" {
				panic(fmt.Sprintf(
					"Failed to get username (%s) for %s",
					l.cfg.usernameAttribute,
					e.DN,
				))
			}
		}
		if l.cfg.emailAttribute != "" {
			member.email = e.GetAttributeValue(l.cfg.emailAttribute)
			if member.email == "" {
				panic(fmt.Sprintf(
					"Failed to get e-mail (%s) for %s",
					l.cfg.emailAttribute,
					e.DN,
				))
			}
		}

		members = append(
			members,
			User{
				username: e.GetAttributeValue(l.cfg.usernameAttribute),
				email:    e.GetAttributeValue(l.cfg.emailAttribute),
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
