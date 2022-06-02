BUCKET ?= nexus-template-downloads
TAG ?= $(shell git rev-parse --verify HEAD)
build:
	$(MAKE) archive_runtime_manifest
archive_runtime_manifest:
	rm -f runtime-manifests.tar && \
	tar -czvf runtime-manifests.tar runtime-manifests ;

publish:
	gsutil cp runtime-manifests.tar gs://${BUCKET}/${TAG}/;