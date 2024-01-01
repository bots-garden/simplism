# hello-world

> Add Extism dependency
```bash
go get github.com/extism/go-pdk
```

> Build the wasm plug-in:
```bash
tinygo build -scheduler=none --no-debug \
-o hello-world.wasm \
-target wasi main.go
```

> Serve the wasm plug-in with Simplism:
```bash
simplism listen \
hello-world.wasm handle --http-port 8080 --log-level info
```

> Query the wasm plug-in:
```bash
curl http://localhost:8080/hello/world \
-H 'content-type: application/json; charset=utf-8' \
-d '{"firstName":"Bob","lastName":"Morane"}'
```

## Dockerize the service

### Build the Docker image

```bash
IMAGE_NAME="hello-world-simplism"
docker login -u ${DOCKER_USER} -p ${DOCKER_PWD}
docker build -t ${IMAGE_NAME} . 

docker images | grep ${IMAGE_NAME}
```

### Create and run the container

```bash
IMAGE_NAME="hello-world-simplism"
docker run \
  -p 8080:8080 \
  --rm ${IMAGE_NAME}
```

### Push the image to the Docker Hub

```bash
IMAGE_NAME="hello-world-simplism"
docker tag ${IMAGE_NAME} ${DOCKER_USER}/${IMAGE_NAME}:0.0.0
docker push ${DOCKER_USER}/${IMAGE_NAME}:0.0.0
```
