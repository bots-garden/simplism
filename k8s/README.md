# Simplism on Kubernetes
> Create a Simply FaaS on Kubernetes with Simplism

## Prerequisites

### Kubernetes

You need a Kubernetes Cluster. I use [Civo](https://www.civo.com/), a Kubernetes as a Service provider. Of course you can use any Kubernetes provider or do it locally (I will do another blog post on this topic). I wrote two blog posts that can help you if you are a beginner with Kube:

- [Create easily your first K8S cluster on Civo.com](https://k33g.hashnode.dev/create-easily-your-first-k8s-cluster-on-civocom-and-do-it-again-and-again-and)
- [Simple Golang Web Application Deployment to K8S on Civo.com](https://k33g.hashnode.dev/simple-golang-web-application-deployment-to-k8s-on-civocom)


You need some tools:
- Kubectl: https://kubernetes.io/docs/tasks/tools/#kubectl
- K9S (it's not mandatory, but it's really helpful to manage a kubernetes cluster): https://k9scli.io/
- `envsubst` to substitute the values of environment variables in a file, installation: https://command-not-found.com/envsubst

### Did you know Simplism?
> If you don't want to read the documentation
I wrote a blog post series: [Simplism: the cloud-native runtime for Extism Wasm plug-ins](https://k33g.hashnode.dev/series/simplism)


## Connect to the cluster

Once you get a Kubernetes cluster, define the value of `KUBECONFIG` and you can connect to it with the following command:

```bash
export KUBECONFIG=$PWD/config/k3s.yaml
kubectl get pods --all-namespaces
```

You should get something like this:

```bash

NAMESPACE     NAME                                 READY   STATUS      RESTARTS   AGE
kube-system   civo-ccm-db67548d-cnxkk              1/1     Running     0          3m43s
kube-system   civo-csi-node-xnc4f                  2/2     Running     0          3m33s
kube-system   coredns-59b4f5bbd5-lxbkw             1/1     Running     0          3m43s
kube-system   metrics-server-7b67f64457-xmngb      1/1     Running     0          3m43s
kube-system   civo-csi-controller-0                4/4     Running     0          3m44s
default       install-traefik2-nodeport-wo-6vkwl   0/1     Completed   0          3m19s
kube-system   traefik-f2hpb                        1/1     Running     0          3m12s
```

If you installed K9S, you can use it with the following command:

```bash
export KUBECONFIG=$PWD/config/k3s.yaml
k9s --all-namespaces
```

For the rest, create a `.env` vile with this content:
```bash
KUBECONFIG=$PWD/config/k3s.yaml
KUBE_NAMESPACE="simplism-faas"
DNS="ffb140f9-7479-4308-9763-9f70628794b1.k8s.civo.com"
IMAGE_NAME="k33g/simplism"
IMAGE_TAG="0.1.3"
```
> - `DNS` is the domain name of your K8S cluster (I'm using Civo.com), set the `DNS` variable with the information of your own cluster.
> - `KUBE_NAMESPACE` is the namespace we will use in our cluster.

## Chapters

- [Deploy a Simplism function on K8S](01-wasm-function-in-a-pod.md)
- [Deploy a Simplism registry on K8S](02-wasm-registry-with-two-pods.md)