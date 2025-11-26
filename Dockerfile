# Build Stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# 编译两个二进制：server 和 artisan
RUN go build -o server ./cmd/server && \
    go build -o artisan ./cmd/artisan

# Runtime Stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/artisan .
COPY --from=builder /app/configs ./configs
EXPOSE 8080
CMD ["./server"]
