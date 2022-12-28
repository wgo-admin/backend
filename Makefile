
# 引入依赖
include scripts/make-rules/common.mk

.PHONY: all
all: format build

.PHONY: build
build: deps # 编译源码，依赖 deps 目标自动添加/移除依赖包.
	@go build -v -ldflags "$(GO_LDFLAGS)" -o $(OUTPUT_DIR)/wgo $(ROOT_DIR)/cmd/app/main.go

.PHONY: dev
dev: deps # 开发模式，热更新
	@air

.PHONY: format
format: # 格式化 Go 源码.
	@gofmt -s -w ./

.PHONY: deps
deps: # 自动添加、移除依赖包
	@go mod tidy

.PHONY: clean
clean: # 清除构建产物、临时文件等
	@echo "===========> Cleaning all build output"
	@rm -vrf $(OUTPUT_DIR)

