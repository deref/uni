#!/usr/bin/env bash

set -euo pipefail

DIR=$(realpath `dirname $0`/../..)

cd $DIR/examples/simple-workspace

bash "$NVM_DIR/nvm.sh" use
monoclean deps
monoclean clean

set +e

monoclean lint
monoclean test
echo "exit code expected=0 actual=$?"
