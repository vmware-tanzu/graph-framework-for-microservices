SHELL := /bin/bash

generate_test_schemas: generate_go_openapi_schema

#protoc-gen-gogo:
#	@path=$$(pwd); \
#	pushd ../../_deps/github.com/gogo/protobuf/protoc-gen-gogo; \
#		go build -o $${path}; \
#	popd;
#
#generate_structs_from_proto: protoc-gen-gogo
#	protoc --gogo_out=. --plugin=protoc-gen-gogo="protoc-gen-gogo" test_data/proto/test.proto

generate_go_openapi_schema:
	openapi-gen \
		-h ../openapi/boilerplate.go.txt \
		-i gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/_generated/openapi-generator/openapi_generator/test_data/... \
		-p gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/_generated/openapi-generator/openapi_generator/test_data/openapi

.PHONY: generate_test_schemas generate_structs_from_proto generate_go_openapi_schema
