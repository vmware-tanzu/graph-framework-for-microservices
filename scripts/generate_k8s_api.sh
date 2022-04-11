#! /bin/bash

set -e

ROOT_PACKAGE="gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git"
GENERATED_PACKAGE="${ROOT_PACKAGE}/_generated"

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
./_deps/github.com/kubernetes/code-generator/generate-groups.sh all "${GENERATED_PACKAGE}/client" "${GENERATED_PACKAGE}/apis" "${API_NAMES}" --go-header-file "./_deps/github.com/kubernetes/code-generator/hack/boilerplate.go.txt"

case $PWD/ in
  $GOPATH/src/${ROOT_PACKAGE}/) echo "we're in GOPATH, no need to copy";;
  *) echo "we're NOT in GOPATH, need to copy generated code to repository path"; \
    cp -r $GOPATH/src/${GENERATED_PACKAGE}/* _generated/;;
esac

if [[ -v CRD_MODULE_PATH ]]; then
  echo "Update import paths with user's CRD module path..."
  find _generated \( -type d -name .git -prune \) -o -type f -print0 | xargs -0 sed -i "s|${GENERATED_PACKAGE}/|${CRD_MODULE_PATH}|g"
fi

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
