#!/usr/bin/env bash

set -o errexit          # Exit on most errors (see the manual)
set -o errtrace         # Make sure any error trap is inherited
set -o nounset          # Disallow expansion of unset variables
set -o pipefail         # Use last non-zero exit code in a pipeline
#set -o xtrace          # Trace the execution of the script (debug)

# Open project directory
cd "$(dirname "$0")/.."

go mod tidy
go mod vendor
exec golangci-lint run --fix -c "./build/ci/lint.yaml" "$@"
