#!/bin/bash

set -e

nexus-openapi-gen \
  -h ./openapi-generator/openapi/boilerplate.go.txt \
  -i gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/_generated/apis/... \
  -p gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/_generated/openapi-generator/openapi

case $PWD/ in
  $GOPATH/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/_generated/) echo "we're in GOPATH, no need to copy";;
  *) echo "we're NOT in GOPATH, need to copy generated code to repository path"; \
    cp $GOPATH/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/_generated/openapi-generator/openapi/openapi_generated.go ./openapi-generator/openapi/;;
esac

sed -i "s|github.com/go-openapi/spec|gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kube-openapi.git/pkg/validation/spec|" ./openapi-generator/openapi/openapi_generated.go
sed -i "s|k8s.io/kube-openapi/pkg/common|gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kube-openapi.git/pkg/common|" ./openapi-generator/openapi/openapi_generated.go

#go run openapi-generator/cmd/generate-openapischema/generate-openapischema.go -yamls-path crds
