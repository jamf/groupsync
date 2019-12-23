#!/usr/bin/env bash

# Build the project, injecting the version variable from git describe output
# This creates a release build.

if [[ -z "$1" ]]; then
    version="$(git describe --tags --always --dirty)"
else
    version=$1    
fi

echo "Building a release version $version"

GOOS=darwin go build -i -v -ldflags="-X github.com/jamf/groupsync/cmd.version=$version" -o groupsync-darwin github.com/jamf/groupsync
GOOS=linux go build -i -v -ldflags="-X github.com/jamf/groupsync/cmd.version=$version" -o groupsync-linux github.com/jamf/groupsync
