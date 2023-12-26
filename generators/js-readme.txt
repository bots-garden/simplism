# <name>

> Install Extism-js
```bash
EXTISM_JS_VERSION=1.0.0-rc3
EXTISM_JS_ARCH=x86_64
EXTISM_JS_OS=linux

curl -L -O "https://github.com/extism/js-pdk/releases/download/v${EXTISM_JS_VERSION}/extism-js-${EXTISM_JS_ARCH}-${EXTISM_JS_OS}-v${EXTISM_JS_VERSION}.gz"
gunzip extism-js*.gz
chmod +x extism-js-*
sudo mv extism-js-* /usr/local/bin/extism-js
```

> Install Binaryen
```bash
BINARYEN_VERSION=version_116
BINARYEN_ARCH=x86_64
BINARYEN_OS=linux

wget https://github.com/WebAssembly/binaryen/releases/download/${BINARYEN_VERSION}/binaryen-${BINARYEN_VERSION}-${BINARYEN_ARCH}-${BINARYEN_OS}.tar.gz
tar -xf binaryen-${BINARYEN_VERSION}-${BINARYEN_ARCH}-${BINARYEN_OS}.tar.gz

sudo cp binaryen-${BINARYEN_VERSION}/bin/* /usr/bin
rm -rf binaryen-${BINARYEN_VERSION}
rm binaryen-${BINARYEN_VERSION}-${BINARYEN_ARCH}-${BINARYEN_OS}.tar.gz

wasm2js --version
```
> There is no linux arm version of Binaryen

> Build the wasm plug-in:
```bash
extism-js index.js -i index.d.ts -o <name>.wasm
```

> Serve the wasm plug-in with Simplism:
```bash
simplism listen \
<name>.wasm handle --http-port 8080 --log-level info
```

> Query the wasm plug-in:
```bash
curl http://localhost:8080 \
-H 'content-type: text/plain; charset=utf-8' \
-d 'Bob Morane'
```

## Dockerize the service

### Build the Docker image

```bash
IMAGE_NAME="<name>-simplism"
docker login -u ${DOCKER_USER} -p ${DOCKER_PWD}
docker build -t ${IMAGE_NAME} . 

docker images | grep ${IMAGE_NAME}
```

### Create and run the container

```bash
IMAGE_NAME="<name>-simplism"
docker run \
  -p 8080:8080 \
  --rm ${IMAGE_NAME}
```

### Push the image to the Docker Hub

```bash
IMAGE_NAME="<name>-simplism"
docker tag ${IMAGE_NAME} ${DOCKER_USER}/${IMAGE_NAME}:0.0.0
docker push ${DOCKER_USER}/${IMAGE_NAME}:0.0.0
```