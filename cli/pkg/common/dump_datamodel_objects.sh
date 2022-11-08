#!/bin/bash

apiserver=$1
TMP_LOCATION="/tmp"

if [ "$#" -ne 1 ]; then
    echo -e "Insufficient number of args.\nPlease provide valid api-server. Example: <apiserver-hostname:port> --> localhost:5000"
    exit 1
fi

cdate=`date +%Y_%m_%d_%H_%M_%S`
DATAMODEL_DUMP_PATH="${TMP_LOCATION}/datamodel_dump_${cdate}"
DATAMODEL_DATA_FILE_NAME="${TMP_LOCATION}/datamodel_dump_${cdate}.tar.gz"

mkdir -p "${DATAMODEL_DUMP_PATH}"

crdlist=( $(kubectl -s ${apiserver} get crds --no-headers | awk '{print $1}') )
echo "Dumping ${#crdlist[@]} datamodel CR's objects"

for crd in "${crdlist[@]}"
do
   #echo "$crd"
   mkdir -p "${DATAMODEL_DUMP_PATH}/${crd}"
   kubectl -s ${apiserver} get $crd -o yaml > "${DATAMODEL_DUMP_PATH}/${crd}/${crd}.yaml"

   cr_objs=( $(kubectl -s ${apiserver} get $crd --no-headers | awk '{print $1}') )
   for obj in "${cr_objs[@]}"
   do
       kubectl -s ${apiserver} get $crd ${obj} -o yaml > "${DATAMODEL_DUMP_PATH}/${crd}/${obj}.yaml"
   done
done
 
tar -czf ${DATAMODEL_DATA_FILE_NAME} -C ${DATAMODEL_DUMP_PATH} .
echo "Please find the datamodel objects here : ${DATAMODEL_DATA_FILE_NAME}"
rm -rf ${DATAMODEL_DUMP_PATH}
