##!/bin/sh

BASE="1.1"
# Get latest build
BUILD=$(<build_version.txt)
chrlen=${#BUILD}
if [ $chrlen = 0 ]
then
  BUILD=0
fi

echo "START BUILDING WINDOWS 64 VERSION $BASE.$BUILD..."

## windows
env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-H windowsgui" -o ./_build/windows64/datamover.exe ./main.go
