# Deploy a Wasm Simplism Registry on Kubernetes

> Disclaimer: the registry mode of Simplism is a worl in progress (at the beginning of the project, it was a facility to do some tests) and is not ready entirely for production use yet. (== it's not battle tested)

First, add these variables to the `.env` file:

```bash
PRIVATE_REGISTRY_TOKEN="people-are-strange"
ADMIN_REGISTRY_TOKEN="morrison-hotel"
REGISTRY_SIZE="10Mi"
```

Then, we need these manifest files:

```bash
wget https://github.com/bots-garden/simplism/releases/download/v0.1.3/wasm-registry-volume.yaml 
wget https://github.com/bots-garden/simplism/releases/download/v0.1.3/deploy-wasm-registry.yaml
```

## Create a space to store the Wasm files of the Registry

```bash
set -o allexport; source .env; set +o allexport

rm -f tmp/create.wasm.registry.volume.yaml
envsubst < wasm-registry-volume.yaml > tmp/create.wasm.registry.volume.yaml

# Create namespace (if needed)
kubectl create namespace ${KUBE_NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
# Deploy
kubectl apply -f tmp/create.wasm.registry.volume.yaml -n ${KUBE_NAMESPACE}
```

You should see:
```bash
persistentvolume/task-pv-wasm-registry-volume configured
persistentvolumeclaim/task-pv-wasm-registry-claim configured
pod/wasm-registry-store configured
```

## Deploy the Wasm Registry

```bash
set -o allexport; source .env; set +o allexport

rm -f tmp/deploy.wasm.registry.yaml
envsubst < deploy-wasm-registry.yaml > tmp/deploy.wasm.registry.yaml

# Create namespace (if needed)
kubectl create namespace ${KUBE_NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
# Deploy
kubectl apply -f tmp/deploy.wasm.registry.yaml -n ${KUBE_NAMESPACE}
```

You should see:
```bash
service/wasm-registry created
deployment.apps/wasm-registry created
ingress.networking.k8s.io/wasm-registry created
```

You can get the ingress of the service with this command:
```bash
kubectl describe ingress wasm-registry -n ${KUBE_NAMESPACE}
#simplism-faas.registry.ffb140f9-7479-4308-9763-9f70628794b1.k8s.civo.com 
```

And now you can check the Wasm Registry:
```bash
curl http://${KUBE_NAMESPACE}.registry.${DNS}
# you should get: 
🖖 Live long and prosper 🤗 | simplism v0.1.3
# yes I'm a Trekkie
```

## Publish some Wasm plug-ins to the Registry

```bash
curl http://${KUBE_NAMESPACE}.registry.${DNS}/registry/push \
-H 'admin-registry-token: morrison-hotel' \
-F 'file=@./wasm-files/small-cow.wasm'

curl http://${KUBE_NAMESPACE}.registry.${DNS}/registry/push \
-H 'admin-registry-token: morrison-hotel' \
-F 'file=@./wasm-files/small_ant.wasm'

```

Now you can get the list of the wasm files:
```bash
curl http://${KUBE_NAMESPACE}.registry.${DNS}/registry/discover \
-H 'private-registry-token: people-are-strange' \
-H 'content-type:text/plain; charset=UTF-8'

# you should get this output:
+------------------+--------------------------------------+---------+---------------------+
|       NAME       |                 PATH                 |  SIZE   |       CREATED       |
+------------------+--------------------------------------+---------+---------------------+
| small-cow.wasm   | wasm-registry-files/small-cow.wasm   |  194165 | 2024-01-26 07:27:00 |
| small_ant.wasm   | wasm-registry-files/small_ant.wasm   | 2246665 | 2024-01-26 07:27:07 |
+------------------+--------------------------------------+---------+---------------------+
```

If you prefer a JSON output:
```bash
curl http://${KUBE_NAMESPACE}.registry.${DNS}/registry/discover \
-H 'private-registry-token: people-are-strange' \
-H 'content-type:application/json; charset=UTF-8'
```

Now, you can download "manually" the Wasm files:
```bash
curl http://${KUBE_NAMESPACE}.registry.${DNS}/registry/pull/small-cow.wasm -o small-cow.wasm \
-H 'private-registry-token: people-are-strange'
```

So, let's deploy the Wasm functions using the registry.

## Deploy a Wasm function from the registry

Then, we need this manifest file:

```bash
wget https://github.com/bots-garden/simplism/releases/download/v0.1.3/deploy-wasm-from-registry.yaml
```

```bash
set -o allexport; source .env; set +o allexport
export SERVICE_NAME="small-cow"
export WASM_FILE="small-cow.wasm" 
export WASM_URL="http://${KUBE_NAMESPACE}.registry.${DNS}/registry/pull/${WASM_FILE}"
export FUNCTION_NAME="handle"

# as we already deploy this service, remove it before
kubectl delete -f tmp/deploy.${SERVICE_NAME}.yaml -n ${KUBE_NAMESPACE}

rm -f tmp/deploy.${SERVICE_NAME}.yaml
envsubst < deploy-wasm-from-registry.yaml > tmp/deploy.${SERVICE_NAME}.yaml

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

🎉 Let's do it for the other Wasm function:

```bash
set -o allexport; source .env; set +o allexport
export SERVICE_NAME="small-ant"
export WASM_FILE="small_ant.wasm" 
export WASM_URL="http://${KUBE_NAMESPACE}.registry.${DNS}/registry/pull/${WASM_FILE}"
export FUNCTION_NAME="handle"

# as we already deploy this service, remove it before
kubectl delete -f tmp/deploy.${SERVICE_NAME}.yaml -n ${KUBE_NAMESPACE}

rm -f tmp/deploy.${SERVICE_NAME}.yaml
envsubst < deploy-wasm-from-registry.yaml > tmp/deploy.${SERVICE_NAME}.yaml

# Create namespace (if needed)
kubectl create namespace ${KUBE_NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
# Deploy
kubectl apply -f tmp/deploy.${SERVICE_NAME}.yaml -n ${KUBE_NAMESPACE}
```

You should see:
```bash
service/small-ant created
deployment.apps/small-ant created
ingress.networking.k8s.io/small-ant created
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

Get the list of the deployed Wasm functions:
```bash
kubectl get ingress -l component=simplism-function --namespace simplism-faas
kubectl get service -l component=simplism-function --namespace simplism-faas
kubectl get deployment -l component=simplism-function --namespace simplism-faas
```

