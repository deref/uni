#!/usr/bin/env bash

set -euo pipefail

monoclean deps
monoclean clean
monoclean build
