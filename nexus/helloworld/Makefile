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

ifeq ($(DOCKER_REPO),)
DOCKER_REPO=$(NAME)
endif

ifeq ($(VERSION),)
VERSION="latest"
endif

.PHONY: datamodel_build
datamodel_build:
	@echo "CRD and API Generated Output Directory: ${BUILD_DIR}"
	@echo "OPENAPISpec Generated Output Directory: ${BUILD_DIR}/crds/"
	rm -rf ${DATAMODEL_LOCAL_PATH}/build
	mkdir -p ${BUILD_DIR}
	if [ -z $(CONTAINER_ID) ]; then \
		docker run --pull=missing\
			--volume $(realpath .):/go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/datamodel/ \
			-v /go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/datamodel/build/ \
			--volume $(realpath .)/build:/go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/generated/ \
			--volume ~/.ssh:/root/.ssh \
			-e CRD_MODULE_PATH=${CRD_MODULE_PATH} \
			-e CONFIG_FILE=datamodel/nexus.yaml \
			-e GOPRIVATE=gitlab.eng.vmware.com \
			harbor-repo.vmware.com/nexus/compiler:$(TAG); \
	else \
		docker run --pull=missing\
			--volumes-from=$(CONTAINER_ID) \
			-e DATAMODEL_PATH=$(DATAMODEL_LOCAL_PATH) \
			-e GENERATED_OUTPUT_DIRECTORY=$(DATAMODEL_LOCAL_PATH)/build \
			-e CONFIG_FILE=${DATAMODEL_LOCAL_PATH}/nexus.yaml \
			-e CRD_MODULE_PATH=${CRD_MODULE_PATH} \
			-e GOPRIVATE=gitlab.eng.vmware.com \
			-e CICD_TOKEN=${CICD_TOKEN} \
			--user root:root \
			harbor-repo.vmware.com/nexus/compiler:$(TAG); \
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
