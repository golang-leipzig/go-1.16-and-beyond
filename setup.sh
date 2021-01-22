#!/usr/bin/env bash

set -euo pipefail

mkdir -p go
GOPATH=$PWD/go go get golang.org/dl/go1.16rc1
./go/bin/go1.16rc1 download
