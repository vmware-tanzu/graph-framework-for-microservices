SHELL := /bin/bash
#
# App Info
#
APP_NAME ?= api-gw
NAMESPACE ?= default
DATAMODEL_DIR ?= "./nexus"
CLUSTER ?= ""
DATAMODEL ?= ""
DATAMODEL_GROUP ?= ""
NEXUS_BIN ?= $(shell which nexus)
BUCKET ?= nexus-template-downloads
PROJECT_DIR ?= $(shell pwd)
#
# Image Info
#
### adding this to test app init..
CI_COMMIT ?= $(shell git rev-parse --verify --short=8 HEAD 2> /dev/null || echo "00000000")
IMAGE_TAG ?= ${CI_COMMIT}
GIT_HEAD  ?= $(shell git rev-parse --verify HEAD 2> /dev/null || echo "0000000000000000")
IMAGE_REGISTRY ?= gcr.io/nsx-sm/nexus/nexus-api-gateway
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

replace:
	if [ -n ${DATAMODEL} ]; then\
			test -s ${DATAMODEL_DIR}/${DATAMODEL} || { echo "Please create datamodel ${DATAMODEL} for go mod replace"; exit 1; } ;\
			go mod edit -replace $(DATAMODEL)=${DATAMODEL_DIR}/${DATAMODEL} ;\
	fi

.SILENT:
.PHONY: datamodel_init
datamodel_init: ## Initialize datamodel
	if [ -z $(NEXUS_BIN) ]; then \
		echo "Please install nexus binary"; \
		exit 1; \
	fi
	if [ -n ${DATAMODEL} ]; then \
		if [ -n ${DATAMODEL_GROUP} ]; then \
			$(NEXUS_BIN) datamodel init --name ${DATAMODEL} --group ${DATAMODEL_GROUP};\
		else \
			$(NEXUS_BIN) datamodel init --name ${DATAMODEL} ;\
		fi \
	else \
		$(NEXUS_BIN) datamodel init ;\
	fi
	$(MAKE) replace

##@ Dev

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: lint
lint: ## lint checks using the make targets
	$(MAKE) fmt
	$(MAKE) vet

go_get:
	GOINSECURE=*.eng.vmware.com GOPRIVATE=*.eng.vmware.com go get . ;

.PHONY: build
build: lint ## Build manager binary.
	mkdir -p .ssh ;\
	if [ -n $(CICD_TOKEN) ]; then \
		DOCKER_BUILDKIT=1 docker build --build-arg APP_NAME=${APP_NAME} \
					--build-arg GIT_HEAD=${GIT_HEAD} \
					--build-arg GIT_TAG=${CI_COMMIT} \
					--build-arg CICD_TOKEN=${CICD_TOKEN} \
					-t ${IMAGE_REGISTRY}:${IMAGE_TAG} --progress=plain --no-cache -f Dockerfile.build;\
	else \
		test -s ~/.ssh || { echo "Please provide CICD_TOKEN if ssh not available."; exit 1; }; \
		cp -rf ~/.ssh/* .ssh ;\
		DOCKER_BUILDKIT=1 docker build --build-arg APP_NAME=${APP_NAME} \
					--build-arg GIT_HEAD=${GIT_HEAD} \
					--build-arg GIT_TAG=${CI_COMMIT} \
					--build-arg USE_SSH="true" \
					-t ${IMAGE_REGISTRY}:${IMAGE_TAG} --progress=plain --no-cache -f Dockerfile.build;\
	fi

##@ Test
.PHONY: init-unit-test
init-unit-test:
	if [[ ! -f test/bin/etcd || ! -f test/bin/kube-apiserver || ! -f test/bin/kubectl ]]; then \
		curl -SLo envtest-bins.tar.gz "https://go.kubebuilder.io/test-tools/1.24.2/$$(go env GOOS)/$$(go env GOARCH)"; \
        tar --strip-components=1 -xvf envtest-bins.tar.gz -C test; \
        rm envtest-bins.tar.gz; \
  	fi

.PHONY: unit-test
unit-test: init-unit-test
	ginkgo -cover ./controllers/...
	ginkgo -cover ./pkg/...
	ginkgo -cover ./internal/...

.PHONY: race-unit-test
race-unit-test: init-unit-test
	ginkgo -cover ./controllers/...
	ginkgo -race -cover ./pkg/... ./internal/...

.PHONY: test-fmt
test-fmt:
	test -z $$(goimports -w -l cmd pkg)

.PHONY: test
test: race-unit-test

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

##@ Coverage
.PHONY: coverage
coverage: init-unit-test
	go test --json -coverprofile=coverage.out -coverpkg=./... ./...   | tee report.json ; exit $${PIPESTATUS[0]} ;\


##@ Publish

.PHONY: publish
publish: build ## Push docker image with the manager.
	docker push ${IMAGE_REGISTRY}:${IMAGE_TAG}

##@ Security Scan

.PHONY: image_scan
image_scan:
	flash docker scan image ${IMAGE_REGISTRY}:${IMAGE_TAG}

add_operator: install-nexus-kubebuilder
	if [ -n $(CRD_DATAMODEL_NAME) ]; then \
		if [ -n $(CRD_GROUP) ]; then \
			if [ -n $(CRD_VERSION) ]; then \
				if [ -n $(CRD_KIND) ]; then \
					$(NEXUS-KUBEBUILDER) create api --group $(CRD_GROUP) --kind $(CRD_KIND) --version $(CRD_VERSION) --controller --resource=false --import $(CRD_DATAMODEL_NAME) ;\
				else \
					echo "Please provide CRD_KIND"; exit 1;\
				fi \
			else \
				echo "Please provide CRD_VERSION"; exit 1; \
			fi \
		else \
			echo "Please provide CRD_GROUP"; exit 1; \
		fi \
	else \
		echo "Please provide CRD_DATAMODEL_NAME"; exit 1; \
	fi

NEXUS-KUBEBUILDER = $(PROJECT_DIR)/bin/nexus-kubebuilder
install-nexus-kubebuilder:
	test -s ${PROJECT_DIR}/bin/nexus-kubebuilder || { mkdir -p ${PROJECT_DIR}/bin; cd ${PROJECT_DIR}/bin; GOBIN=${PROJECT_DIR}/bin go install gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kubebuilder.git/cmd/nexus-kubebuilder@master ; }


build_manifests:
	cd manifests && \
	rm -rf api-gw-manifests.tar && \
	cd ..&& \
	tar -czvf api-gw-manifests.tar manifests/*

publish_manifests:
	gsutil cp api-gw-manifests.tar gs://${BUCKET}/${IMAGE_TAG}/;

.PHONY: deploy_kind
deploy_kind:
	kind load docker-image --name ${KIND_NAME} ${IMAGE_REGISTRY}:${IMAGE_TAG}
	kubectl -n ${NAMESPACE} set image deployment/nexus-api-gw ${APP_NAME}=${IMAGE_REGISTRY}:${IMAGE_TAG}

CURRENT_TIMESTAMP = $(shell date +%s)
.PHONY: build_deploy_kind
build_deploy_kind:
	make build IMAGE_TAG=${CURRENT_TIMESTAMP} && make deploy_kind NAMESPACE=${NAMESPACE} KIND_NAME=kind IMAGE_TAG=${CURRENT_TIMESTAMP}

