DEBUG ?= FALSE

GO_PROJECT_NAME ?= validation.git

ECR_DOCKER_REGISTRY ?= 284299419820.dkr.ecr.us-west-2.amazonaws.com

IMAGE_NAME ?= nexus-validation
#TAG ?= $(shell git rev-parse --verify --short=8 HEAD)
TAG = latest

BUILDER_NAME ?= ${IMAGE_NAME}-builder
BUILDER_TAG := $(shell md5sum builder/Dockerfile | awk '{ print $1 }' | head -c 8)

PKG_NAME?=/go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/${GO_PROJECT_NAME}

ifeq ($(CONTAINER_ID),)
define run_in_container
  docker run \
  --volume $(realpath .):/go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/validation.git/ \
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

.PHONY: build
build:
	cd cmd/nexus-validation && \
		CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" .

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
docker: build_in_container ${BUILDER_NAME}\:${BUILDER_TAG}.image.exists
	docker build --no-cache \
		--build-arg BUILDER_TAG=${BUILDER_TAG} \
		-t ${IMAGE_NAME}:${TAG} .

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

.PHONY: create_kind_cluster
create_kind_cluster:
	kind create cluster --config manifests/kind.cluster.yaml

.PHONY: kind_load_image
kind_load_image: docker
	kind load docker-image nexus-validation:latest

.PHONY: kind_deploy_config
kind_deploy_config:
	kubectl apply -f manifests/validating.config.yaml

.PHONY: kind_delete_config
kind_delete_config:
	kubectl delete -f manifests/validating.config.yaml

.PHONY: kind_delete_deploy
kind_delete_deploy:
	kubectl delete -f manifests/webhook.deploy.yaml || true
	kubectl delete -f manifests/webhook.svc.yaml || true
	kubectl delete -f manifests/webhook.tls.yaml || true
	kubectl delete -f manifests/rbac.yaml || true

.PHONY: kind_deploy
kind_deploy: kind_load_image kind_delete_deploy kind_deploy_config
	kubectl apply -f manifests/webhook.svc.yaml
	kubectl apply -f manifests/webhook.tls.yaml
	kubectl apply -f manifests/webhook.deploy.yaml
	kubectl apply -f manifests/rbac.yaml


.PHONY: kind_logs
kind_logs:
	kubectl logs -l app=nexus-validation -f
