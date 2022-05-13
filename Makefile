BUCKET ?= nexus-template-downloads
TAG ?= $(shell git rev-parse --verify HEAD)
build:
	$(MAKE) archive_app_templates
archive_app_templates:
	cd templates && \
	rm -f app-template.tar && \
	tar -czvf app-template.tar *

publish:
	cd templates && \
	gsutil cp app-template.tar gs://${BUCKET}/${TAG}/;
	