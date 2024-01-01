#!/bin/bash
set -o allexport; source ../.env; set +o allexport
set -o allexport; source .env; set +o allexport

# Create namespace (if needed)
kubectl create namespace ${KUBE_NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -

envsubst < templates/deploy.tpl.yaml > tmp/deploy.${IMAGE_TAG}.yaml
kubectl apply -f tmp/deploy.${IMAGE_TAG}.yaml -n ${KUBE_NAMESPACE}

: << comment


comment
