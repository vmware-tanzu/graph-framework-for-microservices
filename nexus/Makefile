UNAME := $(shell uname)
NAMESPACE ?= default
BUILD_DIR ?= build
DATAMODEL ?= "datamodel"
OUTPUT_DIR ?= "${BUILD_DIR}/${DATAMODEL}/output"
DATAMODEL_DIR ?= "./nexus"
APP_DIR ?= ".."
PID=$(shell ps -ef | grep $$(cat /tmp/portforward.pid) | grep -v grep | wc -l| tr -d ' ')
FILE ?= ""
SHELL=bash

.PHONY: runtime_install
runtime_install:
	kubectl create namespace ${NAMESPACE} --dry-run -o yaml | kubectl apply -f -
	mkdir -p /tmp/ && \
	cd /tmp  && \
	curl https://storage.googleapis.com/nexus-template-downloads/runtime-manifests.tar -o manifests.tar && \
	tar -xvf manifests.tar && \
	kubectl apply -f runtime-manifests/crds/ -n ${NAMESPACE} &&\
	kubectl apply -f runtime-manifests/deployment/ -n ${NAMESPACE};
	$(MAKE) post-install-checks

.SILENT:
post-install-checks:
	CURRENT=0
	while [[ -z `kubectl get pods -n ${NAMESPACE} -lapp=etcd 2>/dev/null` ]] ; do \
		CURRENT=$$(( CURRENT + 1 )) ;\
		TRIES_LEFT=$$(( 20 - CURRENT )) ;\
		if [[ $$TRIES_LEFT -ne 0 ]] ; then \
			sleep 6 ;\
		else \
			echo "etcd pod not created"; \
			exit 1 ; \
		fi \
	done ; \
	kubectl wait pods -lapp=etcd -n ${NAMESPACE} --for=condition=ready --timeout=300s || { "etcd pod could not be started"; exit 1; }
	kubectl wait pods -lapp=kube-apiserver -n ${NAMESPACE} --for=condition=ready --timeout=120s || { "k8sapiserver pod could not be started"; exit 1; }
	kubectl wait pods -lname=kube-controllermanager -n ${NAMESPACE} --for=condition=ready --timeout=120s || { "k8scontrollermanager pod could not be started"; exit 1; }

.SILENT:
start-portforward: check-portforward
	( kubectl port-forward svc/nexus-proxy-container -n ${NAMESPACE} 45192:80 & echo $$! > /tmp/portforward.pid )
	CURRENT=0
	while [[ `curl -s -o /dev/null -w "%{http_code}" http://localhost:45192/api/v1/namespaces 2>/dev/null` -ne 200  ]] ; do \
		CURRENT=$$(( CURRENT + 1 )) ;\
		TRIES_LEFT=$$(( 10 - CURRENT )) ;\
		if [[ $$TRIES_LEFT -ne 0 ]] ; then \
			sleep 5 ;\
		else \
			echo "portforwarding failed due to service not reachable"; \
			$(MAKE) kill-portforward ;\
			exit 1 ; \
		fi \
	done

.SILENT:
check-portforward:
	@if [ -z `which telnet` ]; then \
		echo "please install telnet for mac : brew install telnet" ;\
		exit 1 ;\
	fi
	@if [ -z `kubectl get ns ${NAMESPACE} --no-headers 2>/dev/null` ]; then \
		echo "Please install runtime on namespace" ;\
		exit 1 ;\
	fi
	@if echo -e '\x1dclose\x0d' | telnet localhost 45192; then\
		echo "could not portforward to 45192 port - port already open"; exit 1; \
	else \
		echo "Proceeding to start portforwarding.." ;\
	fi

.SILENT:
kill-portforward:
	@if test -s /tmp/portforward.pid ; then \
		{ kill $(shell cat /tmp/portforward.pid); rm /tmp/portforward.pid ;} ;\
	fi

check-file:
	@if [ "$(FILE)" = "" ]; then \
		{ echo 'please call as FILE="<filename> make datamodel_apply'; exit 1; }  ;\
	fi

.PHONY: datamodel_install
datamodel_install: start-portforward
	kubectl -s localhost:45192 apply -f ${DATAMODEL}/build/crds/ || { echo "could not install datamodel"; kill $(shell cat /tmp/portforward.pid); rm /tmp/portforward.pid; exit 1; } && \
	kill $(shell cat /tmp/portforward.pid); rm /tmp/portforward.pid;	

