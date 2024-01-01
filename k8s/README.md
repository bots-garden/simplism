# Simplism on Kubernetes

## Define environment variables

> update the `env` file

```bash
KUBE_NAMESPACE="simplism-faas"
IMAGE_NAME="k33g/simplism"
IMAGE_TAG="0.0.8"
DNS="1f833ec8-f509-46f5-98ad-9e57465fde32.k8s.civo.com"
```

## Create a wasm storage pod
> - Ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-persistent-volume-storage/
> - Not for production use
```bash
set -o allexport; source .env; set +o allexport
# Create namespace (if needed)
kubectl create namespace ${KUBE_NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f ./manifests/wasm-files-volume.yaml -n ${KUBE_NAMESPACE}
```
> deletion: `kubectl delete -f 01-wasm-files-volume.yaml -n ${KUBE_NAMESPACE}`

### Check the wasm storage

```bash
kubectl exec -n ${KUBE_NAMESPACE} -it wasm-store -- /bin/sh
ls
```

## Copy wasm files to the wasm storage

```bash
kubectl cp ./wasm-files/hello.wasm ${KUBE_NAMESPACE}/wasm-store:wasm-files/hello.wasm
kubectl cp ./wasm-files/hello-people.wasm ${KUBE_NAMESPACE}/wasm-store:wasm-files/hello-people.wasm
```

### Check again the wasm storage

```bash
kubectl exec -n ${KUBE_NAMESPACE} -it wasm-store -- /bin/sh
ls wasm-files
```

## Create "hello" Simplism pod

```bash
set -o allexport; source .env; set +o allexport # if needed
APPLICATION_NAME="hello" WASM_FILE="hello.wasm" FUNCTION_NAME="handle" \
envsubst < templates/deploy.from.volume.yaml > tmp/deploy.${APPLICATION_NAME}.yaml
kubectl apply -f tmp/deploy.${APPLICATION_NAME}.yaml -n ${KUBE_NAMESPACE}
```

### Ingress

```bash
kubectl describe ingress ${APPLICATION_NAME} -n ${KUBE_NAMESPACE}
```

```bash
curl http://${APPLICATION_NAME}.${DNS} -d '👋 Hello World 🌍 on Civo'
```

## Deploy a Simplism registry

