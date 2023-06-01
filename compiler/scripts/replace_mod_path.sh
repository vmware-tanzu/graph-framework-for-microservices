#!/bin/bash

set -e

GENERATED_PACKAGE="nexustempmodule"

if [[ -n "$CRD_MODULE_PATH" ]]; then
  pushd _generated
  echo "Update import paths with user's CRD module path..."
  find ./ \( -type d -name .git -prune \) -o \( -type f ! -name "server" \) -print0 | xargs -0 sed -i'.bak' -e "s|${GENERATED_PACKAGE}/|${CRD_MODULE_PATH}|g"
  popd
fi
