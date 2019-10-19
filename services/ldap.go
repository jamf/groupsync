package services

import (
	"crypto/tls"
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
	userBaseDN      string
	groupBaseDN     string
	userClass       string
	searchAttribute string
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

func (l LDAP) members(group string) *ldap.SearchResult {
	if l.conn == nil {
		panic("No LDAP connection!")
	}

	filter := fmt.Sprintf(
		"(&(objectClass=%s)(%s=cn=%s,%s))",
		l.userClass,
		l.searchAttribute,
		group,
		l.groupBaseDN,
	)

	result, err := l.conn.Search(&ldap.SearchRequest{
		BaseDN: l.userBaseDN,
		Filter: filter,
		Scope:  1,
	})
	if err != nil {
		panic(err)
	}

	return result
}

func (l *LDAP) close() {
	if l.conn != nil {
		l.conn.Close()
		l.conn = nil
	}
}
