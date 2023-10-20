DOCKER_GO_PATH ?= /go
DOCKER_BUILD_MOUNT_DIR ?= ${DOCKER_GO_PATH}/src/github.com/vmware-tanzu/graph-framework-for-microservices
API_GW_COMPONENT_NAME ?= api-gw
CLI_DIR ?= cli
HOST_KUBECONFIG ?= ${HOME}/.kube/config
MOUNTED_KUBECONFIG ?= /etc/config/kubeconfig
RUNTIME_NAMESPACE ?= default
UNAME_S ?= $(shell uname -s)
DATAMODEL_NAME ?= $(shell grep groupName ${DATAMODEL_DIR}/nexus.yaml | cut -f 2 -d" " | tr -d '"')
DATAMODEL_IMAGE_NAME ?= $(shell grep dockerRepo ${DATAMODEL_DIR}/nexus.yaml | cut -f 2 -d" ")
TAG ?= $(shell cat TAG | awk '{ print $1 }')

.PHONY: cli.build.darwin
cli.build.darwin:
	docker run \
		--pull=missing \
		--volume $(realpath .):${DOCKER_BUILD_MOUNT_DIR} \
		-w ${DOCKER_BUILD_MOUNT_DIR}/${CLI_DIR} \
		golang:1.19.8 \
		/bin/bash -c "go mod download && make build.darwin";

.PHONY: cli.build.linux
cli.build.linux:
	docker run \
		--pull=missing \
		--volume $(realpath .):${DOCKER_BUILD_MOUNT_DIR} \
		-w ${DOCKER_BUILD_MOUNT_DIR}/${CLI_DIR} \
		golang:1.19.8 \
		/bin/bash -c "go mod download && make build.linux";

.PHONY: cli.install.darwin
cli.install.darwin:
	cd cli; make install.darwin

.PHONY: cli.install.linux
cli.install.linux:
	cd cli; make install.linux

.PHONY: compiler.build
compiler.build:
	cd compiler; BUILDER_TAG=${TAG} make docker.builder
	cd compiler; BUILDER_TAG=${TAG} TAG=${TAG} make docker

.PHONY: api-gw.docker
api-gw.docker:
	docker run \
		--pull=missing \
		--volume $(realpath .):${DOCKER_BUILD_MOUNT_DIR} \
		-w ${DOCKER_BUILD_MOUNT_DIR}/${API_GW_COMPONENT_NAME} \
		golang:1.19.8 \
		/bin/bash -c "go mod download && GOOS=linux GOARCH=amd64 go build -buildvcs=false -o bin/${API_GW_COMPONENT_NAME}";
	docker build -t ${API_GW_COMPONENT_NAME}:${TAG} -f api-gw/Dockerfile .

.PHONY: k0s.install
k0s.install:
	$(realpath .)/nexus-runtime-manifests/k0s/run_k0s.sh

.PHONY: k0s.uninstall
k0s.uninstall:
	$(realpath .)/nexus-runtime-manifests/k0s/stop_k0s.sh

.PHONY: dm.install_init
dm.install_init: HOST_KUBECONFIG=$(realpath .)/nexus-runtime-manifests/k0s/.kubeconfig
dm.install_init:
	docker run \
		--net host \
		--pull=missing \
		--volume $(realpath .):${DOCKER_BUILD_MOUNT_DIR} \
		-w ${DOCKER_BUILD_MOUNT_DIR}/nexus-runtime-manifests/datamodel-install \
		--mount type=bind,source=${HOST_KUBECONFIG},target=${MOUNTED_KUBECONFIG},readonly \
		alpine/helm \
		upgrade --install datamodel-install-scripts . --set global.namespace=${RUNTIME_NAMESPACE} --kubeconfig ${MOUNTED_KUBECONFIG}

.PHONY: api.build
api.build:
	cd api; TAG=${TAG} VERSION=${TAG} make datamodel_build

