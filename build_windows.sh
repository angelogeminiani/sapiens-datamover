##!/bin/sh

BASE="1.1"
# Get latest build
BUILD=$(<build_version.txt)
chrlen=${#BUILD}
if [ $chrlen = 0 ]
then
  BUILD=0
fi

echo "START BUILDING WINDOWS 32 VERSION $BASE.$BUILD..."

## windows
env GOOS=windows GOARCH=386 go build -ldflags="-H windowsgui" -o ./__build/windows32/datamover.exe ./main.go
