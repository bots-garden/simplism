#!/bin/bash
set -o allexport; source .env; set +o allexport
docker exec --workdir /${WORKDIR} -it ${CONTAINER_NAME} \
/bin/bash


