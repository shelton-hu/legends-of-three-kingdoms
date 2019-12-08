#!/bin/bash

########################################################################
########## The image versions of Phoenix and related services ##########
########################################################################

### Ltk version latest
ltk_latest="
  VERSION=latest
"

### Ltk version v0.0.1
ltk_v0_0_1="
  VERSION=v0.0.1
"

## Usage:
## sh version.sh [ltk_${version}]

VERSION=$1
DEFAULT_VERSION="ltk_latest"
if [ "x${VERSION}" == "x" ]; then
  VERSION=${DEFAULT_VERSION}
fi

# get ltk version, eg: ltk_v0_0_1
LTK_VERSION=${VERSION//[.-]/_}
VERSIONS=`eval echo '$'"${LTK_VERSION}"`
# check if the given version exist
if [ "x${VERSIONS}" == "x" ]; then
  echo "The version ${VERSION} of ltk not exist!"
  exit 1
fi

# echo versions
for V in ${VERSIONS} ; do
  export ${V}
  echo {{}}${V}
done

# ltk latest images
ltk_latest_images="
  IMAGE=ranklier/ltk:latest
  LTK_MYSQL_PRISMA_DEPLOY_IMAGE=ranklier/ltk-mysql-prisma-deploy:latest
"

# ltk images
ltk_images="
  IMAGE=ranklier/ltk:${VERSION}
  LTK_MYSQL_PRISMA_DEPLOY_IMAGE=ranklier/ltk-mysql-prisma-deploy:${VERSION}
"

IMAGES=${ltk_images}
if [ "x${PH_VERSION}" == "x${DEFAULT_VERSION}" ]; then
  IMAGES=${ltk_latest_images}
fi

# echo images
for I in ${IMAGES} ; do
  echo ${I}
done

# prisma endpoints
ltk_prisma_ednpoints="
  LTK_MYSQL_PRISMA_ENDPOINT=http://ltk-mysql-prisma:4466/ltk/mysql
"

# echo prisma-endpoints
ENDPOINTS=${ltk_prisma_ednpoints}
for E in ${ENDPOINTS} ; do
  echo ${E}
done

