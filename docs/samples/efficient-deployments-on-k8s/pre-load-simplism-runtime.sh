#!/bin/bash
set -o allexport; source .env; set +o allexport

# Create namespace (if needed)
kubectl create namespace ${KUBE_NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -

envsubst < templates/load.simplism.runtime.yaml > tmp/load.simplism.${IMAGE_TAG}.runtime.yaml



kubectl create -f tmp/load.simplism.${IMAGE_TAG}.runtime.yaml -n ${KUBE_NAMESPACE}
#kubectl apply -f tmp/load.simplism.${IMAGE_TAG}.runtime.yaml -n ${KUBE_NAMESPACE}
kubectl delete -f tmp/load.simplism.${IMAGE_TAG}.runtime.yaml -n ${KUBE_NAMESPACE}
