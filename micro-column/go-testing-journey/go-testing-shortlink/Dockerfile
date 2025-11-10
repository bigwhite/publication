# --- Builder Stage ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 拷贝 go.mod 和 go.sum 文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 拷贝所有源代码
COPY . .

# 编译应用，生成静态链接的二进制文件
RUN CGO_ENABLED=0 go build -o /shortlink-server ./cmd/server

# --- Final Stage ---
FROM gcr.io/distroless/static-debian11

# 从 builder 阶段拷贝二进制文件
COPY --from=builder /shortlink-server /

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["/shortlink-server"]
