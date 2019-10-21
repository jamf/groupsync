package services

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

var client = LDAP{
	cfg: LDAPConfig{
		port:       389,
		server:     "127.0.0.1",
		ssl:        false,
		skipVerify: false,

		bindUser:     "cn=admin,dc=planetexpress,dc=com",
		bindPassword: "GoodNewsEveryone",

		userBaseDN:        "ou=people,dc=planetexpress,dc=com",
		groupBaseDN:       "ou=people,dc=planetexpress,dc=com",
		userClass:         "person",
		searchAttribute:   "memberOf",
		usernameAttribute: "uid",
		emailAttribute:    "mail",
	},
}

var sslClient = LDAP{
	cfg: LDAPConfig{
		port:       636,
		server:     "127.0.0.1",
		ssl:        true,
		skipVerify: true,

		bindUser:     "cn=admin,dc=planetexpress,dc=com",
		bindPassword: "GoodNewsEveryone",

		userBaseDN:        "ou=people,dc=planetexpress,dc=com",
		groupBaseDN:       "ou=people,dc=planetexpress,dc=com",
		userClass:         "person",
		searchAttribute:   "memberOf",
		usernameAttribute: "uid",
		emailAttribute:    "mail",
	},
}

// Test cases

func TestLDAP(t *testing.T) {
	ldapTeardown := setupLDAPService(t)
	defer ldapTeardown(t)

	testClient(t, client)
	testClient(t, sslClient)
}

// Helpers

func testClient(t *testing.T, client LDAP) {
	client.connect()
	defer client.close()

	actualResults, err := client.GroupMembers("ship_crew")
	if err != nil {
		panic(err)
	}

	expectedResults := []User{
		User{
			username: "bender",
			email:    "bender@planetexpress.com",
		},
		User{
			username: "fry",
			email:    "fry@planetexpress.com",
		},
		User{
			username: "leela",
			email:    "leela@planetexpress.com",
		},
	}

	expectedResultsOrig := expectedResults

	if len(actualResults) != len(expectedResults) {
		panic(fmt.Sprintf(
			"Actual and expected results differ.\nActual:\n%+v\nExpected:\n%+v\n",
			actualResults,
			expectedResults,
		))
	}

	for _, member := range actualResults {
		t.Log(member)
		for i := range expectedResults {
			if reflect.DeepEqual(expectedResults[i], member) {
				expectedResults = append(
					expectedResults[:i],
					expectedResults[i+1:]...,
				)
				break
			}
		}
	}

	if len(expectedResults) > 0 {
		panic(fmt.Sprintf(
			"Actual and expected results differ.\nActual:\n%+v\nExpected:\n%+v\n",
			actualResults,
			expectedResultsOrig,
		))
	}
}

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
