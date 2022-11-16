DOCKER_REPO ?= harbor-repo.vmware.com/nexus
IMAGE_NAME ?= nexus-graphql-server
IMAGE_TAG ?= $(shell git rev-parse --verify HEAD)
BUILDER_NAME ?= ${IMAGE_NAME}-builder
BUILDER_TAG := $(shell md5sum builder/Dockerfile | awk '{ print $1 }' | head -c 8)
PKG_NAME ?= /go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/graphql-server.git


ifeq ($(CONTAINER_ID),)
define run_in_container
  docker run \
  --volume ~/.ssh:/root/.ssh \
  --volume $(realpath .):/go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/graphql-server.git/ \
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
	go build -trimpath main.go

docker: build_in_container
	docker build -t $(DOCKER_REPO)/$(IMAGE_NAME):$(IMAGE_TAG) . -f Dockerfile ;\

publish: docker
	docker push $(DOCKER_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)

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

build_in_container: ${BUILDER_NAME}\:${BUILDER_TAG}.image.exists
	$(call run_in_container,make build)