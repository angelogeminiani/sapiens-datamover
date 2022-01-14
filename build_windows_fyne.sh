##!/bin/sh

BASE="1.1"
# Get latest build
BUILD=$(<build_version.txt)
chrlen=${#BUILD}
if [ $chrlen = 0 ]
then
  BUILD=0
fi

echo "START BUILDING WINDOWS VERSION WITH FYNE-CROSS (https://github.com/fyne-io/fyne-cross) $BASE.$BUILD..."

## windows
fyne-cross windows -arch=amd64,386
