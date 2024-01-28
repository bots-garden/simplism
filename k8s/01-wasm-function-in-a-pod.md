## Deploy a Wasm Simplism function in a pod
> wasm function == simplism plug-in

First we need a manifest file:

```bash
wget https://github.com/bots-garden/simplism/releases/download/v0.1.3/deploy-wasm-from-remote.yaml
```

### Deploy functions from remote files

Simplism can download and execute Wasm plug-ins from remote URLs.

There are some Simplism plug-ins available here for demonstration purposes:
- https://github.com/simplism-registry/small-cow/releases/tag/v0.0.0
- https://github.com/simplism-registry/small-ant/releases/tag/v0.0.0


#### Deploy a first function

```bash
set -o allexport; source .env; set +o allexport
export SERVICE_NAME="small-cow"
export WASM_FILE="small-cow.wasm" 
export WASM_URL="https://github.com/simplism-registry/small-cow/releases/download/v0.0.0/small-cow.wasm"
export FUNCTION_NAME="handle"

rm -f tmp/deploy.${SERVICE_NAME}.yaml
envsubst < deploy-wasm-from-remote.yaml > tmp/deploy.${SERVICE_NAME}.yaml

# Create namespace (if needed)
kubectl create namespace ${KUBE_NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
# Deploy
kubectl apply -f tmp/deploy.${SERVICE_NAME}.yaml -n ${KUBE_NAMESPACE}
```

You should see:
```bash
service/small-cow created
deployment.apps/small-cow created
ingress.networking.k8s.io/small-cow created
```

You can get the ingress of the service with this command:
```bash
kubectl describe ingress ${SERVICE_NAME} -n ${KUBE_NAMESPACE}
#simplism-faas.small-cow.ffb140f9-7479-4308-9763-9f70628794b1.k8s.civo.com
```

And now you can call the function:
```bash
curl http://${KUBE_NAMESPACE}.${SERVICE_NAME}.${DNS} -d '👋 Hello World 🌍'
# you should get: 
^__^
(oo)\_______
(__)\       )\/\
    ||----w |
    ||     ||
👋 Hello World 🌍
```

#### Let's try another function

```bash
set -o allexport; source .env; set +o allexport
export SERVICE_NAME="small-ant"
export WASM_FILE="small_ant.wasm" 
export WASM_URL="https://github.com/simplism-registry/small-ant/releases/download/v0.0.0/small_ant.wasm"
export FUNCTION_NAME="handle"
rm -f tmp/deploy.${SERVICE_NAME}.yaml
envsubst < deploy-wasm-from-remote.yaml > tmp/deploy.${SERVICE_NAME}.yaml

# Create namespace (if needed)
kubectl create namespace ${KUBE_NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
# Deploy
kubectl apply -f tmp/deploy.${SERVICE_NAME}.yaml -n ${KUBE_NAMESPACE}
```

You can get the ingress of the service with this command:
```bash
kubectl describe ingress ${SERVICE_NAME} -n ${KUBE_NAMESPACE}
#simplism-faas.small-ant.ffb140f9-7479-4308-9763-9f70628794b1.k8s.civo.com
```

And now you can call the function:
```bash
curl http://${KUBE_NAMESPACE}.${SERVICE_NAME}.${DNS} -d '✋ Hey people 🤗'
# you should get: 
/\/\
  \_\  _..._
  (" )(_..._)
   ^^  // \\
✋ Hey people 🤗
```

### Deploy functions from local files

In this case, we need to create a kubernetes volume to store the Wasm files. Then we will copy the wasm files on this volume, and then we will deploy the Wasm function from the volume.

First, add this line to the `.env` file:
```bash
VOLUME_SIZE="10Mi"
```

First we need these manifest files:

```bash
wget https://github.com/bots-garden/simplism/releases/download/v0.1.3/wasm-files-volume.yaml
wget https://github.com/bots-garden/simplism/releases/download/v0.1.3/deploy-wasm-from-volume.yaml
```

We need some wasm files, let's download them from github:
```bash
wget https://raw.githubusercontent.com/bots-garden/simplism/main/k8s/wasm-files/hello-world.wasm -O ./wasm-files/hello-world.wasm
wget https://raw.githubusercontent.com/bots-garden/simplism/main/k8s/wasm-files/small-cow.wasm -O ./wasm-files/small-cow.wasm
wget https://raw.githubusercontent.com/bots-garden/simplism/main/k8s/wasm-files/small_ant.wasm -O ./wasm-files/small_ant.wasm
``````

Then, to create the volume:

```bash
set -o allexport; source .env; set +o allexport

rm -f tmp/create.wasm.files.volume.yaml
envsubst < wasm-files-volume.yaml > tmp/create.wasm.files.volume.yaml

# Create namespace (if needed)
kubectl create namespace ${KUBE_NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
# Deploy
kubectl apply -f tmp/create.wasm.files.volume.yaml -n ${KUBE_NAMESPACE}
```

You should see:
```bash
persistentvolume/pv-wasm-files-volume created
persistentvolumeclaim/pv-claim-wasm-files created
pod/wasm-files-store created
```

Now, we can copy a local wasm file to the storage volume:

```bash
kubectl cp ./wasm-files/hello-world.wasm ${KUBE_NAMESPACE}/wasm-files-store:wasm-files/hello-world.wasm
```

Let's check if the wasm file is in the storage volume:
```bash
kubectl exec -n ${KUBE_NAMESPACE} -it wasm-files-store -- /bin/sh
ls wasm-files
# you should see: hello-world.wasm
```

Now it's time to deploy the Wasm function from the volume:

```bash
set -o allexport; source .env; set +o allexport
export SERVICE_NAME="hello-world"
export WASM_FILE="hello-world.wasm" 
export FUNCTION_NAME="handle"
rm -f tmp/deploy.${SERVICE_NAME}.yaml
envsubst < deploy-wasm-from-volume.yaml > tmp/deploy.${SERVICE_NAME}.yaml

# Create namespace (if needed)
kubectl create namespace ${KUBE_NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
# Deploy
kubectl apply -f tmp/deploy.${SERVICE_NAME}.yaml -n ${KUBE_NAMESPACE}
```

You can get the ingress of the service with this command:
```bash
kubectl describe ingress ${SERVICE_NAME} -n ${KUBE_NAMESPACE}
#simplism-faas.hello-world.ffb140f9-7479-4308-9763-9f70628794b1.k8s.civo.com
```

And now you can call the function:
```bash
curl http://${KUBE_NAMESPACE}.${SERVICE_NAME}.${DNS} -d 'Bob Morane'
# you should get: 
🤗 Hello Bob Morane
```

Get the list of the deployed Wasm functions:
```bash
set -o allexport; source .env; set +o allexport

kubectl get ingress -l component=simplism-function --namespace simplism-faas
kubectl get service -l component=simplism-function --namespace simplism-faas
kubectl get deployment -l component=simplism-function --namespace simplism-faas
```