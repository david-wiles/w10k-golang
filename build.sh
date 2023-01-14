#!/bin/bash

set -eo pipefail

mkdir -p bin

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/broadcast broadcast/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/client2client client2client/main.go
