#!/bin/bash
set -o allexport; source .env; set +o allexport

kubectl delete -f tmp/deploy.${IMAGE_TAG}.yaml -n ${KUBE_NAMESPACE}