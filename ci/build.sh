#!/usr/bin/env bash

# Build the project, injectingthe version variable from git describe output

go build -i -v -ldflags="-X stash.jamf.build/devops/groupsync/cmd.version=$(git describe --always --dirty)" stash.jamf.build/devops/groupsync
