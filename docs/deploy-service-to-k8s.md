# Deploy a Simplism service to Kubernetes

> prerequisites: 
> - create a Simplism service: read [create-and-serve-wasm-plug-in.md](create-and-serve-wasm-plug-in.md)
> - dockerize your application: read [dockerize-a-simplism-service.md](dockerize-a-simplism-service.md)
> remark: the tests have been done on https://civo.com

## Install kubectl

```bash
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"

chmod +x kubectl
mkdir -p ~/.local/bin
mv ./kubectl ~/.local/bin/kubectl
```

## Kubernetes configuration

> Set the KUBECONFIG environment variable:
```bash
# for example
export KUBECONFIG=$PWD/k3s.yaml
# check
kubectl get pods
```

## Create a yaml manifest to deploy the service

Create a `deploy.yaml` file at the root of the project:
```yaml
---
# Service
apiVersion: v1
kind: Service
metadata:
  name: demo-hello-simplism
spec:
  selector:
    app: demo-hello-simplism
  ports:
    - port: 80
      targetPort: 8080
---
# Deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: demo-hello-simplism
spec:
  replicas: 1
  selector:
    matchLabels:
      app: demo-hello-simplism
  template:
    metadata:
      labels:
        app: demo-hello-simplism
    spec:
      containers:
        - name: demo-hello-simplism
          image: k33g/hello-simplism:0.0.0
          ports:
            - containerPort: 8080
          imagePullPolicy: Always

---
# Ingress
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: demo-hello-simplism
spec:
  rules:
    - host: demo-hello-simplism.89ff64fe-3673-4ce6-82f7-9bc17793af6d.k8s.civo.com
      http:
        paths:
        - path: /
          pathType: Prefix
          backend:
            service: 
              name: demo-hello-simplism
              port: 
                number: 80
```

### Deploy the service

```bash
kubectl apply -f ./deploy.yaml
```

### Get the ingress

```bash
WEBAPP_NAME="demo-hello-simplism"
KUBE_NAMESPACE="default"
kubectl describe ingress ${WEBAPP_NAME} -n ${KUBE_NAMESPACE}
```

You should get something like this:
```bash
Name:             demo-hello-simplism
Labels:           <none>
Namespace:        default
Address:          89ff64fe-3673-4ce6-82f7-9bc17793af6d.k8s.civo.com
Ingress Class:    <none>
Default backend:  <default>
Rules:
  Host                                                                   Path  Backends
  ----                                                                   ----  --------
  demo-hello-simplism.89ff64fe-3673-4ce6-82f7-9bc17793af6d.k8s.civo.com  
                                                                         /   demo-hello-simplism:80 (10.42.0.7:8080)
Annotations:                                                             <none>
Events:                                                                  <none>
```

### Call the Simplism service

```bash
curl http://demo-hello-simplism.89ff64fe-3673-4ce6-82f7-9bc17793af6d.k8s.civo.com -d 'Bob Morane'
```

You should get: `🤗 Hello Bob Morane`

