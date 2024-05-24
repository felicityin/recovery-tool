
# Description

 Cross platform shared library.

```
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "-s -w" -buildmode=c-shared -o recovery_tool.dll main.go 

GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=c-shared -o recovery_tool.dylib main.go

GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -buildmode=c-shared -o lib_recovery_tool.dylib main.go

CGO_ENABLED=1 GOARCH=arm64  GOOS=linux   CC=arm-linux-musleabihf-gcc CGO_LDFLAGS="-static" go build -a -v -installsuffix cgo -o bin/bfs-data-detection .

go build -buildmode=c-shared -o recovery_tool.dylib main.go
```