SHELL := /bin/bash

DEBUG ?= FALSE
LOG_LEVEL ?= error

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

NEXUS_KUBEOPENAPI_VERSION ?= 7416bd4754d3c0dd8b3fa37fff53d36594f11607
NEXUS_GRAPHQLGEN_VERSION ?= 32f028bce22efeb70b47a640195bd969dbb337f0

ifeq ($(CONTAINER_ID),)
define run_in_container
  docker run \
  --volume $(realpath .):${PKG_NAME} \
  --volume ~/.ssh:/root/.ssh \
  --network=host \
  --workdir ${PKG_NAME} \
  "${BUILDER_NAME}:${BUILDER_TAG}" /bin/bash -c "make docker.gitlab_credentials && \
  git config --global --add safe.directory '*' && ${1}"
endef
else
define run_in_container
 docker run \
 --volumes-from ${CONTAINER_ID} \
 --workdir ${PKG_NAME} \
 --env CICD_TOKEN=${CICD_TOKEN} \
 --env PKG_NAME=${PKG_NAME} \
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
docker: init_submodules ${BUILDER_NAME}\:${BUILDER_TAG}.image.exists build_openapigen_in_container build_gqlgen_in_container
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
	go install github.com/onsi/gomega/...@v1.18.0
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/mikefarah/yq/v4@latest
	go install gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kube-openapi.git/cmd/nexus-openapi-gen@${NEXUS_KUBEOPENAPI_VERSION}

.PHONY: build_openapigen_in_container
build_openapigen_in_container:
	$(call run_in_container,make build_openapigen)

.PHONY: build_openapigen
build_openapigen:
	GOBIN=${PKG_NAME}/cmd go install gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kube-openapi.git/cmd/nexus-openapi-gen@${NEXUS_KUBEOPENAPI_VERSION}

build_gqlgen_in_container:
	$(call run_in_container,make build_gqlgen)

build_gqlgen:
	GOBIN=${PKG_NAME}/cmd go install gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git@${NEXUS_GRAPHQLGEN_VERSION}

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
	go test -json -coverprofile=coverage.out -coverpkg=./... ./pkg/... | tee report.json ;

.PHONY: test
test: test-fmt vet lint race-unit-test

.PHONY: test_in_container
test_in_container: ${BUILDER_NAME}\:${BUILDER_TAG}.image.exists
	$(call run_in_container, make test)

.PHONY: generate_code
generate_code:
	@echo "Nexus Compiler: Running compiler code generation"
	@echo "Cleaning up workdir"
	rm -rf _generated ${GOPATH}/src/nexustempmodule
	@echo "Copying generated_base_structure to create directory structure"
	cp -R generated_base_structure _generated
	@echo "Copying go.mod file of datamodel"
	cp ${DATAMODEL_PATH}/go.mod _generated/go.mod
	sed -i'.bak' -e "1s|.*|module nexustempmodule|" _generated/go.mod
	cd _generated/ && go mod edit -go=1.18
	@echo "Nexus Compiler: Generating base nexus code structure"
	CRD_MODULE_PATH=${CRD_MODULE_PATH} go run cmd/nexus-sdk/main.go -config-file ${CONFIG_FILE} -dsl ${DATAMODEL_PATH} -crd-output _generated -log-level ${LOG_LEVEL}
	mv _generated/api_names.sh scripts/
	@echo "Nexus Compiler: Resolving datamodel dependencies"
	cd _generated && ../scripts/pin_deps.sh  && go mod tidy -e
	@echo "Nexus Compiler: Generating kuberenetes APIs"
	./scripts/generate_k8s_api.sh
	@echo "Nexus Compiler: Generating openapi schema"
	./scripts/generate_openapi_schema.sh
	@echo "Nexus Compiler: Generating CRD yamls"
	go run cmd/generate-openapischema/generate-openapischema.go -yamls-path _generated/crds
	git checkout -- pkg/openapi_generator/openapi/openapi_generated.go
	@echo "==> Nexus Compiler: Generating GRAPHQL pkg <=="
	cd _generated && goimports -w .
	cd _generated/nexus-gql && go get github.com/cespare/xxhash/v2@v2.1.2 && go get golang-appnet.eng.vmware.com/nexus-sdk/nexus@NPT-482-add-server-client && gqlgen generate
	cp -rf _generated/* ${GOPATH}/src/nexustempmodule/
	cd ${GOPATH}/src/nexustempmodule && cd nexus-gql && CGO_ENABLED=1 GOOS=linux \
	go build --trimpath -o graphql.so -buildmode=plugin server.go
	@echo "Updating module name"
	./scripts/replace_mod_path.sh
	find . -name "*.bak" -type f -delete
	@echo "Sorting imports"
	cd _generated && goimports -w .
	@echo "Nexus Compiler: Moving files to output directory"
	cp -r _generated/{client,apis,crds,common,nexus-client,helper,nexus-gql} ${GENERATED_OUTPUT_DIRECTORY}
	cp -r ${GOPATH}/src/nexustempmodule/nexus-gql/graphql.so ${GENERATED_OUTPUT_DIRECTORY}/nexus-gql
	@echo "Nexus Compiler: Compiler code generation completed"

.PHONY: test_generate_code_in_container
test_generate_code_in_container: ${BUILDER_NAME}\:${BUILDER_TAG}.image.exists init_submodules
	$(call run_in_container, go install gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kube-openapi.git/cmd/nexus-openapi-gen@${NEXUS_KUBEOPENAPI_VERSION} && \
	go install gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git@${NEXUS_GRAPHQLGEN_VERSION} && \
	mv /go/bin/gqlgen.git /go/bin/gqlgen && \
	make generate_code DATAMODEL_PATH=example/datamodel \
	CONFIG_FILE=example/nexus-sdk.yaml \
	CRD_MODULE_PATH="gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/crd_generated/" \
	LOG_LEVEL=trace \
	GENERATED_OUTPUT_DIRECTORY=example/output/crd_generated && \
	cd example/output/crd_generated && go mod tidy && go vet -structtag=FALSE ./... && golangci-lint run ./...)
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
