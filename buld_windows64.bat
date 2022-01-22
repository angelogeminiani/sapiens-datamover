echo "START BUILDING WINDOWS 32 VERSION $BASE.$BUILD..."
set GOARCH=amd64 go build CGO_ENABLED=1 -ldflags="-H windowsgui" -o ./__build/windows-amd64/datamover.exe ./main.go
