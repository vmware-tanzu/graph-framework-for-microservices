DEBUG ?= FALSE

GO_PROJECT_NAME ?= compiler

ECR_DOCKER_REGISTRY ?= 284299419820.dkr.ecr.us-west-2.amazonaws.com

IMAGE_NAME ?= nexus/compiler
TAG ?= $(shell git rev-parse --verify --short=8 HEAD)

BUILDER_NAME ?= ${IMAGE_NAME}-builder
BUILDER_TAG := $(shell md5sum builder/Dockerfile | awk '{ print $1 }' | head -c 8)

PKG_NAME?=/go/src/gitlab.eng.vmware.com/nexus/${GO_PROJECT_NAME}

DATAMODEL_PATH ?= datamodel
CONFIG_FILE ?= ""
GENERATED_OUTPUT_DIRECTORY ?= generated
DATAMODEL_LOCAL_PATH ?= ""

ifeq ($(CONTAINER_ID),)
define run_in_container
  docker run \
  --volume $(realpath .):/go/src/gitlab.eng.vmware.com/nexus/compiler/ \
  --workdir ${PKG_NAME} \
  "${BUILDER_NAME}:${BUILDER_TAG}" ${1}
endef
else
define run_in_container
 docker run \
 --volumes-from ${CONTAINER_ID} \
 --volume ~/.ssh:/root/.ssh \
 --workdir ${PKG_NAME} \
 "${BUILDER_NAME}:${BUILDER_TAG}" ${1}
endef
endif

%.image.exists:
	@docker inspect $* >/dev/null 2>&1 || \
		(echo "Image $* does not exist. Use 'make docker.builder'." && false)

.PHONY: docker.builder
docker.builder:
	docker build --no-cache -t ${BUILDER_NAME}:${BUILDER_TAG} builder/

.PHONY: docker
docker: init_submodules ${BUILDER_NAME}\:${BUILDER_TAG}.image.exists
	git archive -o compiler.tar --format=tar HEAD
	tar -rf compiler.tar .git
	docker build --no-cache \
		--build-arg BUILDER_TAG=${BUILDER_TAG} \
		-t ${IMAGE_NAME}:${TAG} .

.PHONY: build
build:
	cd cmd/nexus-sdk && \
		CGO_ENABLED=0 go build -ldflags="-w -s" .

.PHONY: build_in_container
build_in_container: ${BUILDER_NAME}\:${BUILDER_TAG}.image.exists
	$(call run_in_container,make build)

.PHONY: tools
tools:
	go get golang.org/x/tools/cmd/goimports
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega/...

.PHONY: unit-test
unit-test:
	ginkgo -cover ./...

.PHONY: race-unit-test
race-unit-test:
	ginkgo -race -cover ./...

.PHONY: test-fmt
test-fmt:
	test -z $$(goimports -w -l cmd pkg)

.PHONY: vet
vet:
	go vet ./cmd/... ./pkg/...

.PHONY: lint
lint:
	golangci-lint run ./cmd/... ./pkg/...

.PHONY: test
test: test-fmt vet lint race-unit-test

.PHONY: test_in_container
test_in_container: ${BUILDER_NAME}\:${BUILDER_TAG}.image.exists
	$(call run_in_container, make test)

.PHONY: generate_code
generate_code:
	if [ -n $(DATAMODEL_LOCAL_PATH) ]; then \
		cp -rf $(DATAMODEL_LOCAL_PATH) /go/src/gitlab.eng.vmware.com/nexus/compiler/datamodel && \
		rm -rf /go/src/gitlab.eng.vmware.com/nexus/compiler/datamodel/build  ;\
	fi
	go run cmd/nexus-sdk/main.go -config-file ${CONFIG_FILE} -dsl ${DATAMODEL_PATH} -crd-output _generated | tee test
	cat test
	mv _generated/api_names.sh ./scripts/
	./scripts/generate_k8s_api.sh
	./scripts/generate_openapi_schema.sh
	$(MAKE) -C pkg/openapi_generator generate_test_schemas
	goimports -w pkg
	if [ -n $(DATAMODEL_LOCAL_PATH) ]; then \
		cp -r _generated/{client,apis,crds} $(DATAMODEL_LOCAL_PATH)/build ;\
	else \
		cp -r _generated/{client,apis,crds} ${GENERATED_OUTPUT_DIRECTORY} ;\
	fi

.PHONY: test_generate_code_in_container
test_generate_code_in_container: ${BUILDER_NAME}\:${BUILDER_TAG}.image.exists init_submodules
	$(call run_in_container, make generate_code DATAMODEL_PATH=example/datamodel \
	CONFIG_FILE=example/nexus-sdk.yaml \
	GENERATED_OUTPUT_DIRECTORY=example/output/_crd_generated)
	@if [ -n "$$(git ls-files --modified --exclude-standard)" ]; then\
		echo "The following changes should be committed:";\
		git status;\
		git diff;\
		return 1;\
	fi

.PHONY: show-image-name
show-image-name:
	@echo ${ECR_DOCKER_REGISTRY}/${IMAGE_NAME}:${TAG}

.PHONY: publish
publish:
	docker tag ${IMAGE_NAME}:${TAG} ${ECR_DOCKER_REGISTRY}/${IMAGE_NAME}:${TAG}
	docker push ${ECR_DOCKER_REGISTRY}/${IMAGE_NAME}:${TAG};

.PHONY: download_builder_image
download_builder_image:
	docker pull ${ECR_DOCKER_REGISTRY}/${BUILDER_NAME}:${BUILDER_TAG}
	docker tag ${ECR_DOCKER_REGISTRY}/${BUILDER_NAME}:${BUILDER_TAG} ${BUILDER_NAME}:${BUILDER_TAG}

.PHONY: publish_builder_image
publish_builder_image:
	docker tag ${BUILDER_NAME}:${BUILDER_TAG} ${ECR_DOCKER_REGISTRY}/${BUILDER_NAME}:${BUILDER_TAG}
	docker push ${ECR_DOCKER_REGISTRY}/${BUILDER_NAME}:${BUILDER_TAG}

.PHONY: init_submodules
init_submodules:
	CONTAINER_ID=${CONTAINER_ID} git submodule update --init --recursive
