#!/bin/bash

set -x

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

docker compose -f ${SCRIPT_DIR}/compose.yaml up -d
sleep 5

while true 
do
  docker compose -f ${SCRIPT_DIR}/compose.yaml ps  --format json  | grep running
  if [ $? -eq 0 ]; then 
    docker exec k0s cat /var/lib/k0s/pki/admin.conf > ${SCRIPT_DIR}/.kubeconfig 
    exit 0
  else
    sleep 2
  fi
done
