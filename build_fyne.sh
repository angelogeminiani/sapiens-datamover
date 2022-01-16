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

# windows
echo "  * WINDOWS (amd64,386)"
fyne-cross windows -arch=amd64,386
echo "    copy executables ..."
cp ./fyne-cross/bin/windows-386/gg-progr-datamover.exe ./__build/windows-386/datamover.exe
cp ./fyne-cross/bin/windows-amd64/gg-progr-datamover.exe ./__build/windows-amd64/datamover.exe

# linux
echo "  * LINUX (amd64,386,arm,arm64)"
fyne-cross linux -arch=amd64,386,arm,arm64
echo "    copy executables ..."
cp ./fyne-cross/bin/linux-386/gg-progr-datamover ./__build/linux-386/datamover
cp ./fyne-cross/bin/linux-amd64/gg-progr-datamover ./__build/linux-amd64/datamover
cp ./fyne-cross/bin/linux-arm/gg-progr-datamover ./__build/linux-arm/datamover
cp ./fyne-cross/bin/linux-arm64/gg-progr-datamover ./__build/linux-arm64/datamover

# mac
echo "  * DARWIN (arm64)"
fyne-cross darwin -arch=arm64 -app-build=1 -app-version="$BASE.$BUILD" -app-id="datamover"
echo "    copy executables ..."
cp ./fyne-cross/dist/darwin-arm64/gg-progr-datamover.app/Contents/MacOS/gg-progr-datamover ./__build/darwin-arm64/datamover

echo "remove temp files ..."
rm -r ./fyne-cross