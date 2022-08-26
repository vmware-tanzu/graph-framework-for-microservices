#!/bin/bash
set -exo pipefail

PUB_TAG=
DOCKER_REPOSITORY="284299419820.dkr.ecr.us-west-2.amazonaws.com/nexus/tsm-openapispec"
COMMON_API_VERSION_TAG=$(cat common-api-version)

PUB_TAG="$COMMON_API_VERSION_TAG"
DOCKER_PUB_TAG=${CI_COMMIT_TAG:-$CI_COMMIT_SHA}
if [ -z $DOCKER_PUB_TAG ]; then
  echo "Please provide docker image tag to publish"
   exit 1
fi

npm i --no-package-lock @allspark/common-apis@$PUB_TAG
docker build -t $DOCKER_REPOSITORY:$DOCKER_PUB_TAG . -f Dockerfile
docker push $DOCKER_REPOSITORY:$DOCKER_PUB_TAG
