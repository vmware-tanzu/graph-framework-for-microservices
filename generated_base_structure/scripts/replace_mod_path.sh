#!/bin/bash

set -e

ROOT_PACKAGE="gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git"
GENERATED_PACKAGE="${ROOT_PACKAGE}/_generated"

if [[ -v CRD_MODULE_PATH ]]; then
  echo "Update import paths with user's CRD module path..."
  find \( -type d -name .git -prune \) -o -type f -print0 | xargs -0 sed -i "s|${GENERATED_PACKAGE}/|${CRD_MODULE_PATH}|g"
fi
