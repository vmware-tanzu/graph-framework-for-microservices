BUILD_DIR ?= build
CRD_MODULE_PATH ?= $(shell go list -m)/${BUILD_DIR}/
TAG ?= "latest"
CONTAINER_ID ?= ""
DATAMODEL_LOCAL_PATH ?= $(realpath .)
BUCKET ?= nexus-template-downloads
TAG ?= $(shell git rev-parse --verify HEAD)

.PHONY: datamodel_build
datamodel_build:
	@echo "CRD and API Generated Output Directory: ${BUILD_DIR}"
	@echo "OPENAPISpec Generated Output Directory: ${BUILD_DIR}/crds/"
	rm -rf ${DATAMODEL_LOCAL_PATH}/build
	mkdir -p ${BUILD_DIR}
	if [ -z $(CONTAINER_ID) ]; then \
		docker run --pull=always\
			--volume $(realpath .):/go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/datamodel/ \
			-v /go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/datamodel/build/ \
			--volume $(realpath .)/build:/go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/generated/ \
			--volume ~/.ssh:/root/.ssh \
			-e CRD_MODULE_PATH=${CRD_MODULE_PATH} \
			-e CONFIG_FILE=datamodel/nexus.yaml \
			-e GOPRIVATE=gitlab.eng.vmware.com \
			harbor-repo.vmware.com/nexus/compiler:$(TAG); \
	else \
		docker run --pull=always\
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

archive_crds:
	cd ${BUILD_DIR}/crds/ &&\
		rm -f api-datamodel-crds.tar && \
		tar -czvf api-datamodel-crds.tar *;

publish: archive_crds
	cd ${BUILD_DIR}/crds && \
	gsutil cp api-datamodel-crds.tar gs://${BUCKET}/${TAG}/;
