##!/bin/sh

BASE="1.1"
# Get latest build
BUILD=$(<build_version.txt)
chrlen=${#BUILD}
if [ $chrlen = 0 ]
then
  BUILD=0
fi

echo "START BUILDING WITH FYNE-CROSS (https://github.com/fyne-io/fyne-cross) $BASE.$BUILD..."

# mac
echo "  * DARWIN (arm64)"
fyne-cross darwin -arch=arm64 -app-build=1 -app-version="$BASE.$BUILD" -app-id="datamover"
echo "    copy executables ..."
cp ./fyne-cross/dist/darwin-arm64/sapiens-datamover.app/Contents/MacOS/gg-progr-datamover ./__build/darwin-arm64/datamover
