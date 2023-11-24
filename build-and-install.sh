#!/bin/bash
set -o allexport; source .release.env; set +o allexport
echo -n $MESSAGE > cmds/version.txt

go build
cp simplism /usr/bin
