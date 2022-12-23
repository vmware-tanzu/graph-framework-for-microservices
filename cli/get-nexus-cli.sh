#!/usr/bin/env bash
set -e 


REPOSITORY="gcr.io/nsx-sm/nexus/nexus-cli"
INSTALL_DIRECTORY="/usr/local/bin"
VERSION=""
SEMVER_REGEX_MASTER="v[0-9]+\.[0-9]+\.[0-9]+$"
CLI_UTILS_DOCKER_NAME="cli-utils"

usage() { echo "Usage: $0 [-r <repository-name>] [-v <version>] [-d <install-directory>] " 1>&2; exit 1; }

must_exist() {
        export PATH=$PATH:$HOME/.local/bin
        if ! command -v "$1" >/dev/null 2>&1; then
                echo -e "$1" is not installed! "\n$2" >&2
                return 1
        fi
}

msg="docker is needed by this script, please install the docker to use this script to download the nexus"
must_exist "docker" "${msg}"

version_tag=$(docker run --name ${CLI_UTILS_DOCKER_NAME} -it gcr.io/nsx-sm/nexus/cli-utils)
if [[ "${version_tag}" =~ ${SEMVER_REGEX_MASTER} ]]; then
	VERSION="${version_tag}"
else
	echo "Unable to find latest nexus version"
	exit 1
fi
docker rm ${CLI_UTILS_DOCKER_NAME} &> /dev/null
docker rmi "gcr.io/nsx-sm/nexus/cli-utils:latest" &> /dev/null 


if [[ $# == 0 ]]; then
    echo -e "Downloading Nexus ...\nVersion: ${VERSION}\nImage repository: ${REPOSITORY}\nDirectory: ${INSTALL_DIRECTORY}\n"
    echo -e "Would you like to customize installation [y/n]:"
    read -r customize

    if [ "${customize}" == "y" ]; then
        echo -n "Image repository [${REPOSITORY}]:"
        read -r repository_name
        echo -n "Version [${VERSION}]:"
        read -r version
        echo -n "Directory [/usr/local/bin]:"
        read -r install_directory
    fi
elif [[ "$1" == "--no-prompt" ]]; then
    echo -e "Downloading Nexus ...\nVersion: ${VERSION}\nImage repository: ${REPOSITORY}\nDirectory: ${INSTALL_DIRECTORY}\n"
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
                install_directory=${OPTARG}
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

if [[ -n "${install_directory}" ]]; then
    INSTALL_DIRECTORY="${install_directory}"
fi

OS=`uname -s`
docker_name="nexus-cli"
darwin_src_path="/nexus/darwin/nexus"
linux_src_path="/nexus/linux/nexus"

docker ps > /dev/null || echo "Unable to run docker command"

docker rm -f ${docker_name} &> /dev/null

docker pull ${REPOSITORY}:${VERSION} 1> /dev/null
docker create --name "$docker_name" ${REPOSITORY}:${VERSION} 1> /dev/null

if [ "$OS" == "Darwin" ]; then
   docker cp "${docker_name}:${darwin_src_path}" "$INSTALL_DIRECTORY" 
elif [ "$OS" == "Linux" ]; then
   docker cp "${docker_name}:${linux_src_path}" "$INSTALL_DIRECTORY"
else
   echo "Unsupported OS type: ${OS}, nexus supported OS types are : 1) Linux and 2) Darwin"
fi

docker rm -f ${docker_name} 1> /dev/null
docker rmi "${REPOSITORY}":"${VERSION}" &> /dev/null 

echo -e "Nexus (${VERSION}) installed in $INSTALL_DIRECTORY/nexus\nRun \"nexus help\" to get started"

