# Adding new services
Services act as sources of group membership information. Before getting to
work, you may want to take a look at how
[one is already defined](../services/ldap.go).

## Code things
Here are the steps to defining a new service:

1. Implement the [Service interface](../services/service.go). You'll need to
   define how a service-specific unique ID is acquired, and how to get a list
   of users for a given group name the `GroupMembers` method.
2. Define a config struct for your service, then hook that up to the global Config type
   in [services/config](../services/config.go). This data will be deserialized
   from the the config `.yaml` file provided by the user -
   [here's an example](../examples/groupsync.yaml).
3. Remember to add your service to the
   [`SvcFromString` function](../services/service.go).
4. If you expect to use this service as a source for sync, go through possible
   targets (like GitHub?) and make sure they know how to convert the user
   identity acquired from your service to the target identity - that logic lives
   in the `acquireIdentity` method; an example implementation is in
   [ldap.go](../services/ldap.go).
5. Write tests specific to your service if at all possible. For inspiration,
   look at what we do for [LDAP](../services/ldap_test.go) - we spin up a
   docker container with an OpenLDAP server that contains some test data,
   and then test against that. The CI environment has a docker daemon available;
   go nuts with it.

## Testing
After all that is done, get a dev build going:

```
ci/build.sh
```

Then test if the following work:

```
./groupsync ls yourservice
./groupsync sync -d "yourservice:yourgroup" "otherservice:othergroup"
```

## Getting help
Stuck? Go ahead! Raise an issue!