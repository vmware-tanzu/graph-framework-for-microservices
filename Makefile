BUCKET ?= nexus-template-downloads
TAG ?= $(shell git rev-parse --verify HEAD)
build:
	$(MAKE) archive_nexus
	$(MAKE) archive_datamodel
	$(MAKE) archive_datamodel_helloworld

archive_nexus:
	cd nexus && \
	rm -f runtime-manifests.tar nexus-template.tar && \
	tar -czvf runtime-manifests.tar runtime-manifests && \
	tar -czvf nexus-template.tar Makefile .gitignore;

archive_datamodel:
	cd nexus/.datamodel.templatedir &&\
		rm -f datamodel-templatedir.tar && \
		tar -czvf datamodel-templatedir.tar *;

archive_datamodel_helloworld:
	cd nexus/helloworld &&\
		rm -f helloworld-example.tar && \
		tar -czvf  helloworld-example.tar *;

publish:
	cd nexus && \
	gsutil cp runtime-manifests.tar gs://${BUCKET}/${TAG}/ && \
	gsutil cp nexus-template.tar gs://${BUCKET}/${TAG}/ && \
	cd .datamodel.templatedir &&\
	gsutil cp datamodel-templatedir.tar gs://${BUCKET}/${TAG}/ && \
	cd ../helloworld  && \
	gsutil cp helloworld-example.tar gs://${BUCKET}/${TAG}/;