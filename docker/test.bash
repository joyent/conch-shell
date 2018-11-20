#!/usr/bin/env bash

: ${PREFIX:=$USER}
: ${NAME:="conch-shell"}

: ${BUILDNUMBER:=0}

: ${BUILDER:=${USER}}
BUILDER=$(echo "${BUILDER}" | sed 's/\//_/g' | sed 's/-/_/g')

: ${LABEL:="latest"}
LABEL=$(echo "${LABEL}" | sed 's/\//_/g')

IMAGE_NAME="${PREFIX}/${NAME}:${LABEL}"

TAG=`git describe`
HASH=`git rev-parse HEAD`

docker build \
	-t ${IMAGE_NAME} \
	--build-arg VERSION=${TAG} \
	--build-arg VCS_REF=${HASH} \
	--file Dockerfile . \
&& \
docker run \
	--rm \
	--name ${BUILDER}_${BUILDNUMBER} \
	${IMAGE_NAME}

