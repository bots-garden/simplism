# Deploy a Wasm Simplism function in a pod
> wasm function == simplism plug-in

First we need a manifest file:

```bash
wget https://github.com/bots-garden/simplism/releases/download/v0.1.3/deploy-wasm-from-remote.yaml
```

There are some Simplism plug-ins available here for demonstration purposes:
- https://github.com/simplism-registry/small-cow/releases/tag/v0.0.0
- https://github.com/simplism-registry/small-ant/releases/tag/v0.0.0


## Deploy a first function

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

## Let's try another function

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

