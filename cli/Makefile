NAME ?= cli
PACKAGE_NAME=$(NAME)
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%S%z)
GIT_VERSION ?= $(shell git describe --always --abbrev=0 --tags)
GIT_REVISION := $(shell git rev-parse --short=8 HEAD || echo unknown)
GIT_BRANCH := $(shell git show-ref | grep "$(GIT_REVISION)" | \
	   grep -v HEAD | awk '{print $$2}' | \
	   sed 's|refs/remotes/origin/||' | \
	   sed 's|refs/heads/||' | sort | head -n 1)

BUILD_VERSION ?= $$(cat BUILD_VERSION)

PKG = github.com/vmware-tanzu/graph-framework-for-microservices/$(PACKAGE_NAME)
COMMON_PACKAGE = $(PKG)/pkg/common


GO_LDFLAGS ?= -X $(COMMON_PACKAGE).VERSION=$(BUILD_VERSION) \
              -X $(COMMON_PACKAGE).BUILT=$(BUILD_DATE) \
              -X $(COMMON_PACKAGE).OS=$(OS) \
			  -X $(COMMON_PACKAGE).GIT_BRANCH=$(GIT_BRANCH) \
              -X $(COMMON_PACKAGE).GIT_COMMIT=$(GIT_VERSION) \
              -s -w
ARTIFACTS_PATH ?= ./artifacts/nexus
build_bin_path := $(shell mkdir -p $(ARTIFACTS_PATH))
UNAME := $(shell uname)

.PHONY: version
version:
	@echo Current version: $(GIT_VERSION)
	@echo Current revision : $(GIT_REVISION)
	@echo Current branch: $(GIT_BRANCH)
	@echo Built: $(BUILD_DATE)

.PHONY: build
build:
	GOOS=darwin GOARCH=amd64 go build --ldflags "$(GO_LDFLAGS)" ./cmd/plugin/nexus
	mv nexus $(ARTIFACTS_PATH)/nexus-darwin_amd64
	GOOS=linux GOARCH=amd64 go build --ldflags "$(GO_LDFLAGS)" ./cmd/plugin/nexus
	mv nexus $(ARTIFACTS_PATH)/nexus-linux_amd64

.PHONY: install
install:
	@if [ "$(UNAME)" = "Darwin" ]; then\
	    cp $(ARTIFACTS_PATH)/nexus-darwin_amd64 /usr/local/bin/nexus;\
		echo "nexus cli installed successfully, please use nexus -h to know more";\
	fi
	@if [ "$(UNAME)" = "Linux" ]; then \
	    cp $(ARTIFACTS_PATH)/nexus-linux_amd64 /usr/local/bin/nexus;\
		echo "nexus installed successfully, please use nexus -h to know more";\
	fi

.PHONY: lint
lint:
	@golangci-lint help linters
	@golangci-lint run ./...

.PHONY: help
help:
	# Commands:
	# make version - show information about current version
	# make build - builds the plugins and stores the targets in artifacts
	# make install - installs the plugins in tanzu framework

print_ldflags:
	@echo $(GO_LDFLAGS)

.PHONY: test
test:
	# TBD

