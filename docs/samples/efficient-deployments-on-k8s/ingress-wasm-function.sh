#!/bin/bash
set -o allexport; source .env; set +o allexport

kubectl describe ingress ${APPLICATION_NAME} -n ${KUBE_NAMESPACE}

