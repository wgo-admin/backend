# 默认执行 all 目标
.DEFAULT_GOAL := all

.PHONY: all
all: go.format go.build

# ==============================================================================
# Includes
include scripts/make-rules/common.mk
include scripts/make-rules/golang.mk


## --------------------------------------
## Cleanup
## --------------------------------------

##@ clean:

.PHONY: clean
clean: ## 清理构建产物、临时文件等.
	@echo "===========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)

## --------------------------------------
## Binaries
## --------------------------------------

##@ build:

.PHONY: build
build: go.tidy ## 编译源码，依赖 deps 目标自动添加/移除依赖包.
	@$(MAKE) go.build

.PHONY: dev
dev: go.tidy ## 开发模式，热更新
	@air

## --------------------------------------
## Hack / Tools
## --------------------------------------

##@ hack/tools:

.PHONY: tidy
tidy: ## 自动添加、移除依赖包
	@$(MAKE) go.tidy

.PHONY: format
format:  ## 格式化 Go 源码.
	@$(MAKE) go.format

.PHONY: help
help: Makefile ## 打印 Makefile help 信息.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<TARGETS> <OPTIONS>\033[0m\n\n\033[35mTargets:\033[0m\n"} /^[0-9A-Za-z._-]+:.*?##/ { printf "  \033[36m%-45s\033[0m %s\n", $$1, $$2 } /^\$$\([0-9A-Za-z_-]+\):.*?##/ { gsub("_","-", $$1); printf "  \033[36m%-45s\033[0m %s\n", tolower(substr($$1, 3, length($$1)-7)), $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' Makefile #$(MAKEFILE_LIST)
	@echo -e "$$USAGE_OPTIONS"
