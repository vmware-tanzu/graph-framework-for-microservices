#! /bin/bash

set -e

GENERATED_PACKAGE="nexustempmodule"

. $(dirname "$0")/api_names.sh

# default client generated to main module
DEFAULT_CLIENT_NAME="$(yq eval .k8s_clients.default manifest.yaml)"
DEFAULT_CLIENT_VERSION_TAG=$(printf "%s" $(yq eval -o=json .k8s_clients.versioned manifest.yaml | jq -c  '.[]' | while read i; do
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
echo "Generating default client"
pushd ./_deps/github.com/kubernetes/code-generator
  git checkout -f $DEFAULT_CLIENT_VERSION_TAG
popd
pushd _generated
../_deps/github.com/kubernetes/code-generator/generate-groups.sh all "${GENERATED_PACKAGE}/client" "${GENERATED_PACKAGE}/apis" "${API_NAMES}" --go-header-file "../_deps/github.com/kubernetes/code-generator/hack/boilerplate.go.txt"
popd

cp -r $GOPATH/src/${GENERATED_PACKAGE}/* _generated

# versioned clients generated to dedicated modules
#yq eval -o=json .k8s_clients.versioned manifest.yaml | jq -c  '.[]?' | while read i; do
#  MAJOR_VERSION=$( jq -r  '.k8s_version' <<< "${i}" )
#  GIT_TAG=$( jq -r  '.k8s_code_generator_git_tag' <<< "${i}" )
#  echo "Generating ${MAJOR_VERSION} client"
#  pushd ./_deps/github.com/kubernetes/code-generator
#    git checkout -f $GIT_TAG
#  popd
#  ./_deps/github.com/kubernetes/code-generator/generate-groups.sh all "${ROOT_PACKAGE}/_generated/k8sclients/${MAJOR_VERSION}/client" "${ROOT_PACKAGE}/_generated/apis" "${API_NAMES}" --go-header-file "./_deps/github.com/kubernetes/code-generator/hack/boilerplate.go.txt"
#done

# restore code-generator to default
pushd ./_deps/github.com/kubernetes/code-generator
  git checkout -f $DEFAULT_CLIENT_VERSION_TAG
popd
