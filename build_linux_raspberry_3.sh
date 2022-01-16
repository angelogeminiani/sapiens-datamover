##!/bin/sh

BASE="1.1"
# Get latest build
BUILD=$(<build_version.txt)
chrlen=${#BUILD}
if [ $chrlen = 0 ]
then
  BUILD=0
fi

echo "START BUILDING LINUX RASPBERRY VERSION $BASE.$BUILD..."

## linux
env GOOS=linux GOARCH=arm GOARM=5 go build -o ./__build/linux_raspberry_3/datamover ./main.go
