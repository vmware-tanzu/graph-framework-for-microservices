#!/bin/bash
set -e

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
       help
       exit 0
       ;;
    *)
       echo "Input is invalid. See '--help' for information on using this command."
       exit 1
       ;;
  esac
done

default_path=""

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
   default_path="$PWD"
   DST_DIR="$default_path"
fi

OS=`uname -s`

docker_name="nexus-cli"
darwin_src_path="/nexus/darwin/nexus"
linux_src_path="/nexus/linux/nexus"

docker create --name "$docker_name" ${REPOSITORY}:${VERSION} 1> /dev/null

if [ "$OS" == "Darwin" ]; then
   docker cp "${docker_name}:${darwin_src_path}" "$DST_DIR/nexus" 
elif [ "$OS" == "Linux" ]; then
   docker cp "${docker_name}:${linux_src_path}" "$DST_DIR/nexus"
else
   echo "Unsupported OS type: ${OS}, nexus cli supported OS types are : 1) Linux and 2) Darwin"
fi

docker rm -f ${docker_name} 1> /dev/null

if [ -n "$default_path" ]
then
    echo -e "The nexus cli binary downloaded here: $default_path/nexus\nPlease place the nexus cli binary in the User PATH, hence it's accessible from anywhere in the shell"
else
    echo "The nexus cli binary downloaded here: $DST_DIR/nexus"
fi


