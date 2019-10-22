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
	Port       int32
	Server     string
	SSL        bool
	SkipVerify bool `mapstructure:"skip_verify"`

	// Auth
	BindUser     string `mapstructure:"bind_user"`
	BindPassword string `mapstructure:"bind_password"`

	// Schema
	UserBaseDN        string `mapstructure:"user_base_dn"`
	GroupBaseDN       string `mapstructure:"group_base_dn"`
	UserClass         string `mapstructure:"user_class"`
	SearchAttribute   string `mapstructure:"search_attribute"`
	UsernameAttribute string `mapstructure:"username_attribute"`
	EmailAttribute    string `mapstructure:"email_attribute"`
}

// NewLDAP creates a new instance of LDAP with the provided configuration.
func NewLDAP() *LDAP {
	return &LDAP{
		cfg: cfg.LDAP,
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

	err = c.Bind(l.BindUser, l.BindPassword)
	if err != nil {
		panic(err)
	}

	return c
}

func (l *LDAP) connect() {
	l.conn = l.cfg.connect()
}

// GroupMembers returns the members of group `group` as a slice of User
// instances. Implements the Service interface.
func (l *LDAP) GroupMembers(group string) ([]User, error) {
	l.connect()
	defer l.close()

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

	// Get the DN of the group.
	group, err := l.findGroup(group)
	if err != nil {
		return nil, err
	}

	filter := fmt.Sprintf(
		"(&(objectClass=%s)(%s=%s))",
		l.cfg.UserClass,
		l.cfg.SearchAttribute,
		group,
	)

	result, err := l.conn.Search(&ldap.SearchRequest{
		BaseDN:       l.cfg.UserBaseDN,
		Filter:       filter,
		Scope:        2,
		DerefAliases: 1,
		Attributes:   attrs,
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

// Returns the DN of an LDAP group or an error if not found.
func (l *LDAP) findGroup(g string) (string, error) {
	if l.conn == nil {
		l.connect()
		defer l.close()
	}

	filter := fmt.Sprintf(
		"(&(objectClass=group)(cn=%s))",
		g,
	)

	result, err := l.conn.Search(&ldap.SearchRequest{
		BaseDN:       l.cfg.GroupBaseDN,
		Filter:       filter,
		Scope:        2,
		DerefAliases: 1,
	})
	if err != nil {
		return "", fmt.Errorf("error looking up group %s: %s", g, err)
	}

	if len(result.Entries) < 0 {
		return "", fmt.Errorf("group `%s` not found", g)
	} else if len(result.Entries) > 1 {
		return "", fmt.Errorf("multiple groups found for `%s`", g)
	}

	return result.Entries[0].DN, nil
}

func (l *LDAP) close() {
	if l.conn != nil {
		l.conn.Close()
		l.conn = nil
	}
}
