#!/bin/bash

set -ex

pushd _generated
nexus-openapi-gen \
  -h ../pkg/openapi_generator/openapi/boilerplate.go.txt \
  -i nexustempmodule/apis/... \
  -p github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/openapi_generator/openapi

case $PWD/ in
  $GOPATH/src/github.com/vmware-tanzu/graph-framework-for-microservices/compiler/_generated/) echo "we're in GOPATH, no need to copy";;
  *) echo "we're NOT in GOPATH, need to copy generated code to repository path"; \
    cp $GOPATH/src/github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/openapi_generator/openapi/openapi_generated.go ../pkg/openapi_generator/openapi/;;
esac

sed -i'.bak' -e "s|github.com/go-openapi/spec|github.com/vmware-tanzu/graph-framework-for-microservices/kube-openapi/pkg/validation/spec|" ../pkg/openapi_generator/openapi/openapi_generated.go
sed -i'.bak' -e "s|k8s.io/kube-openapi/pkg/common|github.com/vmware-tanzu/graph-framework-for-microservices/kube-openapi/pkg/common|" ../pkg/openapi_generator/openapi/openapi_generated.go
popd
