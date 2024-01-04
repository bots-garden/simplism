# Use the registry mode

## Reminder

You can start a remote wasm plugin with the `wasm-url` flag. At start, the Simplism server will download the remote wasm file and then, will serve it.

For example:
  - You can publish a wasm plugin into the [generic package registry of GitLab](https://docs.gitlab.com/ee/user/packages/generic_packages/), 
  - Or as a [release asset on GitHub](https://docs.github.com/en/rest/releases/assets),
  - Or use a simple HTTP server to serve the wasm file: `python3 -m http.server 3333 --directory ./wasm-files`
  - ...

## Wasm registry mode

It's possible to use the `--registry-mode` flag of the listen command to use a Simplism service as a kind of wasm registry:

```shell
simplism listen ./tiny-registry.wasm handle \
  --http-port 9090 \
  --log-level info \
  --registry-mode true \
  --registry-path ./wasm-files
```

- The `--registry-mode` flag is used to enable the wasm registry.
- The `--registry-path` flag is the path where the wasm files will be uploaded.

The "registry mode" activate 3 endpoints:
  - `/registry/push` to upload a wasm file
  - `/registry/pull` to download a wasm file
  - `/registry/remove` to remove a wasm file
  - `/registry/discover` to discover wasm files

> If you don't want to specify a wasm file, use this:

```bash
simplism listen ? ? \
  --http-port 9090 \
  --log-level info \
  --registry-mode true \
  --registry-path ./wasm-files
```

or:

```yaml
tiny-registry-config:
  wasm-file: "?"
  wasm-function: "?"
  http-port: 9090
  log-level: info
  registry-mode: true
  registry-path: ./wasm-files
```

### Upload a wasm file to the registry

```bash
curl http://localhost:9090/registry/push \
-F 'file=@hello.wasm'
```

### Download a wasm file from the registry

```bash
curl http://localhost:9090/registry/pull/hello.wasm -o hello.wasm
```

### Get the list of the wasm files

```bash
curl http://localhost:9090/registry/discover
```

### Remove a wasm file

```bash
curl http://localhost:9090/registry/remove/hello.wasm
```

### Start a remote wasm plugin

```bash
simplism listen ./hello.wasm handle \
  --http-port 8080 \
  --log-level info \
  --wasm-url http://localhost:9090/registry/pull/hello.wasm
```

## Protect the registry

You can protect the registry by using these two flags:
- `--admin-registry-token` for `/registry/push` and `/registry/remove` endpoints
- `--private-registry-token` for `/registry/pull` and `/registry/discover` endpoints

```shell
simplism listen ./tiny-registry.wasm handle \
  --http-port 9090 \
  --log-level info \
  --registry-mode true \
  --registry-path ./wasm-files \
  --admin-registry-token: morrison-hotel \
  --private-registry-token: people-are-strange
```

or:

```bash
export PRIVATE_REGISTRY_TOKEN=people-are-strange
export ADMIN_REGISTRY_TOKEN=morrison-hotel

simplism listen ./tiny-registry.wasm handle \
  --http-port 9090 \
  --log-level info \
  --registry-mode true \
  --registry-path ./wasm-files
```

### Upload a wasm file to the registry

```bash
curl http://localhost:9090/registry/push \
-H 'admin-registry-token: morrison-hotel' \
-F 'file=@hello.wasm'
```

### Download a wasm file from the registry

```bash
curl http://localhost:9090/registry/pull/hello.wasm -o hello.wasm \
-H 'private-registry-token: people-are-strange'
```

### Get the list of the wasm files

```bash
curl http://localhost:9090/registry/discover \
-H 'private-registry-token: people-are-strange'
```

### Remove a wasm file

```bash
curl http://localhost:9090/registry/remove/hello.wasm \
-H 'admin-registry-token: morrison-hotel'
```

### Start a remote wasm plugin

```bash
simplism listen ./hello.wasm handle \
--http-port 8080 \
--log-level info \
--wasm-url http://localhost:9090/registry/pull/hello.wasm \
--wasm-url-auth-header: private-registry-token=people-are-strange
```

