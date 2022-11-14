DEBUG ?= FALSE

GO_PROJECT_NAME ?= connector.git
BUCKET_NAME ?= nexus-template-downloads
CHART_NAME ?= "nexus-connector"
HELM_REGISTRY ?= oci://284299419820.dkr.ecr.us-west-2.amazonaws.com/nexus
ECR_DOCKER_REGISTRY ?= 284299419820.dkr.ecr.us-west-2.amazonaws.com/nexus
DOCKER_REGISTRY ?= gcr.io/nsx-sm/nexus
IMAGE_NAME ?= connector
TAG ?= $(shell git rev-parse --verify HEAD)
VERSION ?= "v0.0.0-$(TAG)"
HARBOR_REPO_URL ?= "https://harbor-repo.vmware.com/chartrepo/nexus"
HARBOR_REPO ?= "harbor-vmware"

BUILDER_NAME ?= ${IMAGE_NAME}-builder
BUILDER_TAG := $(shell md5sum builder/Dockerfile | awk '{ print $1 }' | head -c 8)

PKG_NAME ?= /go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/${GO_PROJECT_NAME}

ifeq ($(CONTAINER_ID),)
define run_in_container
  docker run \
  --volume ~/.ssh:/root/.ssh \
  --volume $(realpath .):/go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/connector.git/ \
  --workdir ${PKG_NAME} \
  --env GOPRIVATE="*.eng.vmware.com" \
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
	cd cmd/connector && \
		CGO_ENABLED=0 GOOS=linux go build .

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
	set -o pipefail && CGO_ENABLED=0 go test -v -tags=unit -p=1 -count=1 -vet=off ./...

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
docker: build_in_container
	docker build --no-cache \
		--build-arg BUILDER_TAG=${BUILDER_TAG} \
		-t ${IMAGE_NAME}:${TAG} .

.PHONY: publish
publish:
	docker tag ${IMAGE_NAME}:${TAG} ${DOCKER_REGISTRY}/${IMAGE_NAME}:${TAG}
	docker push ${DOCKER_REGISTRY}/${IMAGE_NAME}:${TAG};


.PHONY: download_builder_image
download_builder_image:
	docker pull ${ECR_DOCKER_REGISTRY}/${BUILDER_NAME}:${BUILDER_TAG}
	docker tag ${ECR_DOCKER_REGISTRY}/${BUILDER_NAME}:${BUILDER_TAG} ${BUILDER_NAME}:${BUILDER_TAG}

.PHONY: publish_builder_image
publish_builder_image:
	docker tag ${BUILDER_NAME}:${BUILDER_TAG} ${ECR_DOCKER_REGISTRY}/${BUILDER_NAME}:${BUILDER_TAG}
	docker push ${ECR_DOCKER_REGISTRY}/${BUILDER_NAME}:${BUILDER_TAG}

build_template:
	tar -czvf connector-manifests.tar manifests/*

publish_template: build_template
	gsutil cp connector-manifests.tar gs://${BUCKET_NAME}/${TAG}/

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

coverage:
	go test -json -coverprofile=coverage.out -coverpkg=./... ./... | tee report.json ;


submodule:
	git submodule update --init --recursive
build_helm: submodule
	mkdir -p nexus-connector/crds
	cp -rf api/build/crds/* nexus-connector/crds/
	sed "s|__CONNECTOR_TAG___|$(TAG)|g" ./values.yaml > nexus-connector/values.yaml
	helm package $(CHART_NAME) --version $(VERSION)

publish.ecr.helm: build_helm
	helm push $(CHART_NAME)-$(VERSION).tgz $(HELM_REGISTRY)

harbor.login:
	helm repo remove $(HARBOR_REPO) || echo "$(HARBOR_REPO) not present"; \
	helm repo add $(HARBOR_REPO) $(HARBOR_REPO_URL) --username $(HARBOR_USERNAME) --password $(HARBOR_PASSWORD);

publish.harbor.helm: build_helm
	helm cm-push $(CHART_NAME)-$(VERSION).tgz $(HARBOR_REPO)