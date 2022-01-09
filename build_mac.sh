##!/bin/sh

BASE="1.1"
# Get latest build
BUILD=$(<build_version.txt)
chrlen=${#BUILD}
if [ $chrlen = 0 ]
then
  BUILD=0
else
  BUILD=$(($BUILD + 1))
fi
echo $BUILD > ./build_version.txt
echo "START BUILDING MAC OSX VERSION $BASE.$BUILD..."

## mac
go build  -o ./_build/mac/datamover ./main.go
echo "END BUILDING VERSION $BASE.$BUILD."

cp ./build_version.txt ./_build

