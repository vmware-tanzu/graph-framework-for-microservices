#
# App Info
#
APP_NAME ?= sampleapp
NAMESPACE ?= default

#
# Image Info
#
### adding this to test app init..
CI_COMMIT ?= $(shell git rev-parse --verify --short=8 HEAD 2> /dev/null || echo "00000000")
IMAGE_TAG ?= ${APP_NAME}-${CI_COMMIT}
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
	nexus datamodel $(filter-out $@,$(MAKECMDGOALS))

.PHONY: runtime
runtime: ## Invoke runtime operations. Run with 'help' for options.
	nexus runtime $(filter-out $@,$(MAKECMDGOALS))

.PHONY: app
app: ## Invoke app operations. Run with 'help' for options.
	nexus app $(filter-out $@,$(MAKECMDGOALS))


##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

##@lint checks using the make targets.
.PHONY: lint
lint:
	$(MAKE) fmt
	$(MAKE) vet

.PHONY: test
test:
	go test ./...

##@ Space for adding integration test and related cleanup functions...
.PHONY: integration_test
integration_test:
	echo "Add your integration_tests for your app here!!!!..." ;
	echo "Sample Workflow can be.";
	echo "1. Create Environment";
	echo "2. Start built application in the Environment";
	echo "3. Start integration_tests with go test / gingko framework";

.PHONY: integration_test_cleanup
integration_test_cleanup:
	echo "Add your cleanup steps here!!!!....";
	echo "Possible steps you can do.";
	echo "1. Get logs of integration test as artifacts"
	echo "2. Get logs of components in clusters as artifacts for debugging"

.PHONY: teardown_environment
teardown_environment:
	echo "Add cluster cleanup step after integration_tests pass/fail here..";
	echo "Clear clusters created";

##@ Coverage checks using sonar-scanner
.PHONY: coverage
coverage:
    go test -json -coverprofile=coverage.out ./... | tee report.json
	APP_NAME=${APP_NAME} sonar-scanner

##@ Build
.PHONY: build
build: fmt vet ## Build manager binary.
	go mod download ;
	env GOOS=linux GOARCH=amd64 go build -o bin/${APP_NAME} main.go ;

.PHONY: docker-build
docker-build: ## Build docker image with the manager.
	docker build --build-arg APP_NAME=${APP_NAME} -t ${IMAGE_REGISTRY}:${IMAGE_TAG} .

.PHONY: publish
publish: build ## Push docker image with the manager.
	docker push ${IMAGE_REGISTRY}:${IMAGE_TAG}

.PHONY: image_scan
image_scan:
	flash docker scan image ${IMAGE_REGISTRY}:${IMAGE_TAG}

##@ Deployment
.PHONY: deploy
deploy:
	sed -e 's/__APP_NAME__/'"${APP_NAME}"'/g' -e 's|__IMAGE__|'"${IMAGE_REGISTRY}:${IMAGE_TAG}"'|g' config/deployment/deployment.yaml | kubectl apply -f - -n ${NAMESPACE};


.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	sed -e 's/__APP_NAME__/'"${APP_NAME}"'/g' -e 's|__IMAGE__|'"${IMAGE_REGISTRY}:${IMAGE_TAG}"'|g' config/deployment/deployment.yaml | kubectl delete --ignore-not-found=true -f - -n ${NAMESPACE};

#check how to use kustomize for now using sed to replace deployment..
#KUSTOMIZE = $(shell pwd)/bin/kustomize
#.PHONY: kustomize
#kustomize: ## Download kustomize locally if necessary.
#$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

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
