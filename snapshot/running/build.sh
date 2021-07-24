#!/usr/bin/env bash

set -euo pipefail

monoclean deps
monoclean clean

(
  set +e

  monoclean run ./exit.ts 0
  echo "exit code expected=0 actual=$?"

  monoclean run ./exit.ts 5
  echo "exit code expected=5 actual=$?"
)

tmppath=$(monoclean run --build-only ./exit.ts 5)
mv "$tmppath" ./out/tmp/run001
