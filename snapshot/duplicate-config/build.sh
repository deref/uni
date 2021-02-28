#!/usr/bin/env bash

set -euo pipefail

(
  set +e
  uni env
  echo "exit code expected=1 actual=$?"
)
