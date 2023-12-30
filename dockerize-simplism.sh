#!/bin/bash
set -o allexport; source .release.env; set +o allexport
set -o allexport; source .docker.env; set +o allexport
# You need to create a .docker.env file with these variables
# - DOCKER_USER
# - DOCKER_PWD
# ✋ use this if needed: sudo chmod 666 /var/run/docker.sock

env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o linux/arm64/${APPLICATION_NAME}
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o linux/amd64/${APPLICATION_NAME}

docker login -u ${DOCKER_USER} -p ${DOCKER_PWD}
docker buildx create --use
docker buildx build -t ${DOCKER_USER}/${IMAGE_BASE_NAME}:${IMAGE_TAG} --platform=linux/arm64,linux/amd64 . --push

docker pull ${DOCKER_USER}/${IMAGE_BASE_NAME}:${IMAGE_TAG}
docker images | grep ${DOCKER_USER}/${IMAGE_BASE_NAME}