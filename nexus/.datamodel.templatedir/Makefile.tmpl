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

# values for nexus datamodel upgrade
FORCE ?= false
DEFAULT_PREV_SPEC_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
PREV_SPEC_BRANCH ?= ""
DEFAULT_PREV_SPEC_DIR ?= ${BUILD_DIR}/crds
PREV_SPEC_DIR ?= ""
ABSOLUTE_PATH_SPEC_DIR ?= $(shell realpath ${PREV_SPEC_DIR})
ABSOLUTE_PATH_BUILD_DIR ?= $(shell realpath .)/${BUILD_DIR}/crds
USER_DEFINED_BRANCH ?= ""
USER_DEFINED_DIRECTORY ?= ""

ifeq (,$(shell echo ${PREV_SPEC_BRANCH}))
USER_DEFINED_BRANCH=false
PREV_SPEC_BRANCH=${DEFAULT_PREV_SPEC_BRANCH}
else
USER_DEFINED_BRANCH=true
endif

ifeq (,$(shell echo ${PREV_SPEC_DIR}))
USER_DEFINED_DIRECTORY=false
PREV_SPEC_DIR=${DEFAULT_PREV_SPEC_DIR}
ABSOLUTE_PATH_SPEC_DIR ?= $(shell realpath ${PREV_SPEC_DIR})
else
USER_DEFINED_DIRECTORY=true
endif

ifeq ($(DOCKER_REPO),)
DOCKER_REPO=$(NAME)
endif

ifeq ($(VERSION),)
VERSION="latest"
endif

.PHONY: datamodel_build
datamodel_build:
	@echo "USER_DEFINED_BRANCH: ${USER_DEFINED_BRANCH} , USER_DEFINED_DIRECTORY: ${USER_DEFINED_DIRECTORY}"
	@echo "CRD and API Generated Output Directory: ${BUILD_DIR}"
	@echo "OPENAPISpec Generated Output Directory: ${BUILD_DIR}/crds/"
	@echo "PREV_SPEC_DIR of Generated Output Directory: ${PREV_SPEC_DIR}"
	@echo "PREV_SPEC_BRANCH of Generated Output Directory: ${PREV_SPEC_BRANCH}"
	mkdir -p ${BUILD_DIR}
	if find "${ABSOLUTE_PATH_SPEC_DIR}" -mindepth 1 -maxdepth 1 ; then \
		cd ${ABSOLUTE_PATH_SPEC_DIR}/.. &&\
		if git ls-tree ${PREV_SPEC_BRANCH} ${ABSOLUTE_PATH_SPEC_DIR} --name-only 2>/dev/null ; then \
			cd ${ABSOLUTE_PATH_SPEC_DIR} && git checkout ${PREV_SPEC_BRANCH} . ;\
		elif  [ ${USER_DEFINED_BRANCH} = true ] ; then \
			echo "Failed to checkout Directory: ${PREV_SPEC_DIR} from ${PREV_SPEC_BRANCH} ";\
			exit 1; \
		fi ;\
		if [ ${ABSOLUTE_PATH_SPEC_DIR} = ${ABSOLUTE_PATH_BUILD_DIR} ] ; then \
			echo "Skipping copy of ${ABSOLUTE_PATH_SPEC_DIR} to ${ABSOLUTE_PATH_BUILD_DIR} directory since both are the same"; \
		else \
			cp -rf ${ABSOLUTE_PATH_SPEC_DIR}/ ${ABSOLUTE_PATH_BUILD_DIR} ;\
			echo "Copying of ${ABSOLUTE_PATH_SPEC_DIR} to the destination directory ${ABSOLUTE_PATH_BUILD_DIR}"; \
		fi ;\
	elif [ ${USER_DEFINED_DIRECTORY} = true ]; then \
		echo "Please provide a valid artifact directory "; \
		exit 1; \
	fi
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
			-e FORCE=$(FORCE) \
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