.PHONY: api.install
api.install:
	docker run \
		--entrypoint /datamodel_installer.sh \
		--net host \
		--pull=missing \
		--volume $(realpath .)/api/build/crds:/crds \
		--mount type=bind,source=${HOST_KUBECONFIG},target=${MOUNTED_KUBECONFIG},readonly \
		--mount type=bind,source=$(realpath .)/nexus-runtime-manifests/datamodel-install/datamodel_installer.sh,target=/datamodel_installer.sh,readonly \
		-e KUBECONFIG=${MOUNTED_KUBECONFIG} \
		-e NAME=nexus.vmware.com \
		-e IMAGE=gcr.io/nsx-sm/nexus/nexus-api \
		bitnami/kubectl

.PHONY: api-gw.run
api-gw.run: HOST_KUBECONFIG=$(realpath .)/nexus-runtime-manifests/k0s/.kubeconfig
api-gw.run: api-gw.stop
ifeq (${UNAME_S}, Linux)
	docker run -d \
		--name=nexus-api-gw \
		--rm \
                --network host \
		--pull=missing \
		--mount type=bind,source=${HOST_KUBECONFIG},target=${MOUNTED_KUBECONFIG},readonly \
		--mount type=bind,source=$(realpath .)/${API_GW_COMPONENT_NAME}/deploy/config/api-gw-config.yaml,target=/api-gw-config.yaml,readonly \
		-e APIGWCONFIG=/api-gw-config.yaml \
		-e KUBECONFIG=${MOUNTED_KUBECONFIG} \
		-e KUBEAPI_ENDPOINT="127.0.0.1:8001" \
		${API_GW_COMPONENT_NAME}:${TAG}
else
	APIGWCONFIG=$(realpath .)/${API_GW_COMPONENT_NAME}/deploy/config/api-gw-config.yaml KUBECONFIG=${HOST_KUBECONFIG} $(realpath .)/${API_GW_COMPONENT_NAME}/bin/${API_GW_COMPONENT_NAME}
endif

.PHONY: api-gw.stop
api-gw.stop:
	docker rm -f nexus-api-gw > /dev/null || true

.PHONY: runtime.build
runtime.build: compiler.build api.build api-gw.docker

.PHONY: clean.runtime
clean.runtime:
	rm -rf api/build/*

.PHONY: runtime.install.k0s 
runtime.install.k0s: HOST_KUBECONFIG=$(realpath .)/nexus-runtime-manifests/k0s/.kubeconfig
runtime.install.k0s: k0s.install dm.install_init api.install api-gw.run
	$(info )
	$(info ====================================================)
	$(info To access runtime, you can execute kubectl as:)
	$(info     kubectl -s localhost:8082 ...)
	$(info )
	$(info )
	$(info To access nexus aip gateway using kubeconfig, export:)
	$(info     export HOST_KUBECONFIG=${HOST_KUBECONFIG})
	$(info )
	$(info ====================================================)

.PHONY: runtime.uninstall.k0s
runtime.uninstall.k0s: k0s.uninstall
	$(info )
	$(info ====================================================)
	$(info Runtime is now uninstalled)
	$(info ====================================================)

.PHONY: dm.check-datamodel-dir
dm.check-datamodel-dir:
ifndef DATAMODEL_DIR
	$(error DATAMODEL_DIR is mandatory)
endif

.PHONY: dm.install
dm.install: dm.check-datamodel-dir
	docker run \
		--net host \
		--pull=missing \
		--volume ${DATAMODEL_DIR}/build/crds:/crds \
		--mount type=bind,source=${HOST_KUBECONFIG},target=${MOUNTED_KUBECONFIG},readonly \
		--mount type=bind,source=$(realpath .)/nexus-runtime-manifests/datamodel-install/datamodel_installer.sh,target=/datamodel_installer.sh,readonly \
		--entrypoint /datamodel_installer.sh \
		-e KUBECONFIG=${MOUNTED_KUBECONFIG} \
		-e NAME=${DATAMODEL_NAME} \
		-e IMAGE=${DATAMODEL_IMAGE_NAME} \
		bitnami/kubectl

