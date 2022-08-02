SHELL := /bin/bash

DEBUG ?= FALSE

GO_PROJECT_NAME ?= compiler.git

DOCKER_REGISTRY ?= harbor-repo.vmware.com

IMAGE_NAME ?= nexus/compiler
TAG ?= $(shell git rev-parse --verify --short=8 HEAD)

BUILDER_NAME ?= ${IMAGE_NAME}-builder
BUILDER_TAG := $(shell md5sum builder/Dockerfile | awk '{ print $1 }' | head -c 8)

PKG_NAME?=/go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/${GO_PROJECT_NAME}

DATAMODEL_PATH ?= datamodel
CONFIG_FILE ?= ""
GENERATED_OUTPUT_DIRECTORY ?= generated

ifeq ($(CONTAINER_ID),)
define run_in_container
  docker run \
  --volume $(realpath .):${PKG_NAME} \
  --volume ~/.ssh:/root/.ssh \
  --network=host \
  --workdir ${PKG_NAME} \
  "${BUILDER_NAME}:${BUILDER_TAG}" /bin/bash -c "make docker.gitlab_credentials && ${1}"
endef
else
define run_in_container
 docker run \
 --volumes-from ${CONTAINER_ID} \
 --workdir ${PKG_NAME} \
 --env CICD_TOKEN=${CICD_TOKEN} \
 "${BUILDER_NAME}:${BUILDER_TAG}" /bin/bash -c "make docker.gitlab_credentials && ${1}"
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
	go install github.com/onsi/ginkgo/ginkgo@v1.16.0
	go install github.com/onsi/gomega/...@v1.17.0
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/mikefarah/yq/v4@latest
	go install gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kube-openapi.git/cmd/nexus-openapi-gen@latest

.PHONY: unit-test
unit-test:
	ginkgo -cover ./pkg/...
	cd example/tests && ginkgo -cover ./...

.PHONY: race-unit-test
race-unit-test:
	ginkgo -race -cover ./pkg/...
	cd example/tests && ginkgo -race -cover ./...

.PHONY: test-fmt
test-fmt:
	test -z $$(goimports -w -l cmd pkg)

.PHONY: vet
vet:
	go vet ./cmd/... ./pkg/...

.PHONY: lint
lint:
	golangci-lint run ./cmd/... ./pkg/...

coverage:
	go test -json -coverprofile=coverage.out -coverpkg=./... ./... | tee report.json ;

.PHONY: test
test: test-fmt vet lint race-unit-test

.PHONY: test_in_container
test_in_container: ${BUILDER_NAME}\:${BUILDER_TAG}.image.exists
	$(call run_in_container, make test)

.PHONY: generate_code
generate_code:
	echo "Cleaning up workdir"
	rm -rf _generated ${GOPATH}/nexustempmodule
	cp -R generated_base_structure _generated
	cp ${DATAMODEL_PATH}/go.mod _generated/go.mod
	sed -i "1s|.*|module nexustempmodule|" _generated/go.mod
	echo "Generating base nexus code structure"
	CRD_MODULE_PATH=${CRD_MODULE_PATH} go run cmd/nexus-sdk/main.go -config-file ${CONFIG_FILE} -dsl ${DATAMODEL_PATH} -crd-output _generated
	mv _generated/api_names.sh scripts/
	echo "Resolving datamodel dependencies"
	cd _generated && ../scripts/pin_deps.sh  && go mod tidy -e
	echo "Generating kuberenetes APIs"
	./scripts/generate_k8s_api.sh
	echo "Generating openapi schema"
	go install gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kube-openapi.git/cmd/nexus-openapi-gen@latest
	cd _generated/ && go mod tidy -e
	./scripts/generate_openapi_schema.sh
	echo "Generating CRD yamls"
	go run cmd/generate-openapischema/generate-openapischema.go -yamls-path _generated/crds
	git checkout -- pkg/openapi_generator/openapi/openapi_generated.go
	echo "Updating module name"
	./scripts/replace_mod_path.sh
	echo "Sorting imports"
	cd _generated && goimports -w .
	echo "Moving files to output directory"
	cp -r _generated/{client,apis,crds,common,nexus-client,helper} ${GENERATED_OUTPUT_DIRECTORY}

.PHONY: test_generate_code_in_container
test_generate_code_in_container: ${BUILDER_NAME}\:${BUILDER_TAG}.image.exists init_submodules
	$(call run_in_container, make generate_code DATAMODEL_PATH=example/datamodel \
	CONFIG_FILE=example/nexus-sdk.yaml \
	CRD_MODULE_PATH="gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/crd_generated/" \
	GENERATED_OUTPUT_DIRECTORY=example/output/crd_generated && \
	cd example/output/crd_generated && go mod tidy && go vet ./... && golangci-lint run ./...)
	@if [ -n "$$(git ls-files --modified --exclude-standard)" ]; then\
		echo "The following changes should be committed:";\
		git status;\
		git diff;\
		exit 1;\
	fi

.PHONY: generate_example
generate_example:
	$(MAKE) generate_code DATAMODEL_PATH=example/datamodel CONFIG_FILE=example/nexus-sdk.yaml CRD_MODULE_PATH="gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/crd_generated/" GENERATED_OUTPUT_DIRECTORY=example/output/crd_generated

.PHONY: show-image-name
show-image-name:
	@echo ${DOCKER_REGISTRY}/${IMAGE_NAME}:${TAG}

.PHONY: publish
publish:
	docker tag ${IMAGE_NAME}:${TAG} ${DOCKER_REGISTRY}/${IMAGE_NAME}:${TAG}
	docker push ${DOCKER_REGISTRY}/${IMAGE_NAME}:${TAG};

.PHONY: download_builder_image
download_builder_image:
	docker pull ${DOCKER_REGISTRY}/${BUILDER_NAME}:${BUILDER_TAG}
	docker tag ${DOCKER_REGISTRY}/${BUILDER_NAME}:${BUILDER_TAG} ${BUILDER_NAME}:${BUILDER_TAG}

.PHONY: publish_builder_image
publish_builder_image:
	docker tag ${BUILDER_NAME}:${BUILDER_TAG} ${DOCKER_REGISTRY}/${BUILDER_NAME}:${BUILDER_TAG}
	docker push ${DOCKER_REGISTRY}/${BUILDER_NAME}:${BUILDER_TAG}

.PHONY: init_submodules
init_submodules:
	CONTAINER_ID=${CONTAINER_ID} git submodule update --init --recursive

.PHONY: render_templates
render_templates:
	go run cmd/nexus-sdk/main.go -config-file example/nexus-sdk.yaml -dsl example/datamodel -crd-output example/output/_crd_base

.PHONY: test_render_templates
test_render_templates: render_templates
	@if [ -n "$$(git ls-files --modified --exclude-standard)" ]; then\
    	echo "The following changes should be committed:";\
    	git status;\
    	git diff;\
    	return 1;\
    fi

.PHONY: docker.gitlab_credentials
docker.gitlab_credentials:
	echo "Updating gitlab settings"
	@if test -z $$CICD_TOKEN; then\
		echo "Using ssh authentication" && \
        git config --global --add url."git@gitlab.eng.vmware.com:".insteadOf "https://gitlab.eng.vmware.com/"; \
	else \
		echo "Using https authentication" && \
        git config --global credential.helper store && \
        echo -e  "https://gitlab-ci-token:${CICD_TOKEN}@gitlab.eng.vmware.com/" >> ~/.git-credentials && \
        git config --global url."https://gitlab.eng.vmware.com/".insteadOf "git@gitlab.eng.vmware.com:"; \
	fi
