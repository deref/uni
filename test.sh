#!/usr/bin/env bash

set -euo pipefail

go build

export PATH="$PWD:$PATH"

if [[ $# == 0 ]]; then
  snapshots=$(ls snapshot)
else
  snapshots="$@"
fi

for snapshot in $snapshots; do
  (
    cd "snapshot/$snapshot"
    pwd
    if ! ./build.sh 2>stderr >stdout; then
      echo "ERROR"
      exit 1
    fi

    # Cleanup some non-determinsim in the output.
    perl -pi -e 's/^audited (\d+) packages in (.*?)s$/audited $1 packages in SOME_AMOUNT_OF_TIME/g' stdout
    perl -pi -e 's/^(\d+) packages are looking for funding$/SOME_NUMBER_OF packages are looking for funding/g' stdout
  )
done

echo "done"
