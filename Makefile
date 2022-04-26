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