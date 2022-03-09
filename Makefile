DEBUG ?= FALSE

GO_PROJECT_NAME ?= compiler

ECR_DOCKER_REGISTRY ?= 284299419820.dkr.ecr.us-west-2.amazonaws.com/

IMAGE_NAME ?= nexus-compiler
TAG ?= $(shell git rev-parse --verify --short=8 HEAD)

BUILDER_NAME ?= ${IMAGE_NAME}-builder
BUILDER_TAG := $(shell md5sum builder/Dockerfile | awk '{ print $1 }' | head -c 8)

BASE_NAME ?= ${IMAGE_NAME}-base
BASE_TAG := $(shell md5sum base/Dockerfile | awk '{ print $1 }' | head -c 8)

PKG_NAME?=/go/src/gitlab.eng.vmware.com/nexus/${GO_PROJECT_NAME}

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
 --user $(id -u ${USER}):$(id -g ${USER}) \
 "${BUILDER_NAME}:${BUILDER_TAG}" ${1}
endef
endif

%.image.exists:
	@docker inspect $* >/dev/null 2>&1 || \
		(echo "Image $* does not exist. Use 'make docker.builder' or 'make docker.base'." && false)

.PHONY: docker.builder
docker.builder:
	docker build --no-cache -t ${BUILDER_NAME}:${BUILDER_TAG} builder/

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
	./scripts/update-codegen.sh
	./scripts/generate_openapi_schema.sh
	$(MAKE) -C pkg/openapi_generator generate_test_schemas
	goimports -w pkg

.PHONY: generate_code_in_container
generate_code_in_container: ${BUILDER_NAME}\:${BUILDER_TAG}.image.exists init_submodules
	$(call run_in_container,make generate_code)

.PHONY: show-image-name
show-image-name:
	@echo ${DOCKER_REGISTRY}/${IMAGE_NAME}:${TAG}

.PHONY: docker.base
docker.base:
	docker build --no-cache\
		--build-arg DEBUG=${DEBUG} \
		--tag ${BASE_NAME}:${BASE_TAG} base/

.PHONY: docker
docker: ${BASE_NAME}\:${BASE_TAG}.image.exists
	docker build \
		--build-arg BASE_TAG=${BASE_TAG} \
		--build-arg ASM_SUPERVISOR_TAG=${ASM_SUPERVISOR_TAG} \
		--tag ${IMAGE_NAME}:${TAG} .

.PHONY: publish
publish:
	if gcloud container images describe ${DOCKER_REGISTRY}/${IMAGE_NAME}:${TAG} >/dev/null 2>/dev/null; \
	then \
		echo "Image with the commit tag ${TAG} already exists in the repository!"; \
	else \
		docker tag ${IMAGE_NAME}:${TAG} ${DOCKER_REGISTRY}/${IMAGE_NAME}:${TAG}; \
		docker push ${DOCKER_REGISTRY}/${IMAGE_NAME}:${TAG}; \
	fi

.PHONY: publish.ecr.registry
publish.ecr.registry:
	docker tag ${IMAGE_NAME}:${TAG} ${ECR_DOCKER_REGISTRY}/${IMAGE_NAME}:${TAG}; \
	docker push ${ECR_DOCKER_REGISTRY}/${IMAGE_NAME}:${TAG};

.PHONY: download_builder_image
download_builder_image:
	docker pull ${DOCKER_REGISTRY}/${BUILDER_NAME}:${BUILDER_TAG}
	docker tag ${DOCKER_REGISTRY}/${BUILDER_NAME}:${BUILDER_TAG} ${BUILDER_NAME}:${BUILDER_TAG}

.PHONY: download_base_image
download_base_image:
	docker pull ${DOCKER_REGISTRY}/${BASE_NAME}:${BASE_TAG}
	docker tag ${DOCKER_REGISTRY}/${BASE_NAME}:${BASE_TAG} ${BASE_NAME}:${BASE_TAG}

.PHONY: publish_base_image
publish_base_image:
	docker tag ${BASE_NAME}:${BASE_TAG} ${DOCKER_REGISTRY}/${BASE_NAME}:${BASE_TAG}
	docker push ${DOCKER_REGISTRY}/${BASE_NAME}:${BASE_TAG}

.PHONY: publish_builder_image
publish_builder_image:
	docker tag ${BUILDER_NAME}:${BUILDER_TAG} ${DOCKER_REGISTRY}/${BUILDER_NAME}:${BUILDER_TAG}
	docker push ${DOCKER_REGISTRY}/${BUILDER_NAME}:${BUILDER_TAG}

.PHONY: init_submodules
init_submodules:
	CONTAINER_ID=${CONTAINER_ID} git submodule update --init --recursive
