
# Description

 Cross platform shared library.

```
//动态链接库
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "-s -w" -buildmode=c-shared -o recovery_tool.dll main.go 

GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "-s -w" -buildmode=c-shared -o lib_recovery_tool.dylib main.go

GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -ldflags "-s -w" -buildmode=c-shared -o lib_recovery_tool_arm64.dylib main.go

CGO_ENABLED=1 GOARCH=arm64  GOOS=linux   CC=arm-linux-musleabihf-gcc CGO_LDFLAGS="-static" go build -a -v -installsuffix cgo -o bin/bfs-data-detection .

//静态链接库
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=c-archive -o lib_recovery_tool.a main.go


```