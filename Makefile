CURRENT_DIR := $(shell pwd)
WORK_DIR_BACKEND := $(CURRENT_DIR)/backend
DIST_DIR_BACKEND := $(WORK_DIR_BACKEND)/dist
WORK_DIR_FRONTEND := $(CURRENT_DIR)/front

VERSION ?= local
COMMIT_ID ?= local

BINARY_NAME := moments
LINUX_AMD64_BINARY_NAME := $(BINARY_NAME)-linux-amd64-$(VERSION)
LINUX_ARM64_BINARY_NAME := $(BINARY_NAME)-linux-arm64-$(VERSION)
LINUX_ARMV7_BINARY_NAME := $(BINARY_NAME)-linux-armv7-$(VERSION)
DARWIN_AMD64_BINARY_NAME := $(BINARY_NAME)-darwin-amd64-$(VERSION)
DARWIN_ARM64_BINARY_NAME := $(BINARY_NAME)-darwin-arm64-$(VERSION)
WINDOWS_AMD64_BINARY_NAME := $(BINARY_NAME)-windows-amd64-$(VERSION).exe
WINDOWS_ARM64_BINARY_NAME := $(BINARY_NAME)-windows-arm64-$(VERSION).exe

BUILD_CMD_DEV := go build -ldflags="-X main.version=$(VERSION) -X main.commitId=$(COMMIT_ID)"
BUILD_CMD_PROD := go build -tags prod -ldflags="-s -w -X main.version=$(VERSION) -X main.commitId=$(COMMIT_ID)"

.PHONY: frontend-install frontend-dev backend-dev clean build frontend backend zip checksums

frontend-install:
	cd $(WORK_DIR_FRONTEND) && pnpm i

backend-install:
	cd $(WORK_DIR_BACKEND) && go mod download

frontend-dev:
	cd $(WORK_DIR_FRONTEND) && pnpm run dev

backend-dev:
	cd $(WORK_DIR_BACKEND) && $(BUILD_CMD_DEV) -o $(DIST_DIR_BACKEND)/$(LINUX_AMD64_BINARY_NAME)
	$(DIST_DIR_BACKEND)/$(LINUX_AMD64_BINARY_NAME)

build: clean frontend-install frontend backend-install backend

clean:
	cd $(WORK_DIR_BACKEND) && go clean
	cd $(WORK_DIR_BACKEND) && rm -rf ./public ./dist
	cd $(WORK_DIR_FRONTEND) && rm -rf ./.output ./dist

frontend:
	cd $(WORK_DIR_FRONTEND) && pnpm generate
	cp -r $(WORK_DIR_FRONTEND)/.output/public $(WORK_DIR_BACKEND)

backend:
	cd $(WORK_DIR_BACKEND) && GOOS=linux GOARCH=amd64 $(BUILD_CMD_PROD) -o $(DIST_DIR_BACKEND)/$(LINUX_AMD64_BINARY_NAME)
	cd $(WORK_DIR_BACKEND) && GOOS=linux GOARCH=arm64 $(BUILD_CMD_PROD) -o $(DIST_DIR_BACKEND)/$(LINUX_ARM64_BINARY_NAME)
	cd $(WORK_DIR_BACKEND) && GOOS=linux GOARCH=arm GOARM=7 $(BUILD_CMD_PROD) -o $(DIST_DIR_BACKEND)/$(LINUX_ARMV7_BINARY_NAME)
	cd $(WORK_DIR_BACKEND) && GOOS=darwin GOARCH=amd64 $(BUILD_CMD_PROD) -o $(DIST_DIR_BACKEND)/$(DARWIN_AMD64_BINARY_NAME)
	cd $(WORK_DIR_BACKEND) && GOOS=darwin GOARCH=arm64 $(BUILD_CMD_PROD) -o $(DIST_DIR_BACKEND)/$(DARWIN_ARM64_BINARY_NAME)
	cd $(WORK_DIR_BACKEND) && GOOS=windows GOARCH=amd64 $(BUILD_CMD_PROD) -o $(DIST_DIR_BACKEND)/$(WINDOWS_AMD64_BINARY_NAME)
	cd $(WORK_DIR_BACKEND) && GOOS=windows GOARCH=arm64 $(BUILD_CMD_PROD) -o $(DIST_DIR_BACKEND)/$(WINDOWS_ARM64_BINARY_NAME)

zip:
	cd $(DIST_DIR_BACKEND) && find . -type f -name "$(BINARY_NAME)-*" | xargs -I {} zip {}.zip {}

checksums:
	cd $(DIST_DIR_BACKEND) && md5sum $(LINUX_AMD64_BINARY_NAME) > $(LINUX_AMD64_BINARY_NAME)-checksum.txt
	cd $(DIST_DIR_BACKEND) && md5sum $(LINUX_ARM64_BINARY_NAME) > $(LINUX_ARM64_BINARY_NAME)-checksum.txt
	cd $(DIST_DIR_BACKEND) && md5sum $(LINUX_ARMV7_BINARY_NAME) > $(LINUX_ARMV7_BINARY_NAME)-checksum.txt
	cd $(DIST_DIR_BACKEND) && md5sum $(DARWIN_AMD64_BINARY_NAME) > $(DARWIN_AMD64_BINARY_NAME)-checksum.txt
	cd $(DIST_DIR_BACKEND) && md5sum $(DARWIN_ARM64_BINARY_NAME) > $(DARWIN_ARM64_BINARY_NAME)-checksum.txt
	cd $(DIST_DIR_BACKEND) && md5sum $(WINDOWS_AMD64_BINARY_NAME) > $(WINDOWS_AMD64_BINARY_NAME)-checksum.txt
	cd $(DIST_DIR_BACKEND) && md5sum $(WINDOWS_ARM64_BINARY_NAME) > $(WINDOWS_ARM64_BINARY_NAME)-checksum.txt
