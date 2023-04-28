BUILD_DIR ?= build
CRD_MODULE_PATH ?= $(shell go list -m)/${BUILD_DIR}/
TAG ?= "latest"
CONTAINER_ID ?= ""
DATAMODEL_LOCAL_PATH ?= $(realpath .)
BUCKET ?= nexus-template-downloads
DOCKER_REPO ?= $(shell cat nexus.yaml | grep dockerRepo | awk '{print $$2}' | tr -d '"')
VERSION ?= $(shell git rev-parse --verify HEAD 2>/dev/null)
NAME ?= $(shell cat nexus.yaml | grep groupName | awk '{print $$2}' | tr -d '"')
NAMESPACE ?= ""
FORCE ?= false

# values for nexus datamodel upgrade
FORCE ?= true

ifeq (,$(shell echo ${ARTIFACT_REPO}))
ARTIFACT_REPO = $(shell cat nexus.yaml | grep artifactRepo | awk '{print $$2}' | tr -d '"')
endif

ifneq (,$(shell echo ${PREV_SPEC_BRANCH}))
# backward compatability check is required.
FORCE=false

# check if ARTIFACT_REPO is configured
ifeq (,$(shell echo ${ARTIFACT_REPO}))
# Error and exit because user provided the branch and artifact repo details are not available
$(error Unable to determine ARTIFACT_REPO to validate backward compatibility with ${PREV_SPEC_BRANCH} branch)
endif

TARGET_REPO_NAME = $(shell echo ${ARTIFACT_REPO} | awk -F'/' '{print $$2}' | tr -d '"'| awk -F'.git' '{print $$1}')

ifeq (,$(shell echo $(shell realpath ${TARGET_REPO_NAME}/crds)))
ARTIFACT_REPO_GITCMD=$(shell echo git clone $(shell echo ${ARTIFACT_REPO}) --branch $(shell echo ${PREV_SPEC_BRANCH}) --single-branch $(shell echo ${TARGET_REPO_NAME}))
else
ARTIFACT_REPO_GITCMD=$(shell echo git -C $(shell realpath ${TARGET_REPO_NAME}) pull origin $(shell echo ${PREV_SPEC_BRANCH}))
endif

endif

ifeq ($(DOCKER_REPO),)
DOCKER_REPO=$(NAME)
endif

ifeq ($(VERSION),)
VERSION="latest"
endif

.PHONY: datamodel_build
datamodel_build:
	mkdir -p ${BUILD_DIR}
	@echo "FORCE: ${FORCE}"
	@echo "PREV_SPEC_BRANCH: ${PREV_SPEC_BRANCH}"
	@echo "ARTIFACT_REPO: ${ARTIFACT_REPO}"
	@echo "TARGET_REPO_NAME: ${TARGET_REPO_NAME}"
	@echo "ARTIFACT_REPO_GITCMD: ${ARTIFACT_REPO_GITCMD}"

	if [ ${FORCE} = false ] ; then \
        		bash -c "${ARTIFACT_REPO_GITCMD}";\
          		cp -rf ${TARGET_REPO_NAME}/crds ${BUILD_DIR} ;\
            	echo "Copying of ${TARGET_REPO_NAME}/crds to the destination directory ${BUILD_DIR}/crds"; \
    fi ;\

	if [ -z $(CONTAINER_ID) ]; then \
		docker run --pull=missing\
			--volume $(realpath .):/go/src/github.com/vmware-tanzu/graph-framework-for-microservices/compiler/datamodel/ \
			-v /go/src/github.com/vmware-tanzu/graph-framework-for-microservices/compiler/datamodel/build/ \
			--volume $(realpath .)/build:/go/src/github.com/vmware-tanzu/graph-framework-for-microservices/compiler/generated/ \
			--volume ~/.ssh:/root/.ssh \
			-e CRD_MODULE_PATH=${CRD_MODULE_PATH} \
			-e CONFIG_FILE=datamodel/nexus.yaml \
			-e GOPRIVATE=*.eng.vmware.com \
			-e GOINSECURE=*.eng.vmware.com \
			-e FORCE=$(FORCE) \
			gcr.io/nsx-sm/nexus/compiler:$(TAG) /bin/bash -c 'make generate_code && chown -R $(shell id -u):$(shell id -g) .;'; \
	else \
		docker run --pull=missing\
			--volumes-from=$(CONTAINER_ID) \
			-e DATAMODEL_PATH=$(DATAMODEL_LOCAL_PATH) \
			-e GENERATED_OUTPUT_DIRECTORY=$(DATAMODEL_LOCAL_PATH)/build \
			-e CONFIG_FILE=${DATAMODEL_LOCAL_PATH}/nexus.yaml \
			-e CRD_MODULE_PATH=${CRD_MODULE_PATH} \
			-e GOPRIVATE=*.eng.vmware.com \
			-e GOINSECURE=*.eng.vmware.com \
			-e CICD_TOKEN=${CICD_TOKEN} \
			--user root:root \
			-e FORCE=${FORCE} \
			gcr.io/nsx-sm/nexus/compiler:$(TAG); \
	fi

docker_build:
	@if [ -n "$(NAME)" ]; then \
				echo "Building image: $(DOCKER_REPO):$(VERSION)" ;\
				docker build --build-arg IMAGE_NAME=$(DOCKER_REPO):$(VERSION) --build-arg NAME=$(NAME) -t $(DOCKER_REPO):$(VERSION) . -f Dockerfile ; \
	else \
			echo "Please provide NAME (GroupName) when running docker_build" ;\
			exit 1;\
	fi

docker_publish: docker_build
	docker push $(DOCKER_REPO):$(VERSION) ;
