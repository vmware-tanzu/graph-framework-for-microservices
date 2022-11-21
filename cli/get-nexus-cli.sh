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
       echo "Input not understood. See '--help' for information on using this command."
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
   #VERSION="v0.0.134"
   VERSION="latest"
fi

if [ -z "$DST_DIR" ]
then
   DST_DIR="/usr/local/bin"
fi

OS=`uname -s`

docker_name="nexus-cli"
darwin_src_path="/nexus/darwin/nexus"
linux_src_path="/nexus/linux/nexus"

docker create --name "$docker_name" ${REPOSITORY}:${VERSION}

if [ "$OS" == "Darwin" ]; then
   docker cp "${docker_name}:${darwin_src_path}" "$DST_DIR/nexus" 
elif [ "$OS" == "Linux" ]; then
   docker cp "${docker_name}:${linux_src_path}" "$DST_DIR/nexus" 
else
   echo "Unsupported OS type: ${OS}, nexus cli supported OS types are : 1) Linux and 2) Darwin"
fi

docker rm -f ${docker_name}

