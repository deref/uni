#!/usr/bin/env bash

set -euo pipefail

(
  set +e
  monoclean env
  echo "exit code expected=1 actual=$?"
)
