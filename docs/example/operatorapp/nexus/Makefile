UNAME := $(shell uname)
NAMESPACE ?= default
BUILD_DIR ?= build
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

runtime_uninstall:
	mkdir -p /tmp/ && \
	cd /tmp  && \
	curl https://storage.googleapis.com/nexus-template-downloads/runtime-manifests.tar -o manifests.tar && \
	tar -xvf manifests.tar && \
	kubectl delete --ignore-not-found=true -f runtime-manifests/crds/ -n ${NAMESPACE} &&\
	kubectl delete --ignore-not-found=true -f runtime-manifests/deployment/ -n ${NAMESPACE};

check_pod:
	CURRENT=0
	while [ `kubectl get pods -n ${NAMESPACE} ${POD} --no-headers 2>/dev/null | wc -l` -eq 0 ] ; do \
		CURRENT=$$(( CURRENT + 1 )) ;\
		TRIES_LEFT=$$(( 20 - CURRENT )) ;\
		if [ $$TRIES_LEFT -ne 0 ] ; then \
			sleep 6 ;\
		else \
			echo "${POD} pod not created"; \
			exit 1 ; \
		fi \
	done ;
post-install-checks:
	$(MAKE) check_pod POD="-lapp=nexus-etcd" ;\
	kubectl wait pods -lapp=nexus-etcd -n ${NAMESPACE} --for=condition=ready --timeout=300s || { echo "etcd pod could not be started"; exit 1; } ;\
	$(MAKE) check_pod POD="-lapp=nexus-kube-apiserver" ;\
	kubectl wait pods -lapp=nexus-kube-apiserver -n ${NAMESPACE} --for=condition=ready --timeout=120s || { echo "k8sapiserver pod could not be started"; exit 1; } ;\
	$(MAKE) check_pod POD="-lname=nexus-kube-controllermanager" ;\
	kubectl wait pods -lname=nexus-kube-controllermanager -n ${NAMESPACE} --for=condition=ready --timeout=120s || { echo "k8scontrollermanager pod could not be started"; exit 1; } ;

.SILENT:
start-portforward: check-portforward
	( kubectl port-forward svc/nexus-proxy-container -n ${NAMESPACE} 45192:80 & echo $$! > /tmp/portforward.pid )
	CURRENT=0
	while [ `curl -s -o /dev/null -w "%{http_code}" http://localhost:45192/api/v1/namespaces 2>/dev/null` -ne 200  ] ; do \
		CURRENT=$$(( CURRENT + 1 )) ;\
		TRIES_LEFT=$$(( 10 - CURRENT )) ;\
		if [ $$TRIES_LEFT -ne 0 ] ; then \
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
	@if [ `kubectl get ns ${NAMESPACE} --no-headers 2>/dev/null | wc -l` -eq 0 ]; then \
		echo "Please install runtime on namespace" ;\
		exit 1 ;\
	fi
	@if echo -e '\x1dclose\x0d' | telnet localhost 45192 2>/dev/null; then\
		echo "could not portforward to 45192 port - port already open"; exit 1; \
	else \
		echo "Proceeding to start portforwarding.." ;\
	fi
	@if [ `kubectl get svc -n $(NAMESPACE) nexus-proxy-container --no-headers 2>/dev/null | wc -l` -eq 0 ]; then \
		echo "Please install runtime components.." ;\
		exit 1;\
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

check-datamodel:
	@if [ "${DATAMODEL}" == "" ]; then \
		{ echo "please provide DATAMODEL name as DATAMODEL=<<>> "; exit 1; } ;\
	fi

.PHONY: datamodel_install
datamodel_install: check-datamodel start-portforward
	kubectl -s localhost:45192 apply -f ${DATAMODEL}/build/crds/ || { echo "could not install datamodel"; kill $(shell cat /tmp/portforward.pid); rm /tmp/portforward.pid; exit 1; } && \
	kill $(shell cat /tmp/portforward.pid); rm /tmp/portforward.pid;

