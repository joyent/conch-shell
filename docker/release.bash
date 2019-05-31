#!/usr/bin/env bash

PWD=$(pwd)

PREFIX="joyentbuildops"
NAME="conch-shell"

: ${BUILDNUMBER:=0}

: ${BUILDER:=${USER}}
BUILDER=$(echo "${BUILDER}" | sed 's/\//_/g' | sed 's/-/_/g')

: ${LABEL:="latest"}
LABEL=$(echo "${LABEL}" | sed 's/\//_/g')

if test $LABEL == "master"; then
	LABEL="latest"
fi

: ${BRANCH:="master"}

IMAGE_NAME="${PREFIX}/${NAME}:${LABEL}"

mkdir -p release

docker build \
	-t ${IMAGE_NAME} \
	--build-arg BRANCH=${BRANCH} \
	--build-arg UID=$(id -u) \
	--file Dockerfile.release . \
&& \
docker run --rm \
	--name ${BUILDER}_${BUILDNUMBER} \
	-u $(id -u):$(id -g) \
	--mount type=bind,source="${PWD}/release",target="/home/app/go/src/github.com/joyent/conch-shell/release" \
	${IMAGE_NAME} \
&& \
docker rmi ${IMAGE_NAME}
