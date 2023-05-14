#!/bin/bash
set -exo pipefail

DOCKER_REPOSITORY="gcr.io/nsx-sm/nexus/nexus-cli"
DOCKER_PUB_TAG=${CI_COMMIT_TAG}

if [ -z "$DOCKER_PUB_TAG" ]; then
  echo "Please provide docker image tag to publish nexus-cli"
  exit 1
fi

SEMVER_REGEX_MASTER="v[0-9]+\.[0-9]+\.[0-9]+$"
ALPHANUM='[0-9]*[A-Za-z-][0-9A-Za-z-]*'
SEMVER_REGEX_PRIVATE="^v[0-9]+.[0-9]+\.[0-9]+-($ALPHANUM)$"
release_tag=""


if [[ "$DOCKER_PUB_TAG" =~ $SEMVER_REGEX_MASTER ]]; then
   echo "The tag is from master"
   release_tag="true"
elif [[ "$DOCKER_PUB_TAG" =~ $SEMVER_REGEX_PRIVATE ]]; then
   echo "The tag is from private branch"
else
   echo "The nexus cli tag is not generated from master."
   exit 1
fi

make build

docker build -t "$DOCKER_REPOSITORY":"$DOCKER_PUB_TAG" . -f Dockerfile.cli
docker push "$DOCKER_REPOSITORY":"$DOCKER_PUB_TAG"

if [ -n "$release_tag" ]
then
    docker tag "$DOCKER_REPOSITORY":"$DOCKER_PUB_TAG" "$DOCKER_REPOSITORY":"latest"
    docker push "$DOCKER_REPOSITORY":"latest"
fi
