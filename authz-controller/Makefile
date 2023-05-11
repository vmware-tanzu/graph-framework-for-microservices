DEBUG ?= FALSE

GO_PROJECT_NAME ?= authz-controller.git
BUCKET_NAME ?= nexus-template-downloads

ECR_DOCKER_REGISTRY ?= 284299419820.dkr.ecr.us-west-2.amazonaws.com/nexus
DOCKER_REGISTRY ?= harbor-repo.vmware.com/nexus
IMAGE_NAME ?= authz
TAG ?= $(shell git rev-parse --verify HEAD)

BUILDER_NAME ?= ${IMAGE_NAME}-builder
BUILDER_TAG := $(shell md5sum builder/Dockerfile | awk '{ print $1 }' | head -c 8)

PKG_NAME ?= /go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/${GO_PROJECT_NAME}

ifeq ($(CONTAINER_ID),)
define run_in_container
  docker run \
  --volume ~/.ssh:/root/.ssh \
  --volume $(realpath .):/go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/authz-controller.git/ \
  --workdir ${PKG_NAME} \
  --env GOPRIVATE="*.eng.vmware.com" \
  --env GOINSECURE=*.eng.vmware.com \
  --env CICD_TOKEN=${CICD_TOKEN} \
  "${BUILDER_NAME}:${BUILDER_TAG}" /bin/bash -c "make cred_setup && ${1}"
endef
else
define run_in_container
 docker run \
 --volumes-from ${CONTAINER_ID} \
 --workdir ${PKG_NAME} \
 --env CICD_TOKEN=${CICD_TOKEN} \
 --env GOPRIVATE=*.eng.vmware.com \
 --env GOINSECURE=*.eng.vmware.com \
 "${BUILDER_NAME}:${BUILDER_TAG}" /bin/bash -c "make cred_setup && ${1}"
endef
endif

%.image.exists:
	@docker inspect $* >/dev/null 2>&1 || \
		(echo "Image $* does not exist. Use 'make docker.builder'." && false)

.PHONY: docker.builder
docker.builder:
	docker build --no-cache -t ${BUILDER_NAME}:${BUILDER_TAG} builder/

.PHONY: build
build:
	cd cmd/authz-controller && \
		CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" .

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
	CGO_ENABLED=0 GOOS=linux go test -cover ./...

.PHONY: race-unit-test
race-unit-test:
	go test -race -cover ./...

.PHONY: test-fmt
test-fmt:
	test -z $$(goimports -w -l cmd pkg)

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test: test-fmt vet lint race-unit-test

.PHONY: test_in_container
test_in_container: ${BUILDER_NAME}\:${BUILDER_TAG}.image.exists
	$(call run_in_container, make test)

.PHONY: show-image-name
show-image-name:
	@echo ${DOCKER_REGISTRY}/${IMAGE_NAME}:${TAG}

.PHONY: docker
docker: build
	docker build --no-cache \
		--build-arg BUILDER_TAG=${BUILDER_TAG} \
		-t ${IMAGE_NAME}:${TAG} .

.PHONY: publish
publish:
	docker tag ${IMAGE_NAME}:${TAG} ${DOCKER_REGISTRY}/${IMAGE_NAME}:${TAG}
	docker push ${DOCKER_REGISTRY}/${IMAGE_NAME}:${TAG};

.PHONY: publish.ecr
publish.ecr:
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

.PHONY: coverage
coverage:
	go test -json -coverprofile=coverage.out -coverpkg=./... ./... | tee report.json ;

build_template:
	tar -czvf authz-controller-manifests.tar manifests/*

publish_template: build_template
	gsutil cp authz-controller-manifests.tar gs://${BUCKET_NAME}/${TAG}/

cred_setup:
	if [ -z ${CICD_TOKEN} ]	;then \
		chmod -R 0600 ~/.ssh/* &&\
        git config --global --add url."git@gitlab.eng.vmware.com:".insteadOf "https://gitlab.eng.vmware.com/" &&\
        go mod tidy && go mod download ;\
	else \
        echo "https://gitlab-ci-token:${CICD_TOKEN}@gitlab.eng.vmware.com" >> ~/.git-credentials && \
        git config --global credential.helper store && \
        go mod tidy && go mod download ;\
    fi

