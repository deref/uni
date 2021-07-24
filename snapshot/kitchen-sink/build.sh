#!/bin/bash

set -euo pipefail

bash "$NVM_DIR/nvm.sh" use

monoclean deps
monoclean clean
monoclean build
