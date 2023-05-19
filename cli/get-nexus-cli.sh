#!/usr/bin/env bash
set -e 


REPOSITORY="gcr.io/nsx-sm/nexus/nexus-cli"
INSTALL_DIRECTORY="/usr/local/bin"
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

if [[ -z "${VERSION}" ]]; then
   echo "Unable to get the nexus image tag, Please contact Nexus Support"
   exit 1
fi

usage() { echo "Usage: $0 [-r <repository-name>] [-v <version>] [-d <install-directory>] " 1>&2; exit 1; }

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

echo "Installing Nexus..."
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

