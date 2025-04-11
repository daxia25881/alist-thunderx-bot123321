FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制go.mod和go.sum文件（如果存在）
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# 复制源代码
COPY *.go ./

# 构建应用
RUN go build -o telegram-bot .

# 使用较小的基础镜像
FROM alpine:latest

WORKDIR /app

# 安装CA证书
RUN apk --no-cache add ca-certificates

# 从builder阶段复制编译好的应用
COPY --from=builder /app/telegram-bot /app/
COPY config.json /app/

# 运行应用
CMD ["/app/telegram-bot"] 