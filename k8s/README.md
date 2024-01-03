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

> **Only for experiments: hostPath PersistentVolume version**:
> - Ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-persistent-volume-storage/
> - Not for production use
```bash
set -o allexport; source .env; set +o allexport
# Create namespace (if needed)
kubectl create namespace ${KUBE_NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f ./manifests/wasm-files-volume.yaml -n ${KUBE_NAMESPACE}
kubectl apply -f ./manifests/civowasm-files-volume.yaml -n ${KUBE_NAMESPACE}
```

>  **The production way: native storage class version**:
> - Ref: https://www.civo.com/docs/kubernetes/kubernetes-volumes
```bash
set -o allexport; source .env; set +o allexport
# Create namespace (if needed)
kubectl create namespace ${KUBE_NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f ./manifests/civo-wasm-files-volume.yaml -n ${KUBE_NAMESPACE}
# ⏳ wait for a moment...
```

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

## Create a "hello" Simplism pod

```bash
set -o allexport; source .env; set +o allexport # if needed
export APPLICATION_NAME="hello" 
export WASM_FILE="hello.wasm" 
export FUNCTION_NAME="handle"
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

### Create a wasm registry storage pod

> **Only for experiments: hostPath PersistentVolume version**:
```bash
set -o allexport; source .env; set +o allexport
# Create namespace (if needed)
kubectl create namespace ${KUBE_NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f ./manifests/wasm-registry-volume.yaml -n ${KUBE_NAMESPACE}
```

>  **The production way: native storage class version**:
```bash
set -o allexport; source .env; set +o allexport
# Create namespace (if needed)
kubectl create namespace ${KUBE_NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f ./manifests/civo-wasm-registry-volume.yaml -n ${KUBE_NAMESPACE}
# ⏳ wait for a moment...
```

#### Check the wasm registry storage

```bash
kubectl exec -n ${KUBE_NAMESPACE} -it wasm-registry-store -- /bin/sh
ls wasm-registry-files
```

### Create a Simplism registry pod

```bash
set -o allexport; source .env; set +o allexport

export APPLICATION_NAME="registry" 
export FUNCTION_NAME="handle"
export PRIVATE_REGISTRY_TOKEN="people-are-strange"
export ADMIN_REGISTRY_TOKEN="morrison-hotel"
rm -f tmp/deploy.${APPLICATION_NAME}.yaml
envsubst < templates/deploy.wasm.registry.yaml > tmp/deploy.${APPLICATION_NAME}.yaml
kubectl apply -f tmp/deploy.${APPLICATION_NAME}.yaml -n ${KUBE_NAMESPACE}
```

#### Ingress

```bash
kubectl describe ingress ${APPLICATION_NAME} -n ${KUBE_NAMESPACE}
```

> Check:
```bash
curl http://${APPLICATION_NAME}.${DNS}
# you should get: `🖖 Live long and prosper 🤗` (yes I'm a Trekkie)
```

#### Upload wasm files to the registry

```bash
curl http://${APPLICATION_NAME}.${DNS}/registry/push \
-H 'admin-registry-token: morrison-hotel' \
-F 'file=@wasm-files/hello.wasm'

curl http://${APPLICATION_NAME}.${DNS}/registry/push \
-H 'admin-registry-token: morrison-hotel' \
-F 'file=@wasm-files/hello-people.wasm'
```

> Check:
```bash
kubectl exec -n ${KUBE_NAMESPACE} -it wasm-registry-store -- /bin/sh
ls wasm-registry-files
```

## Create a "hello" Simplism pod with remote wasm file

```bash
set -o allexport; source .env; set +o allexport # if needed
export APPLICATION_NAME="hello-remote" 
export WASM_FILE="hello.wasm" 
export WASM_URL="http://registry.${DNS}/registry/pull/${WASM_FILE}"
export FUNCTION_NAME="handle"
export WASM_URL_AUTH_HEADER="private-registry-token=people-are-strange"
rm -f tmp/deploy.${APPLICATION_NAME}.yaml
envsubst < templates/deploy.from.remote.yaml > tmp/deploy.${APPLICATION_NAME}.yaml
kubectl apply -f tmp/deploy.${APPLICATION_NAME}.yaml -n ${KUBE_NAMESPACE}
```

### Ingress

```bash
kubectl describe ingress ${APPLICATION_NAME} -n ${KUBE_NAMESPACE}
```

> Call the function
```bash
curl http://${APPLICATION_NAME}.${DNS} -d '👋 Hello World 🌍 on Civo'
```
