#!/usr/bin/env bash

set -o errexit          # Exit on most errors (see the manual)
set -o errtrace         # Make sure any error trap is inherited
set -o nounset          # Disallow expansion of unset variables
set -o pipefail         # Use last non-zero exit code in a pipeline
#set -o xtrace          # Trace the execution of the script (debug)

# Open project directory
cd "$(dirname "$0")/.."

# Prevent file permission issues between container and host.
# By default, files created inside a container in a mounted volume are owned by root.
# To avoid this, the container is run without root privileges,
# using the same UID/GID as the local user, even if Docker is in root mode (default).
# Podman and Docker in rootless mode don't have this issue.
USER="$(id -u):$(id -g)"

exec docker compose -f ./.devcontainer/compose.yaml run --rm -u "$USER" --service-ports dev
