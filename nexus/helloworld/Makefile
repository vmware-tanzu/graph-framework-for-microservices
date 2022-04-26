BUILD_DIR ?= build
CRD_MODULE_PATH ?= $(shell go list -m)/${BUILD_DIR}/
DATAMODEL_LOCAL_PATH ?= ""
TAG ?= "latest"

.PHONY: datamodel_build
datamodel_build:
	@echo "CRD and API Generated Output Directory: ${BUILD_DIR}"
	@echo "OPENAPISpec Generated Output Directory: ${BUILD_DIR}/crds/"
	mkdir -p ${BUILD_DIR}
	if [ -z $(DATAMODEL_LOCAL_PATH) ]; then \
		docker run\
			--volume $(realpath .):/go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/datamodel/ \
			-v /go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/datamodel/build/ \
			--volume $(realpath .)/build:/go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/generated/ \
			-e CRD_MODULE_PATH=${CRD_MODULE_PATH} \
			-e CONFIG_FILE=datamodel/nexus.yaml \
			--workdir /go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/ \
			harbor-repo.vmware.com/nexus/compiler:$(TAG); \
	else \
		docker run\
			--volumes-from=$(CONTAINER_ID) \
			-e DATAMODEL_LOCAL_PATH=$(DATAMODEL_LOCAL_PATH) \
			-e CRD_MODULE_PATH=${CRD_MODULE_PATH} \
			-e CONFIG_FILE=datamodel/nexus.yaml \
			--workdir /go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/ \
			--user root:root \
			harbor-repo.vmware.com/nexus/compiler:$(TAG); \
	fi
