#!/bin/bash
set -e
cd "$(dirname "$0")"
source ../init.sh

VERSION=$(cat ../VERSION)

GITREV=0
if [ -d ../.git ]; then
  GITREV=$(git -C .. rev-parse --short=10 HEAD)
fi

go build \
  -buildmode=exe \
  -ldflags "-X main.versionGit=$GITREV -X main.version=$VERSION" \
  -o $SRCDIR/build/fontctrl

# -a to force rebuild all packages
# -v to print package names as they are built
