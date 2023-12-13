# Dockerize a Simplism service

> prerequisite: read [create-and-serve-wasm-plug-in.md](create-and-serve-wasm-plug-in.md)

## Create a Dockerfile

Create a `Dockerfile` at the root of the project
```Dockerfile
FROM k33g/simplism:0.0.7
COPY hello.wasm .
EXPOSE 8080
CMD ["/simplism", "listen", "hello.wasm", "handle", "--http-port", "8080", "--log-level", "info"]
```

## Build the image

```bash
IMAGE_NAME="hello-simplism"
docker login -u ${DOCKER_USER} -p ${DOCKER_PWD}
docker build -t ${IMAGE_NAME} . 

docker images | grep ${IMAGE_NAME}
```

You should get something like that:
```bash
hello-simplism   latest            e96bc4f2511f   5 seconds ago   8.8MB
```

## Create and run the container

```bash
IMAGE_NAME="hello-simplism"
docker run \
  -p 8080:8080 \
  --rm ${IMAGE_NAME}
```

## Call the service

```bash
curl http://localhost:8080 \
-d 'Bob Morane'
```

## Push the image to the Docker Hub

```bash
IMAGE_NAME="hello-simplism"
docker tag ${IMAGE_NAME} ${DOCKER_USER}/${IMAGE_NAME}:0.0.0
docker push ${DOCKER_USER}/${IMAGE_NAME}:0.0.0
```
