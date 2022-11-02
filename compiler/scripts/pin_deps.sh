#! /bin/bash

set -e

DEFAULT_CLIENT_NAME="$(yq eval .k8s_clients.default $( dirname "$0" )/../manifest.yaml)"
DEFAULT_CLIENT_VERSION_TAG=$(printf "%s" $(yq eval -o=json .k8s_clients.versioned $( dirname "$0" )/../manifest.yaml | jq -c  '.[]' | while read i; do
  NAME=$( jq -r  '.name' <<< "${i}" )
  if [ $NAME = $DEFAULT_CLIENT_NAME ]; then
    echo $( jq -r  '.k8s_code_generator_git_tag' <<< "${i}" )
    break
  fi
done
))
if [[ -z $DEFAULT_CLIENT_VERSION_TAG ]]; then
  echo "Could not determine default k8s client, exiting..."
  exit 1
fi

COMPILER_SRC_DIRECTORY=$1

go mod edit -require github.com/elliotchance/orderedmap@v1.4.0
go mod edit -require github.com/cespare/xxhash/v2@v2.1.2
go mod edit -require k8s.io/apimachinery@$DEFAULT_CLIENT_VERSION_TAG
go mod edit -require k8s.io/client-go@$DEFAULT_CLIENT_VERSION_TAG

go mod edit -replace k8s.io/api=k8s.io/api@$DEFAULT_CLIENT_VERSION_TAG
go mod edit -replace k8s.io/apiextensions-apiserver=k8s.io/apiextensions-apiserver@$DEFAULT_CLIENT_VERSION_TAG
go mod edit -replace k8s.io/apimachinery=k8s.io/apimachinery@$DEFAULT_CLIENT_VERSION_TAG
go mod edit -replace k8s.io/client-go=k8s.io/client-go@$DEFAULT_CLIENT_VERSION_TAG

go mod edit -replace github.com/vmware-tanzu/graph-framework-for-microservices/gqlgen=${COMPILER_SRC_DIRECTORY}/../gqlgen
go mod edit -replace github.com/vmware-tanzu/graph-framework-for-microservices/kube-openapi=${COMPILER_SRC_DIRECTORY}/../kube-openapi
go mod edit -replace github.com/vmware-tanzu/graph-framework-for-microservices/nexus=${COMPILER_SRC_DIRECTORY}/../nexus
