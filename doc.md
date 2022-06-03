
```go
// pb命令
//go:generate protoc --go_out=. data.proto

// grpc命令
//go:generate protoc --go_out=. --go-grpc_out=. data.proto
```
