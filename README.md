.PHONY: run build artisan-controller lint test

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go
	go build -o bin/artisan cmd/artisan/main.go

# 模拟 Laravel 命令体验
# make controller name=User
controller:
	go run cmd/artisan/main.go make:controller $(name)

lint:
	golangci-lint run ./...

test:
	go test -v ./...
