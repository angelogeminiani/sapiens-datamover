##!/bin/sh

BASE="1.1"
# Get latest build
BUILD=$(<build_version.txt)
chrlen=${#BUILD}
if [ $chrlen = 0 ]
then
  BUILD=0
fi

echo "START BUILDING LINUX VERSION $BASE.$BUILD..."

## linux
env GOOS=linux GOARCH=386 go build -o ./_build/linux/datamover ./main.go

