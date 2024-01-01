#!/bin/bash
set -o allexport; source ../.env; set +o allexport

# Create namespace (if needed)
kubectl create namespace ${KUBE_NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -

#kubectl create -f pvc.yaml
kubectl apply -f volume.yaml -n ${KUBE_NAMESPACE}

kubectl get pv task-pv-volume


: << comment
This script creates a volume on a kubernetes cluster
Question: hot to get the list of the existing volumes?
same question with K9S
How to copy files?

comment
