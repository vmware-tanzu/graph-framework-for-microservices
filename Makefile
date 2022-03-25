#
# App Info
#
APP_NAME ?= sampleapp

#
# Image Info
#
GIT_TAG ?= $(shell git rev-parse --verify --short=8 HEAD)
IMAGE_TAG ?= ${APP_NAME}-${GIT_TAG}
IMAGE_REGISTRY ?= 284299419820.dkr.ecr.us-west-2.amazonaws.com/nexus/playground

#
# Platform
#
NEXUS_CLI_TAG ?= latest

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Platform
.PHONY: nexus-cli
nexus-cli: ## Install Nexus CLI
	curl https://storage.googleapis.com/nexus-cli-downloads/nexus-$(uname | tr '[:upper:]' '[:lower:]')_amd64 --output nexus
	chmod 755 nexus
	mv nexus /usr/local/bin/nexus


.PHONY: datamodel
datamodel: ## Invoke datamodel operations. Run with 'help' for options.
	nexus datamodel ${ARGS}

.PHONY: runtime
runtime: ## Invoke runtime operations. Run with 'help' for options.
	nexus runtime ${ARGS}

.PHONY: app
app: ## Invoke app operations. Run with 'help' for options.
	nexus app ${ARGS}


##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

##@ Build

.PHONY: build
build: fmt vet ## Build manager binary.
	go build -o bin/${APP_NAME} main.go

.PHONY: run
run: fmt vet ## Run a controller from your host.
	go run ./main.go

.PHONY: docker-build
docker-build: ## Build docker image with the manager.
	docker build -t ${IMAGE_TAG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push ${IMAGE_REGISTRY}:${IMAGE_TAG}

##@ Deployment

.PHONY: deploy
deploy: kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/${APP_NAME} && $(KUSTOMIZE) edit set image controller=${IMAGE_TAG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/default | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

KUSTOMIZE = $(shell pwd)/bin/kustomize
.PHONY: kustomize
kustomize: ## Download kustomize locally if necessary.
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef
