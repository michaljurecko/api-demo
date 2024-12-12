#!/usr/bin/env bash

set -o errexit          # Exit on most errors (see the manual)
set -o errtrace         # Make sure any error trap is inherited
set -o nounset          # Disallow expansion of unset variables
set -o pipefail         # Use last non-zero exit code in a pipeline
#set -o xtrace          # Trace the execution of the script (debug)

TEST_TIMEOUT=${TEST_TIMEOUT:-600s}

# Ensure at least one argument is provided
if [[ $# -eq 0 ]]; then
  echo "Error: At least one directory path must be provided as an argument."
  exit 1
fi

# Open project directory
cd "$(dirname "$0")/.."

# Iterate through each specified directory and run tests
for dir in "$@"; do
  if [[ -d "$dir" ]]; then
    # 'race' flag enables detection of issues in locking primitives
    # 'tparse' tool is used to pretty print test output
    echo "Running tests in '$dir'..."
    go test -timeout "$TEST_TIMEOUT" -race -bench=. -json "$dir/..." | tparse -progress -format plain
  else
    echo "Error: '$dir' is not a valid directory."
    exit 1
  fi
done
