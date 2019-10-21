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
	Port       int32 `mapstructure:"port"`
	Server     string
	SSL        bool
	SkipVerify bool

	// Auth
	BindUser     string
	BindPassword string

	// Schema
	UserBaseDN        string
	GroupBaseDN       string
	UserClass         string
	SearchAttribute   string
	UsernameAttribute string
	EmailAttribute    string
}

// NewLDAP creates a new instance of LDAP with the provided configuration.
func NewLDAP(cfg LDAPConfig) LDAP {
	return LDAP{
		cfg: cfg,
	}
}

func (l *LDAPConfig) connect() *ldap.Conn {
	var c *ldap.Conn
	var err error

	if l.SSL {
		c, err = ldap.DialTLS(
			"tcp",
			fmt.Sprintf("%s:%d", l.Server, l.Port),
			&tls.Config{
				InsecureSkipVerify: l.SkipVerify,
			},
		)
	} else {
		c, err = ldap.Dial(
			"tcp",
			fmt.Sprintf("%s:%d", l.Server, l.Port),
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
	for _, attr := range []string{l.cfg.UsernameAttribute, l.cfg.EmailAttribute} {
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
		l.cfg.UserClass,
		l.cfg.SearchAttribute,
		group,
		l.cfg.GroupBaseDN,
	)

	result, err := l.conn.Search(&ldap.SearchRequest{
		BaseDN:     l.cfg.UserBaseDN,
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
		if l.cfg.UsernameAttribute != "" {
			member.username = e.GetAttributeValue(l.cfg.UsernameAttribute)
			if member.username == "" {
				panic(fmt.Sprintf(
					"Failed to get username (%s) for %s",
					l.cfg.UsernameAttribute,
					e.DN,
				))
			}
		}
		if l.cfg.EmailAttribute != "" {
			member.email = e.GetAttributeValue(l.cfg.EmailAttribute)
			if member.email == "" {
				panic(fmt.Sprintf(
					"Failed to get e-mail (%s) for %s",
					l.cfg.EmailAttribute,
					e.DN,
				))
			}
		}

		members = append(
			members,
			User{
				username: e.GetAttributeValue(l.cfg.UsernameAttribute),
				email:    e.GetAttributeValue(l.cfg.EmailAttribute),
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
