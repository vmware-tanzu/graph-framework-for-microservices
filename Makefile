HELM_REGISTRY ?= oci://284299419820.dkr.ecr.us-west-2.amazonaws.com/nexus
CHART_NAME ?= nexus-runtime
VERSION ?= "v0.0.0-$(shell git rev-parse --verify HEAD)"
HARBOR_REPO_URL ?= "https://harbor-repo.vmware.com/chartrepo/nexus"
HARBOR_REPO ?= "harbor-vmware"
submodule:
	git submodule update --init --remote

build: submodule
	if [ -z $(VERSION) ]; then \
		echo "Please provide VERSION=<>" ; \
		exit 1; \
	fi
	helm dependency update $(CHART_NAME)
	helm dependency build $(CHART_NAME)
	helm package $(CHART_NAME) --version $(VERSION)

publish.ecr: build
	helm push $(CHART_NAME)-$(VERSION).tgz $(HELM_REGISTRY)

harbor.login:
	helm repo remove $(HARBOR_REPO) || echo "$(HARBOR_REPO) not present"; \
	helm repo add $(HARBOR_REPO) $(HARBOR_REPO_URL) --username $(HARBOR_USERNAME) --password $(HARBOR_PASSWORD);

publish.harbor: build
	helm cm-push $(CHART_NAME)-$(VERSION).tgz $(HARBOR_REPO)