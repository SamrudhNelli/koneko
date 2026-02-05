#!/bin/bash

set -e

rm -f koneko
CGO_ENABLED=0 go build -ldflags "-s -w" -o koneko
./koneko