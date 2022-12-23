#!/usr/bin/env bash
set -e 

VERSION=""
SEMVER_REGEX_MASTER="v[0-9]+\.[0-9]+\.[0-9]+$"

json_resp=$(curl -s https://gcr.io/v2/nsx-sm/nexus/nexus-cli/tags/list | jq  -r '.manifest[] | .tag | select(.[]=="latest") | .[]')
declare -a tags_array=($(echo "${json_resp}" | tr "\n" " "))
for version in "${tags_array[@]}"; do
       if [[ "${version}" =~ ${SEMVER_REGEX_MASTER} ]]; then
                VERSION="${version}"
                break
       fi
done

echo "$VERSION" | tr -d '[:space:]'
