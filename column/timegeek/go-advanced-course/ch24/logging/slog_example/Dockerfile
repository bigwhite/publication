# --- Stage 1: Build Stage ---
FROM golang:1.21.7-alpine3.19 AS builder

# 设置Go环境变量，确保模块模式开启，并为静态链接做准备
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# 设置工作目录，所有后续命令都将在此目录下执行
WORKDIR /app
COPY go.mod ./
RUN go mod download && go mod verify

COPY . .

# 编译Go应用。
RUN go build -ldflags="-w -s" -o /app/myslogapp .


# --- Stage 2: Runtime Stage ---
# 使用一个极小的基础镜像。
# `scratch` 是一个完全空的镜像，不包含任何文件系统和用户空间工具。
FROM scratch

# 设置最终镜像中的工作目录
WORKDIR /app

# 从构建阶段 (builder) 拷贝编译好的二进制文件到当前阶段的/app目录下
COPY --from=builder /app/myslogapp /app/myslogapp

EXPOSE 8088

# 定义容器启动时执行的命令
ENTRYPOINT ["/app/myslogapp"]
