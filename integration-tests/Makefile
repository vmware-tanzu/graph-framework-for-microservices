CLI_BRANCH ?= "master"
DOCS_VERSION ?= ""
run_integration_test:
	git submodule update --init --remote
	if [ -n $(DOCS_VERSION) ]; then\
		cd docs && git checkout $(DOCS_VERSION) && cd .. ;\
	fi
	nexus version
	nexus config set --disable-upgrade-prompt --debug-always --skip-upgrade-check
	go run . workflows.yaml
