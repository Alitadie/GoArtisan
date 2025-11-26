APP_NAME = server
VERSION ?= $(shell git describe --tags --always --dirty || echo "dev")
COMMIT_HASH ?= $(shell git rev-parse --short HEAD || echo "none")
BUILD_TIME ?= $(shell date +%FT%T%z)

# Linker flags: -X 动态修改 pkg/version 包里的变量
LDFLAGS := -X 'go-artisan/pkg/version.GitTag=${VERSION}' \
           -X 'go-artisan/pkg/version.GitCommit=${COMMIT_HASH}' \
           -X 'go-artisan/pkg/version.BuildTime=${BUILD_TIME}'


.PHONY: run build

run:
	go run -ldflags "${LDFLAGS}" cmd/server/main.go

build:
	@echo "📦 Building ${VERSION}..."
	@mkdir -p bin
	@go build -ldflags "-s -w ${LDFLAGS}" -o bin/${APP_NAME} cmd/server/main.go
	@echo "✅ Build success: bin/${APP_NAME}"

# 模拟 Laravel 命令体验
# make controller name=User
controller:
	go run cmd/artisan/main.go make:controller $(name)

lint:
	golangci-lint run ./...

test:
	go test -v ./...
