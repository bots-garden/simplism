#!/bin/bash
set -o allexport; source ../../.release.env; set +o allexport

IMAGE_BASE_NAME="${GITPOD_IMAGE_BASE_NAME}"
IMAGE_TAG="${GITPOD_IMAGE_TAG}"
docker login -u ${DOCKER_USER} -p ${DOCKER_PWD}
docker buildx build \
--build-arg="GO_ARCH=amd64" \
--build-arg="GO_VERSION=1.21.3" \
--build-arg="TINYGO_ARCH=amd64" \
--build-arg="TINYGO_VERSION=0.30.0" \
--build-arg="EXTISM_VERSION=1.0.2" \
--build-arg="EXTISM_ARCH=amd64" \
--build-arg="EXTISM_JS_VERSION=1.0.0-rc3" \
--build-arg="EXTISM_JS_ARCH=x86_64" \
--build-arg="EXTISM_JS_OS=linux" \
--build-arg="BINARYEN_VERSION=version_116" \
--build-arg="BINARYEN_ARCH=x86_64" \
--build-arg="BINARYEN_OS=linux" \
--build-arg="NODE_MAJOR=20" \
--platform linux/amd64 \
--push -t ${DOCKER_USER}/${IMAGE_BASE_NAME}:${IMAGE_TAG} .

echo "to test it:"
echo "docker run -it ${DOCKER_USER}/${IMAGE_BASE_NAME}:${IMAGE_TAG}"
