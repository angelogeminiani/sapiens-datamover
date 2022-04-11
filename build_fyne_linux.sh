##!/bin/sh

BASE="1.1"
# Get latest build
BUILD=$(<build_version.txt)
chrlen=${#BUILD}
if [ $chrlen = 0 ]
then
  BUILD=0
fi

echo "START BUILDING WITH FYNE-CROSS LINUX $BASE.$BUILD..."

# linux
echo "  * LINUX (amd64,386,arm,arm64)"
fyne-cross linux -arch=amd64,386,arm,arm64
echo "    copy executables ..."
cp ./fyne-cross/bin/linux-386/sapiens-datamover ./__build/linux-386/datamover
cp ./fyne-cross/bin/linux-amd64/sapiens-datamover ./__build/linux-amd64/datamover
cp ./fyne-cross/bin/linux-arm/sapiens-datamover ./__build/linux-arm/datamover
cp ./fyne-cross/bin/linux-arm64/sapiens-datamover ./__build/linux-arm64/datamover
