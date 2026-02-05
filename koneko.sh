#!/bin/bash

set -e

rm -f bin/koneko
CGO_ENABLED=0 go build -mod=vendor -ldflags "-s -w" -o bin/koneko ./cmd/koneko
./bin/koneko