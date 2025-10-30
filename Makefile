.PHONY: help build run clean test deps tidy

help: ## 显示帮助信息
	@echo "可用的命令:"
	@echo "  make build    - 编译项目"
	@echo "  make run      - 运行项目"
	@echo "  make clean    - 清理编译文件"
	@echo "  make test     - 运行测试"
	@echo "  make deps     - 下载依赖"
	@echo "  make tidy     - 整理依赖"

build: ## 编译项目
	@echo "编译项目..."
	@mkdir -p bin
	@go build -o bin/api cmd/api/main.go
	@echo "编译完成！可执行文件: bin/api"

run: ## 运行项目
	@echo "运行项目..."
	@go run cmd/api/main.go

clean: ## 清理编译文件
	@echo "清理编译文件..."
	@rm -rf bin
	@echo "清理完成！"

test: ## 运行测试
	@echo "运行测试..."
	@go test -v ./...

deps: ## 下载依赖
	@echo "下载依赖..."
	@go mod download
	@echo "依赖下载完成！"

tidy: ## 整理依赖
	@echo "整理依赖..."
	@go mod tidy
	@echo "依赖整理完成！"

.DEFAULT_GOAL := help
