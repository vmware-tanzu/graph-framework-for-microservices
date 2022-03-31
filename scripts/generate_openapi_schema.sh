#!/bin/bash

set -e

openapi-gen \
  -h ./pkg/openapi/boilerplate.go.txt \
  -i gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/_generated/apis/... \
  -p gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/openapi

case $PWD/ in
  $GOPATH/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/) echo "we're in GOPATH, no need to copy";;
  *) echo "we're NOT in GOPATH, need to copy generated code to repository path"; \
    cp $GOPATH/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/openapi/openapi_generated.go pkg/openapi/;;
esac

go run cmd/generate-openapischema/generate-openapischema.go -yamls-path _generated/crds

git checkout -- pkg/openapi/openapi_generated.go
