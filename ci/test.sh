#!/usr/bin/env bash

# Test the project, injecting the debug tag.

go test -tags debug -ldflags="-X github.com/jamf/groupsync/cmd.version=$(git describe --always --dirty)" ./...
