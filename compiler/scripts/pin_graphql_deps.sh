#! /bin/bash

set -e

COMPILER_SRC_DIRECTORY=$1
go mod edit -require sigs.k8s.io/yaml@v1.3.0
go mod edit -replace github.com/vmware-tanzu/graph-framework-for-microservices/gqlgen=${COMPILER_SRC_DIRECTORY}/../gqlgen
go mod edit -replace github.com/vmware-tanzu/graph-framework-for-microservices/kube-openapi=${COMPILER_SRC_DIRECTORY}/../kube-openapi
go mod edit -replace github.com/vmware-tanzu/graph-framework-for-microservices/nexus=${COMPILER_SRC_DIRECTORY}/../nexus
go mod edit -require github.com/cespare/xxhash/v2@v2.1.2
go mod edit -require github.com/imdario/mergo@v0.3.12