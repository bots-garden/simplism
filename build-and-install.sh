#!/bin/bash
set -o allexport; source .release.env; set +o allexport
echo -n $MESSAGE > cmds/version.txt

echo -n $SIMPLISM_IMAGE > generators/docker.image.txt
echo -n $GITPOD_IMAGE > generators/gitpod.image.txt
echo -n $TAG > generators/simplism.version.txt


go build
ls -lh simplism
#cp simplism /usr/bin
sudo cp simplism /usr/bin
