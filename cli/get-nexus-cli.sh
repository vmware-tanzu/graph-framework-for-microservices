#!/usr/bin/env bash
set -e 

REPOSITORY="gcr.io/nsx-sm/nexus/nexus-cli"
VERSION="latest"
DST_DIR="$PWD"
CURR_NEXUS_PATH=$(which nexus) || true

usage() { echo "Usage: $0 [-r <repository-name>] [-v <version>] [-d <destination-path>] " 1>&2; exit 1; }

if [[ $# == 0 ]]; then
    echo "Installing nexus cli"
    echo -n "Enter the nexus cli image repository(${REPOSITORY}):"
    read -r repository_name
    echo -n "Enter the nexus cli version to be installed(${VERSION}):"
    read -r version
    if [[ -n ${CURR_NEXUS_PATH} ]]; then
       echo -n "Do you want to over-write the existing nexus cli(${CURR_NEXUS_PATH} [yes/no]):"
       read -r overwrite_path
    else
       echo -n "Enter the Path to install nexus cli(${DST_DIR}):"
       read -r dest_path
    fi
else
    while getopts ":r:v:d:" o; do
        case "${o}" in
            r)
                repository_name=${OPTARG}
                ;;
            v)
                version=${OPTARG}
                ;;
            d)
                dest_path=${OPTARG}
                ;;
            *)
                usage
                ;;
        esac
    done
    shift "$((OPTIND-1))"
fi


if [[ -n "${repository_name}" ]]; then
    REPOSITORY="${repository_name}"
fi

if [[ -n "${version}" ]]; then
    VERSION="${version}"
fi

if [[ -n "${dest_path}" ]]; then
    DST_DIR="${dest_path}"
elif [[ "${overwrite_path}" == "yes" ]]; then
    
    DST_DIR=`echo "${CURR_NEXUS_PATH}" | awk -F'/' 'BEGIN{OFS=FS} {$NF=""; print}'`
elif [[ "${overwrite_path}" == "no" ]]; then
    DST_DIR="${PWD}"
fi

echo "The nexus CLI will be downloaded from ${REPOSITORY}:${VERSION} and get installed to ${DST_DIR}"


OS=`uname -s`

docker_name="nexus-cli"
darwin_src_path="/nexus/darwin/nexus"
linux_src_path="/nexus/linux/nexus"

docker rm -f ${docker_name} &> /dev/null

docker pull ${REPOSITORY}:${VERSION} 1> /dev/null
docker create --name "$docker_name" ${REPOSITORY}:${VERSION} 1> /dev/null

if [ "$OS" == "Darwin" ]; then
   docker cp "${docker_name}:${darwin_src_path}" "$DST_DIR" 
elif [ "$OS" == "Linux" ]; then
   docker cp "${docker_name}:${linux_src_path}" "$DST_DIR"
else
   echo "Unsupported OS type: ${OS}, nexus cli supported OS types are : 1) Linux and 2) Darwin"
fi

docker rm -f ${docker_name} 1> /dev/null
docker rmi "${REPOSITORY}":"${VERSION}" &> /dev/null 

echo "The nexus cli binary downloaded here: $DST_DIR/nexus"

