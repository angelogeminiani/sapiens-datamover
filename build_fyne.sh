##!/bin/sh


echo "START BUILDING WITH FYNE-CROSS..."

# windows
sh ./build_fyne_windows.sh

# linux
sh ./build_fyne_linux.sh

# mac
sh ./build_fyne_mac.sh

echo "remove temp files ..."
rm -r ./fyne-cross