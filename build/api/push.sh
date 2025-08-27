#!/bin/bash

set -o errexit          # Exit on most errors (see the manual)
set -o errtrace         # Make sure any error trap is inherited
set -o nounset          # Disallow expansion of unset variables
set -o pipefail         # Use last non-zero exit code in a pipeline
#set -o xtrace          # Trace the execution of the script (debug)


# Verify required environment variables
if [[ -z "${GITHUB_USERNAME:-}" ]]; then
    echo "Error: GITHUB_USERNAME environment variable is not set"
    exit 1
fi

if [[ -z "${GITHUB_PAT:-}" ]]; then
    echo "Error: GITHUB_PAT environment variable is not set"
    exit 1
fi

# Open project directory
cd "$(dirname "$0")/../.."

echo "$GITHUB_PAT" | docker login ghcr.io -u "$GITHUB_USERNAME" --password-stdin

REMOTE="ghcr.io/$GITHUB_USERNAME/api-demo:latest"
docker tag api-demo:latest "$REMOTE"
docker push "$REMOTE"
