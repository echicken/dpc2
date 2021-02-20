#!/bin/bash
# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

# Change into that directory
cd "$DIR"

go build -o ./build/linux_x64/doorparty-connector ./cmd/doorparty-connector.go
GOARM=6 GOARCH=arm go build -o ./build/linux_arm6/doorparty-connector ./cmd/doorparty-connector.go
GOOS=windows GOARCH=386 go build -o ./build/win32/doorparty-connector.exe ./cmd/doorparty-connector.go
