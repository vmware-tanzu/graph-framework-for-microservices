#!/bin/bash

set -ex

NAME=${NAME:-}
DATAMODEL_IMAGE=${IMAGE}
IMAGE="file://${NAME}"
TITLE=${TITLE:-}
SKIP_CRD_INSTALLATION=${SKIP_CRD_INSTALLATION:-false}
GRAPHQL_ENABLED=${GRAPHQL_ENABLED:-false}

### User can pass the custom HTTP URL where the graphql plugin can be downloaded via internet
GRAPHQL_PATH=${GRAPHQL_PATH:-NA}

### This is to support older way of installing datamodel from local folder
if [[ $SKIP_CRD_INSTALLATION == "false" ]]; then
    kubectl apply -f /crds --recursive
    [[ $GRAPHQL_PATH != NA ]] && GRAPHQL_ENABLED=true
    if  test -f /build/server; then
        GRAPHQL_ENABLED=true
    fi
fi
echo $NAME
### We will create datamodel object
if [[ -n $NAME ]] && [[ -n $IMAGE ]]; then
  if [[ -n $TITLE ]]; then
    echo '
      apiVersion: nexus.vmware.com/v1
      kind: Datamodel
      metadata:
        name: '"$NAME"'
      spec:
        name: '"$NAME"'
        url: '"$IMAGE"'
        title: '"$TITLE"'
        enableGraphql: '"$GRAPHQL_ENABLED"'' | kubectl apply -f -
  else
    echo '
    apiVersion: nexus.vmware.com/v1
    kind: Datamodel
    metadata:
      name: '"$NAME"'
    spec:
      name: '"$NAME"'
      url: '"$IMAGE"'
      enableGraphql: '"$GRAPHQL_ENABLED"'' | kubectl apply -f -
  fi
fi
