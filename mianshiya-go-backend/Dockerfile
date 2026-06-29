# 第一阶段：编译 Go 二进制
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 允许自动下载所需 Go 版本
ENV GOTOOLCHAIN=auto

# 先复制依赖文件，利用 Docker 缓存层
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并编译
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/server/

# 第二阶段：最小运行镜像
FROM alpine:3.20

# HTTPS 请求需要 CA 证书（AI API）
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# 复制二进制和配置文件
COPY --from=builder /app/server .
COPY --from=builder /app/config.yaml .

# 创建上传文件目录
RUN mkdir -p /app/uploads

EXPOSE 8101

CMD ["./server"]
