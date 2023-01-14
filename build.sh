#!/bin/bash

set -eo pipefail

if [ -z $1 ]; then
  echo "Usage: ./build.sh [program]"
  exit 1
fi

mkdir -p bin

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/$1 $1/main.go