.PHONY: datamodel_apply
datamodel_apply: check-file start-portforward
	kubectl -s localhost:45192 apply -f $(FILE) -n default || { echo "could not apply datamodel crds"; kill $(shell cat /tmp/portforward.pid); rm /tmp/portforward.pid; exit 1; } && \
	{ kill $(shell cat /tmp/portforward.pid); rm /tmp/portforward.pid; } ;\
	

.PHONY: datamodel_build
datamodel_build:
	@echo "CRD and API Generated Output Directory: ${OUTPUT_DIR}"
	@echo "OPENAPISpec Generated Output Directory: ${BUILD_DIR}/${DATAMODEL}/policymodel/crds/controller"
	mkdir -p ${BUILD_DIR}/${DATAMODEL}
	@@test -s ${BUILD_DIR}/${DATAMODEL}/policymodel || git clone git@gitlab.eng.vmware.com:nsx-allspark_users/m7/policymodel.git -b nexus-sdk-dev --single-branch --depth 1 ${BUILD_DIR}/${DATAMODEL}/policymodel
	@if [ "$(UNAME)" = "Darwin" ]; then\
		cp bin/nexus-sdk-darwin-x64 bin/nexus-sdk; \
	fi
	@if [ "$(UNAME)" = "Linux" ]; then \
		cp bin/nexus-sdk-linux-x64 bin/nexus-sdk;\
	fi
	@@test -s bin/nexus-sdk || { echo "Please use mac/linux system for building datamodel"; exit 1; }
	echo "starting to build datamodel crds and api for openapispec generation."
	@@bin/nexus-sdk -config-file ${DATAMODEL}/nexus-sdk.yaml -crd-output ${OUTPUT_DIR}  -dsl ${DATAMODEL} || { echo "Failure in generating crds and apis" ; exit 1;}
	@@mkdir -p ${BUILD_DIR}/${DATAMODEL}/policymodel/pkg/apis
	@@mkdir -p ${BUILD_DIR}/${DATAMODEL}/policymodel/crds/controller
	echo "Copying files to policymodel for openapispec generation"
	@@cp -rf ${OUTPUT_DIR}/*/*.yaml ${BUILD_DIR}/${DATAMODEL}/policymodel/crds/controller/
	@@cp -rf ${OUTPUT_DIR}/api_names.sh ${BUILD_DIR}/${DATAMODEL}/policymodel/scripts
	@@cp -rf ${OUTPUT_DIR}/* ${BUILD_DIR}/${DATAMODEL}/policymodel/pkg/apis
	cd	${BUILD_DIR}/${DATAMODEL}/policymodel && make generate_code_in_container || { echo "Failure in generation of openapispec."; exit 1; }
	@@rm -rf ${OUTPUT_DIR};


.PHONY: app_init
app_init:
	if [[ -z ${APP_DIR} ]]; then \
		echo "could not initialize app please provide app directory using APP_DIR variable"; \
		exit 1 ; \
	fi
	cd ${APP_DIR} && \
		echo "Intializing go mod" && \
		test -s go.mod || go mod init && \
		test -s main.go || { echo "Populating sample main.go"; git archive --remote=ssh://git@gitlab.eng.vmware.com/nsx-allspark_users/m7/policymodel.git nexus-sdk-dev sample-main.go.template | tar -x ; mv sample-main.go.template main.go; } && \
		if [ -n ${DATAMODEL} ]; then\
			go mod edit -replace gitlab.eng.vmware.com/nsx-allspark_users/m7/policymodel.git=${DATAMODEL_DIR}/build/${DATAMODEL}/policymodel && \
			go mod edit -replace gitlab.eng.vmware.com/nsx-allspark_users/m7/policymodel.git/pkg/k8sclients/1.22=${DATAMODEL_DIR}/build/${DATAMODEL}/policymodel/pkg/k8sclients/1.22 ; \
		fi

.PHONY: app_run
app_run: kill-portforward start-portforward
	cd ${APP_DIR} ; \
	{ go get .  ; go mod vendor  ; go run main.go -host http://localhost:45192; } || { kill $(shell cat /tmp/portforward.pid); rm /tmp/portforward.pid; exit 1; } ; 
