#!/usr/bin/env bash

set -o errexit          # Exit on most errors (see the manual)
set -o errtrace         # Make sure any error trap is inherited
set -o nounset          # Disallow expansion of unset variables
set -o pipefail         # Use last non-zero exit code in a pipeline
#set -o xtrace          # Trace the execution of the script (debug)

# Open project directory
cd "$(dirname "$0")/.."

FORMAT="{{level_style (bold (uppercase (fixed_size 5 level)))}} {{blue (bold logger)}} {{msg}}"

# Go application logs JSON logs.
# fblog utility is used to print logs in a human-readable format.
# "go run" must run in a foreground process to receive signals.
set -m
exec 10> >(fblog --main-line-format "$FORMAT")
go run "$@" 2>&1 1>&10
