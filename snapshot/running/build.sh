#!/usr/bin/env bash

set -euo pipefail

uni deps
uni clean

(
  set +e

  uni run ./exit.ts 0
  echo "exit code expected=0 actual=$?"

  uni run ./exit.ts 5
  echo "exit code expected=5 actual=$?"
)

tmppath=$(uni run --build-only ./exit.ts 5)
mv "$tmppath" ./out/tmp/run001
