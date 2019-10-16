#!/usr/bin/env bash

if [[ $# -eq 0 ]] ; then
    echo 'Please provide the version to tag groupsync with.'
    exit 1
fi

# Tag groupsync with the provided version.

VERSION=$1

git tag -a $VERSION -m "$VERSION"
