#!/bin/bash

# Copyright 2021 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

source $(dirname "$0")/common.sh

echo "Checking for license header..."
allfiles=$(listFiles|grep -v ./internal/bindata/...)
licRes=""
for file in $allfiles; do
  if ! head -n4 "${file}" | grep -Eq "(Copyright|generated|GENERATED|Licensed)" ; then
    licRes="${licRes}\n"$(echo -e "  ${file}")
  fi
done
if [ -n "${licRes}" ]; then
  echo -e "license header checking failed:\n${licRes}"
  exit 255
fi
