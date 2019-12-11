#!/usr/bin/env bash

# Build the project, injecting the version variable from git describe output.
# This creates a dev (debug) build meant for both CI and local testing.

go build -i -v -tags debug -ldflags="-X github.com/jamf/groupsync/cmd.version=$(git describe --always --dirty)" github.com/jamf/groupsync
