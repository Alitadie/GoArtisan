# -------------------------------------------------------------------
# Builder Stage: 负责编译代码
# -------------------------------------------------------------------
FROM golang:1.25-alpine AS builder

# 安装构建必要的工具 (Git用于获取版本信息, Make用于执行构建脚本)
RUN apk add --no-cache git make

WORKDIR /app

# 1. 先下载依赖 (利用 Docker 缓存层，代码变动只要不改 go.mod 就不会重新下载依赖)
COPY go.mod go.sum ./
RUN go mod download

# 2. 拷贝源代码
COPY . .

# 3. 编译
# 这里我们将 Build 参数作为 Docker 构建参数传入
ARG VERSION=unknown
ARG COMMIT=unknown
ARG TIME=unknown

# 使用 -ldflags 注入版本信息 (手动复刻 Makefile 里的逻辑，确保 Windows 用户 docker build 也没问题)
RUN go build -ldflags "-s -w \
    -X 'go-artisan/pkg/version.GitTag=${VERSION}' \
    -X 'go-artisan/pkg/version.GitCommit=${COMMIT}' \
    -X 'go-artisan/pkg/version.BuildTime=${TIME}'" \
    -o /bin/server ./cmd/server

RUN go build -o /bin/artisan ./cmd/artisan

# -------------------------------------------------------------------
# Runner Stage: 纯净的运行环境
# -------------------------------------------------------------------
FROM alpine:latest

WORKDIR /app

# 安装基础库 (tzdata用于时区, ca-certificates用于HTTPS请求)
RUN apk add --no-cache tzdata ca-certificates

# 从构建层复制二进制
COPY --from=builder /bin/server .
COPY --from=builder /bin/artisan .
# 复制默认配置和迁移文件
COPY configs/ ./configs/
COPY migrations/ ./migrations/

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["./server"]
