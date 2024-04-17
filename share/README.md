
# Description

 Cross platform shared library.

```
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "-s -w" -buildmode=c-shared -o recovery_tool.dll main.go 

GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=c-shared -o recovery_tool.dylib main.go

```