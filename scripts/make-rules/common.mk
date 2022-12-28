# ==============================================================================
# 定义全局 Makefile 变量方便后面引用
SHELL := /bin/bash

COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
# 项目根目录
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/../../ && pwd -P))
# 构建产物目录
OUTPUT_DIR := $(ROOT_DIR)/_output
# 定义包名
ROOT_PACKAGE=github.com/wgo-admin/backend

# ==============================================================================
# 定义版本相关变量

# 使用指定的 version 包通过 `-ldflags -X` 注入我们的版本信息
VERSION_PACKAGE=$(ROOT_PACKAGE)/pkg/version

# 获取版本号信息
ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --tags --always --match='v*')
endif

# 检查代码仓库是否dirty，默认是dirty
GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
  GIT_TREE_STATE="clean"
endif
# 获取构建时的 Commit ID
GIT_COMMIT:=$(shell git rev-parse HEAD)

# -ldflags
GO_LDFLAGS += \
  -X $(VERSION_PACKAGE).GitVersion=$(VERSION) \
  -X $(VERSION_PACKAGE).BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') \
  -X $(VERSION_PACKAGE).GitCommit=$(GIT_COMMIT) \
  -X $(VERSION_PACKAGE).GitTreeState=$(GIT_TREE_STATE)
