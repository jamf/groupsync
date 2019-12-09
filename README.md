# groupsync
*groupsync* is a CLI tool for syncing user group membership info from a directory (like LDAP) to a service like GitHub.

## Status
Preparing for the first (`0.1.0`) release. Still got to work some kinks out, refactor, test. Use at your own peril!

## Installation
Right now *groupsync* is available only from source. The plan is to add binaries and a docker image soon.

You'll need *git* and the *[go toolchain](https://golang.org/doc/install)*.

```sh
go get github.com/jamf/groupsync
```

## Config
In order to let groupsync know how to access the services you're trying to sync data between, you'll need to provide a config file. On UNIX-ish systems this configuration file should be `~/.groupsync/groupsync.yaml`.

Here's an [example config file](examples/groupsync.yaml).

The `groupsync ls` subcommand is ideal for testing the connection.

## Usage
### List users in a group
```sh
groupsync ls ldap my-group
```

```sh
groupsync ls github my-team1 my-team2 my-team3
```

*Note:* the failure exit code (`1`) will only be returned when looking up a single user group. If looking up multiples, the exit code will always be `0`.

### Perform a dry run of sync
```sh
groupsync sync -d "ldap:my-group" "github:my-team"
```

If all looks good, remove the `-d` flag to actually commit the changes.

### Sync from multiple sources
```sh
groupsync sync "ldap:my-group1" "ldap:my-group2" "github:my-source-team" "github:my-target-team"
```

All the groups/teams provided are treated as sources except for the last one, which is the target.

The members of all source group are collected and then the resulting list of accounts is synced into the target group.

### Mapping files
Instead of providing the group/team names to sync using command line arguments, you can provide a file with all the mappings like so:

```
groupsync sync -m mappings.yaml
```

Here's an [example mappings file](examples/mappings.yaml). Note that it contains multiple mappings.

## Hacking
Some developer's documentation (and possibly reference, but probably only as an afterthought) is in the plans - particularly a guide to adding new services. For now if you wish to implement any specific functionality, you might have to inquire a bit. Don't hesitate to get in touch by opening an issue!

A refactor is probably in order too. Please don't hesitate to point out unreadable portions of code.