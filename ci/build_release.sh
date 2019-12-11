#!/usr/bin/env bash

# Build the project, injecting the version variable from git describe output
# This creates a release build.

GOOS=darwin go build -i -v -ldflags="-X github.com/jamf/groupsync/cmd.version=$(git describe --always --dirty)" -o groupsync-darwin github.com/jamf/groupsync
GOOS=linux go build -i -v -ldflags="-X github.com/jamf/groupsync/cmd.version=$(git describe --always --dirty)" -o groupsync-linux github.com/jamf/groupsync
