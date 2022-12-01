#!/usr/bin/env bash
set -e 

REPOSITORY="gcr.io/nsx-sm/nexus/nexus-cli"
VERSION="latest"
DST_DIR="/usr/local/bin"

usage() { echo "Usage: $0 [-r <repository-name>] [-v <version>] [-d <destination-path>] " 1>&2; exit 1; }

if [[ $# == 0 ]]; then
    echo -e "Downloading Nexus ...\nVersion: ${VERSION}\nImage repository: ${REPOSITORY}\nDirectory: ${DST_DIR}\n"
    echo -e "Would you like to customize installation [y/n]:"
    read -r customize

    if [ "${customize}" == "y" ]; then
        echo -n "Image repository [${REPOSITORY}]:"
        read -r repository_name
        echo -n "Version [${VERSION}]:"
        read -r version
        echo -n "Directory [/usr/local/bin]:"
        read -r dest_path
    fi
elif [[ "$1" == "--no-prompt" ]]; then
    echo -e "Downloading Nexus ...\nVersion: ${VERSION}\nImage repository: ${REPOSITORY}\nDirectory: ${DST_DIR}\n"
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
fi

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
   echo "Unsupported OS type: ${OS}, nexus supported OS types are : 1) Linux and 2) Darwin"
fi

docker rm -f ${docker_name} 1> /dev/null
docker rmi "${REPOSITORY}":"${VERSION}" &> /dev/null 

echo -e "Nexus (${VERSION}) installed in $DST_DIR/nexus\nRun \"nexus help\" to get started"

