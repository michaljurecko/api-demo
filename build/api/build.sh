#!/bin/bash

set -o errexit          # Exit on most errors (see the manual)
set -o errtrace         # Make sure any error trap is inherited
set -o nounset          # Disallow expansion of unset variables
set -o pipefail         # Use last non-zero exit code in a pipeline
#set -o xtrace          # Trace the execution of the script (debug)


# Open project directory
cd "$(dirname "$0")/../.."

docker build -f ./.devcontainer/Containerfile -t api-demo-dev:latest .

docker build -f ./build/api/Containerfile -t api-demo:latest --build-arg DEV_IMAGE=api-demo-dev:latest .

