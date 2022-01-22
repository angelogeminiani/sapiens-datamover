echo "START BUILDING WINDOWS 32 VERSION $BASE.$BUILD..."
set GOARCH=386 go build -ldflags="-H windowsgui" -o ./__build/windows-386/datamover.exe ./main.go
