
```go
// pb命令
//go:generate protoc --go_out=. data.proto

// grpc命令
//go:generate protoc --go_out=. --go-grpc_out=. data.proto
```


```sh
# 交叉编译cgo程序, 需要安装filosottile/musl-cross/musl-cross
GOOS=linux CC="/usr/local/bin/arm-linux-musleabi-gcc" GOARCH=arm CGO_ENABLED=1 go build -o koe_pi_arm -ldflags "-linkmode external -extldflags -static" main.go
```