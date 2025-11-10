FROM golang:1.25-alpine

WORKDIR /app

# 将 user-service 的源代码复制进去
COPY . .

# 编译 user-service
RUN go build -o /user-service .

EXPOSE 8090
CMD ["/user-service"]
