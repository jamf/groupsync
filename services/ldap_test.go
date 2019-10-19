package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

// Test cases

var sslClient = LDAP{
	port:       636,
	server:     "127.0.0.1",
	ssl:        true,
	skipVerify: true,

	bindUser:     "cn=admin,dc=planetexpress,dc=com",
	bindPassword: "GoodNewsEveryone",

	userBaseDN:      "ou=people,dc=planetexpress,dc=com",
	groupBaseDN:     "ou=people,dc=planetexpress,dc=com",
	userClass:       "person",
	searchAttribute: "memberOf",
}

func TestLDAP(t *testing.T) {
	ldapTeardown := setupLDAPService(t)
	defer ldapTeardown(t)

	sslClient.connect()
	defer sslClient.close()

	actualResults := sslClient.members("ship_crew").Entries

	expectedResults := map[string]bool{
		"cn=Bender Bending Rodríguez,ou=people,dc=planetexpress,dc=com": false,
		"cn=Philip J. Fry,ou=people,dc=planetexpress,dc=com":            false,
		"cn=Turanga Leela,ou=people,dc=planetexpress,dc=com":            false,
	}

	for _, member := range actualResults {
		t.Log(member.DN)
		expectedResults[member.DN] = true
	}

	for person, exists := range expectedResults {
		if !exists {
			panic(fmt.Sprintf(
				"One of the expected search results wasn't there:\n%s",
				person,
			))
		}
	}
}

// Helpers

func setupLDAPService(t *testing.T) func(t *testing.T) {
	t.Log("Setting up an LDAP server container...")

	d, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}

	port, err := nat.NewPort("tcp", "389")
	if err != nil {
		panic(err)
	}

	sslPort, err := nat.NewPort("tcp", "636")
	if err != nil {
		panic(err)
	}

	_, err = d.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: "rroemhild/test-openldap",
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				port: []nat.PortBinding{nat.PortBinding{
					HostIP:   "127.0.0.1",
					HostPort: "389",
				}},
				sslPort: []nat.PortBinding{nat.PortBinding{
					HostIP:   "127.0.0.1",
					HostPort: "636",
				}},
			},
		},
		&network.NetworkingConfig{},
		"ldap_test_server",
	)
	if err != nil {
		panic(err)
	}

	err = d.ContainerStart(
		context.Background(),
		"ldap_test_server",
		types.ContainerStartOptions{},
	)
	if err != nil {
		panic(err)
	}

	// Wait for the container to be ready. This isn't ideal, I know.
	time.Sleep(5 * time.Second)

	// Return a teardown function.
	return func(t *testing.T) {
		t.Log("Tearing down the LDAP server container...")

		err := d.ContainerRemove(
			context.Background(),
			"ldap_test_server",
			types.ContainerRemoveOptions{
				Force: true,
			},
		)
		if err != nil {
			panic(err)
		}
	}
}
