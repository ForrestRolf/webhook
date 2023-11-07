SHELL:=/bin/sh
.PHONY: build clean

export GO111MODULE=on

MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))
RELEASE_DIR := ${MKFILE_DIR}dist
GO_PATH := $(shell go env | grep GOPATH | awk -F '"' '{print $$2}')
ENABLE_CGO= CGO_ENABLED=0

build: clean \
		build_webhook_ui \
		build_webhook_linux_amd64 \
		build_webhook_windows_amd64 \
		build_webhook_darwin_amd64 \
		build_webhook_darwin_arm64

build_webhook_linux_amd64:
	@echo "build webhook - linux amd64"
	cd ${MKFILE_DIR} && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -trimpath \
	-o ${RELEASE_DIR}/webhook-linux-amd64 ${MKFILE_DIR}

build_webhook_windows_amd64:
	@echo "build webhook - windows amd64"
	cd ${MKFILE_DIR} && \
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -v -trimpath \
	-o ${RELEASE_DIR}/webhook-win-amd64 ${MKFILE_DIR}

build_webhook_darwin_amd64:
	@echo "build webhook - darwin amd64"
	cd ${MKFILE_DIR} && \
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -v -trimpath \
	-o ${RELEASE_DIR}/webhook-darwin-amd64 ${MKFILE_DIR}

build_webhook_darwin_arm64:
	@echo "build webhook - darwin arm64"
	cd ${MKFILE_DIR} && \
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -v -trimpath \
	-o ${RELEASE_DIR}/webhook-darwin-arm64 ${MKFILE_DIR}

build_webhook_ui:
	@echo "build webhook - UI"
	cd web && npm run build

clean:
	rm -rf ${RELEASE_DIR}
