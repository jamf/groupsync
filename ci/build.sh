go build -i -v -ldflags="-X stash.jamf.build/devops/groupsync/cmd.version=$(git describe --always --dirty)"
