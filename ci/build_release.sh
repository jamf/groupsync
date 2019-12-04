#!/usr/bin/env bash

# Build the project, injecting the version variable from git describe output
# This creates a release build.

go build -i -v -ldflags="-X github.com/jamf/groupsync/cmd.version=$(git describe --always --dirty)" github.com/jamf/groupsync
