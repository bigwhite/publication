FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main /app/main
RUN chmod +x /app/main
# COPY config.yaml /app/config.yaml # 如果有配置文件也可以拷贝
EXPOSE 8080
CMD ["/app/main"]
