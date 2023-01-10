.PHONY: race-unit-test
race-unit-test:
	ginkgo -race -cover ./pkg/...

.PHONY: test-fmt
test-fmt:
	test -z $$(goimports -w -l cmd pkg)

.PHONY: vet
vet:
	go vet ./cmd/... ./pkg/...

.PHONY: lint
lint:
	golangci-lint run ./cmd/... ./pkg/...

.PHONY: test
test: test-fmt vet lint race-unit-test

.PHONY: build
build:
	cd cmd/validate && \
		CGO_ENABLED=0 go build -ldflags="-w -s" .
