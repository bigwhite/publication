# Dockerfile

# --- 第一阶段: 构建环境 (builder) ---
FROM golang:1.25.4-alpine3.22 AS builder
             
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/server ./main.go


# --- 第二阶段: 生产环境 (final) ---
FROM gcr.io/distroless/static-debian11

COPY --from=builder /app/server /server

EXPOSE 8080

CMD ["/server", "web", "--port", "8080", "api"]
