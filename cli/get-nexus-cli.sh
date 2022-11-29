#!/bin/bash
set -e

usage ()
{
   echo "======================================================================================================================"
   echo "Download and install latest Nexus CLI: bash $0 "
   echo "Download and install Specific Version of Nexus CLI: bash $0 --version v0.0.134"
   echo "Download and install Specific Version of Nexus CLI in User Defined Path: bash $0 --version v0.0.134 --dst_dir /usr/local/bin"
   echo "If --dst_dir is configured --> if the given path is required sudo persmission, please run the script with sudo bash $0"
   echo "======================================================================================================================="
}


while [[ $# -gt 0 ]]; do
 case $1 in
    --repository)
       REPOSITORY="$2"
       shift
       shift
       ;;
    --version)
       VERSION="$2"
       shift
       shift
       ;;
    --dst_dir)
       DST_DIR="$2"
       shift
       shift
       ;;
    --help)
       usage
       exit 0
       ;;
    *)
       echo "Input is invalid. See '--help' for information on using this command."
       exit 1
       ;;
  esac
done


if [ -z "$REPOSITORY" ]
then
   REPOSITORY="gcr.io/nsx-sm/nexus/nexus-cli"
fi

if [ -z "$VERSION" ]
then
   VERSION="latest"
fi

if [ -z "$DST_DIR" ]
then
   DST_DIR="$PWD"
fi

OS=`uname -s`

docker_name="nexus-cli"
darwin_src_path="/nexus/darwin/nexus"
linux_src_path="/nexus/linux/nexus"

docker rm -f ${docker_name} &> /dev/null

docker pull ${REPOSITORY}:${VERSION} 1> /dev/null
docker create --name "$docker_name" ${REPOSITORY}:${VERSION} 1> /dev/null

if [ "$OS" == "Darwin" ]; then
   docker cp "${docker_name}:${darwin_src_path}" "$DST_DIR/nexus" 
elif [ "$OS" == "Linux" ]; then
   docker cp "${docker_name}:${linux_src_path}" "$DST_DIR/nexus"
else
   echo "Unsupported OS type: ${OS}, nexus cli supported OS types are : 1) Linux and 2) Darwin"
fi

docker rm -f ${docker_name} 1> /dev/null

echo "The nexus cli binary downloaded here: $DST_DIR/nexus"


