#!/bin/bash
SRCDIR_REL=$(dirname "${BASH_SOURCE[0]}")

# BUILD_DIR=$SRCDIR_REL/build
# if [[ "${BUILD_DIR:0:2}" == "./" ]]; then
#   BUILD_DIR=${BUILD_DIR:2}
# fi

pushd "$SRCDIR_REL" >/dev/null
export SRCDIR=$(pwd)
popd >/dev/null

PREV_GOPATH=$GOPATH

export GOPATH=$SRCDIR/build/gopath
export GOBIN=$GOPATH/bin

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  # Subshell
  set -e
  cd "$SRCDIR_REL"

  pushd client >/dev/null
  echo "client"
  go get -d -v .
  popd >/dev/null

  echo "DONE."

  if [[ "$GOPATH" != "$PREV_GOPATH" ]]; then
    echo "Source this script to setup the environment:" >&2
    if [ "$0" == "./init.sh" ]; then
      # pretty format for common case
      echo "  source init.sh"
    else
      echo "  source '$0'"
    fi
  fi
# else
#   # sourced
fi
