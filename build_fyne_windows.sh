##!/bin/sh

BASE="1.1"
# Get latest build
BUILD=$(<build_version.txt)
chrlen=${#BUILD}
if [ $chrlen = 0 ]
then
  BUILD=0
fi

echo "START BUILDING WITH FYNE-CROSS WINDOWS $BASE.$BUILD..."

# windows
echo "  * WINDOWS (amd64,386)"
fyne-cross windows -arch=amd64,386
echo "    copy executables ..."
cp ./fyne-cross/bin/windows-386/sapiens-datamover.exe ./__build/windows-386/datamover.exe
cp ./fyne-cross/bin/windows-amd64/sapiens-datamover.exe ./__build/windows-amd64/datamover.exe