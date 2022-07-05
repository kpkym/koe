# 编译镜像
FROM golang:1.18.3-alpine3.16 AS build
RUN apk add -U --no-cache ca-certificates

WORKDIR /work
COPY . .
RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go build -tags postgresql -o /work/server main.go

# 运行镜像
FROM alpine
ENV TZ Asia/Shanghai

# 定义工作目录为work
WORKDIR /work

COPY --from=build /work/server ./server
COPY --from=build /work/config.toml ./config.toml
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 开放http 8080端口
EXPOSE 8080
# 启动http服务
ENTRYPOINT ["./server", "web", "8080"]
