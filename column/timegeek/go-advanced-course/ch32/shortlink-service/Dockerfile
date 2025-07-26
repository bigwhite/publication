# Dockerfile

# --- Build Stage (builder) ---
# 使用一个包含Go编译环境的官方Alpine镜像作为构建基础，它体积较小。
# 请确保这里的Go版本与你项目go.mod中指定的版本或你本地开发使用的版本一致或兼容。
FROM golang:1.24.5-alpine3.21 AS builder
# (或者你可以用一个更精确的 patch 版本，如 golang:1.21.5-alpine)

# 设置Go环境变量，确保模块模式开启，并为静态链接做准备
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
# (GOARCH可以根据你的目标部署平台调整，amd64是常见的)

# 设置工作目录
WORKDIR /build

# 优化依赖下载：先拷贝go.mod和go.sum，运行go mod download。
# 这样，如果这两个文件没有变化，Docker可以利用缓存，跳过重新下载依赖。
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# 拷贝项目的其他所有源代码到工作目录
COPY . .

# 编译Go应用。
# 我们将编译 cmd/server/main.go 作为我们服务的主入口。
# -ldflags="-w -s":
#   -w: 省略DWARF调试信息。
#   -s: 省略符号表。
#   这两个标志可以显著减小最终二进制文件的大小。
# -o /app/shortlink-server: 指定编译输出的二进制文件名为shortlink-server，并放在/app目录下。
#   (注意：/app目录是在这个builder阶段创建的，与最终运行阶段的WORKDIR可以不同)
RUN go build -ldflags="-w -s" -o /app/shortlink-server ./cmd/server/main.go


# --- Runtime Stage (final) ---
# 使用一个极小的基础镜像。
# scratch 是一个完全空的镜像，如果你的二进制是完全静态链接且无任何外部文件依赖，这是最佳选择。
FROM scratch
# 或者，如果你的应用需要一些基础的系统工具或CA证书，可以使用alpine：
# FROM alpine:latest
# RUN apk --no-cache add ca-certificates tzdata # 添加CA证书和时区数据
# 或者使用Google的distroless镜像，它只包含运行时依赖，安全性更高：
# FROM gcr.io/distroless/static:nonroot # 用于静态链接的Go应用，并以非root用户运行

# 设置最终镜像中的工作目录
WORKDIR /app

# 从构建阶段 (builder) 拷贝编译好的二进制文件到当前阶段的/app目录下
COPY --from=builder /app/shortlink-server /app/shortlink-server

# (可选) 如果你的应用依赖配置文件，并且你希望将它们打包到镜像中
# (虽然更推荐的做法是通过ConfigMap或Secret在运行时挂载配置)
# COPY --from=builder /build/configs/config.production.yaml /app/configs/config.yaml

# (可选) 如果你在builder阶段创建了非root用户，并希望用它运行
# 例如，在builder阶段: RUN addgroup -S appgroup && adduser -S appuser -G appgroup
# 然后在这里:
# COPY --from=builder /etc/passwd /etc/passwd
# COPY --from=builder /etc/group /etc/group
# USER appuser:appgroup

# 声明应用在容器内监听的端口 (与应用配置中的端口一致)
# 这只是一个元数据声明，实际端口映射在docker run或Kubernetes Service中定义
EXPOSE 8080 
# (假设我们的shortlink-service默认或通过配置监听8080端口)

# 定义容器启动时执行的命令
# ENTRYPOINT使得容器像一个可执行文件一样运行。
# CMD可以为ENTRYPOINT提供默认参数，或在ENTRYPOINT未定义时作为主命令。
# 对于Go应用，通常直接将编译好的二进制文件作为ENTRYPOINT。
ENTRYPOINT ["/app/shortlink-server"]
# (可选) 如果你的应用需要启动参数，可以在这里提供默认值
# CMD ["-config", "/app/configs/config.yaml"] 
# 或者，这些参数通常在运行时（docker run 或 K8s manifest）中提供。
