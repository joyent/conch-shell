#!/usr/bin/env bash

PREFIX="joyent"
NAME="conch-shell"
: ${BUILDNUMBER:=0}

: ${BUILDER:=${USER}}
BUILDER=$(echo "${BUILDER}" | sed 's/\//_/g' | sed 's/-/_/g')

: ${LABEL:="latest"}
LABEL=$(echo "${LABEL}" | sed 's/\//_/g')

IMAGE_NAME="${PREFIX}/${NAME}:${LABEL}"

PWD=$(pwd)
TAG=`git describe`
HASH=`git rev-parse HEAD`

docker build \
	-t ${IMAGE_NAME} \
	--build-arg VERSION=${TAG} \
	--build-arg VCS_REF=${HASH} \
	--file Dockerfile . \
&& \
docker run --rm \
	--name ${BUILDER}_${BUILDNUMBER} \
	--mount type=bind,source="${PWD}/release",target="/go/src/github.com/joyent/conch-shell/release" \
	--entrypoint=make \
	${IMAGE_NAME} \
	release checksums
