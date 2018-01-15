#!/bin/bash
set -e
cd "$(dirname "$0")"
source ../init.sh

go test
