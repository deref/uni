#!/usr/bin/env bash

set -euo pipefail

go build

export PATH="$PWD:$PATH"

for snapshot in $(ls snapshot); do
  (
    cd "snapshot/$snapshot"
    pwd
    ./build.sh 2>stderr >stdout || true
  )
done
