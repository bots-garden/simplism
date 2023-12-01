# Simplism: a tiny HTTP server for Extism Plug-ins

![image](imgs/simplism-small-logo.jpeg)

## What is Simplism?

**Simplism** is a tiny HTTP server to serve [Extism](https://extism.org/) WebAssembly plug-ins and execute/call a single WebAssembly function.

> It's like the official [Extism CLI](https://github.com/extism/cli), but **Simplism** is "serving" the Extism WebAssembly plug-in instead of running it, and call a function at every HTTP request.

## 🚀 Getting started

### Install Simplism

```bash
SIMPLISM_DISTRO="Linux_arm64.tar" # 👀 https://github.com/bots-garden/simplism/releases
VERSION="0.0.5"
wget https://github.com/bots-garden/simplism/releases/download/v${VERSION}/simplism_${SIMPLISM_DISTRO}.tar.gz -O simplism.tar.gz 
tar -xf simplism.tar.gz -C /usr/bin
rm simplism.tar.gz
simplism version
```

### Generate a (GoLang) wasm plug-in

```bash
simplism generate golang hello ./

# hello
# ├── go.mod
# ├── main.go
# └── README.md
```

#### Build the wasm plug-in
> you can follow the instructions into the `hello/README.md` file

```bash
cd hello
tinygo build -scheduler=none --no-debug \
-o hello.wasm \
-target wasi main.go
```

#### Serve the wasm plug-in

```bash
simplism listen \
hello.wasm handle --http-port 8080 --log-level info
```

#### Query the wasm plug-in:

```bash
curl http://localhost:8080/hello/world \
-H 'content-type: application/json; charset=utf-8' \
-d '{"firstName":"Bob","lastName":"Morane"}'
```

## Run Simplism

```text
Usage:
  simplism [command] [arguments]

Available Commands:
  listen      Serve an Extism plug-in function
              Arguments: [wasm file path] [function name]
  version     Display the Minism version
              Arguments: nothing
  generate    Generate a source code project of an Extism plug-in
              Arguments: [plug-in language] [project name] [project path]
  config      Serve an Extism plug-in function using a yaml configuration file
              Arguments: [yaml file path] [config key]
  flock       Serve several Extism plug-in functions using a yaml configuration file
              Arguments: [yaml file path] [config key]

Flags for listen command:
  --http-port             string   HTTP port of the Simplism server (default: 8080)
  --log-level             string   Log level to print message
                                   Possible values: error, warn, info, debug, trace
  --allow-hosts           string   Hosts for HTTP request (json array) 
                                   Default: ["*"]
  --allow-paths           string   Allowed paths to write and read files (json string) 
                                   Default: {}
  --config                string   Configuration data (json string)
                                   Default: {}
  --env                   string   Environment variables to forward to the wasm plug-in
                                   Default: []
  --wasi                  bool     Default: true
  --wasm-url              string   Url to download the wasm file
  --wasm-url-auth-header  string   Authentication header to download the wasm file, ex: "PRIVATE-TOKEN=IlovePandas"
                                   Or use this environment variable: WASM_URL_AUTH_HEADER='PRIVATE-TOKEN=IlovePandas'
  --cert-file             string   Path to certificate file (https)
  --key-file              string   Path to key file (https)
  --admin-reload-token    string   Admin token to be authorized to reload the wasm-plugin
                                   Or use this environment variable: ADMIN_RELOAD_TOKEN
                                   Use the /reload endpoint to reload the wasm-plugin
  --discovery             bool     The current Simplism server is a service discovery server
                                   Default: false
  --discovery-endpoint    string   The endpoint of the service discovery server
                                   It always ends with /discovery
                                   Example: http://localhost:9000/discovery
  --admin-discovery-token string   Admin token to be authorized to post information to the service discovery server
                                   Or use this environment variable: ADMIN_DISCOVERY_TOKEN
                                   Use the /discovery endpoint to post information to the service discovery server
```
> *Remarks: look at the `./samples` directory*

> **Examples**:

```bash
simplism listen ./samples/golang/simple-plugin/simple.wasm say_hello
```

```bash
simplism listen ./samples/golang/hello-plugin/simple.wasm say_hello \
--http-port 9090 \
--log-level info \
--allow-hosts '["*","*.google.com"]' \
--config '{"message":"👋 hello world 🌍"}' \
--allow-paths '{"data":"/mnt"}'
```

> **Configuration example**:

```yaml
# config.yml
hello-plugin:
  wasm-file: ./hello.wasm
  wasm-function: say_hello
  http-port: 8080
  log-level: info
```

Run the server like this: `simplism config ./config.yml hello-plugin`

> **Run Simplism in "flock" mode**:

```yaml
# config.yml
hello-1:
  wasm-file: ./hello.wasm
  wasm-function: say_hello
  http-port: 8081
  log-level: info
hello-2:
  wasm-file: ./hello.wasm
  wasm-function: say_hello
  http-port: 8082
  log-level: info
hello-3:
  wasm-file: ./hello.wasm
  wasm-function: say_hello
  http-port: 8083
  log-level: info
```

Run the server**s** like this: `simplism flock ./config.yml`. It will start **3** instances of Simplism.

> See `samples/flock` repository for a more complex example.


## Reload remotely a wasm plug-in without stopping the Simplism server

### Start the Simplism server

```bash
simplism listen ./hey-one.wasm handle --http-port 8080  --admin-reload-token "1234567890"
```

or

```bash
export ADMIN_RELOAD_TOKEN="1234567890"
simplism listen ./hey-one.wasm handle --http-port 8080
```

### Reload the wasm plug-in with the /reload api

```bash
curl -v -X POST \
http://localhost:8080/reload \
-H 'content-type: application/json; charset=utf-8' \
-H 'admin-reload-token:1234567890' \
-d '{"wasm-url":"http://0.0.0.0:3333/hey-two/hey-two.wasm", "wasm-file": "./hey-two.wasm", "wasm-function": "handle"}'
```

## Service discovery

> 🚧 this is a work in progress


## Generate Extism plug-in projects for Simplism

You can use **Simplism** to generate a project skeleton of an **Extism** plug-in with the following languages:

- Golang
- Rustlang

### Generate a Golang project

```bash
simplism generate golang hello my-projects
```
This command will create this tree structure:
```bash
├── my-projects
│  ├── hello
│  │  ├── go.mod
│  │  ├── main.go
│  │  └── README.md
```

### Generate a Rustlang project

```bash
simplism generate rustlang hello my-projects
```
This command will create this tree structure:
```bash
├── my-projects
│  ├── hello
│  │  ├── src
│  │  │  └── lib.rs
│  │  ├── Cargo.toml
│  │  └── README.md
```

✋ more languages to come very soon

## How is Simplism developed?

Simplism is developed in Go with **[Wazero](https://wazero.io/)**[^1] as the Wasm runtime and **[Extism](https://extism.org/)**[^2], which offers a Wazero-based Go SDK and a Wasm plugin system.

### Prerequisites
> 🚧 work in progress

To develop on the Simplism project and/or create Extism plug-ins, look at `docker-dev-environment/Dockerfile`, you will find the list of the necessary softwares, libraries, tools...

### 👋 Or you can use ready to use environments

- [🍊 Open it with Gitpod](https://gitpod.io/#https://github.com/bots-garden/simplism)
- [🐳 Open it with Docker Dev Environment (**✋ arm version**)](https://open.docker.com/dashboard/dev-envs?url=https://github.com/bots-garden/simplism/tree/main)
  - Prerequisites:
    - https://docs.docker.com/desktop/dev-environments/create-dev-env/#prerequisites
    - [Visual Studio Code Remote Containers Extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)


### Build Simplism

```bash
go build
./simplism version
```

### Write an Extism plug-in

- Let's have a look at the official Extism documentation https://extism.org/docs/category/write-a-plug-in 
- Look into the `samples` directory of this repository:
  ```bash
  samples
  ├── golang
  │  ├── hello-plugin
  │  └── simple-plugin
  └── rustlang
    ├── hello-plugin
    └── simple-plugin
  ```

> ✋ **important**: you can write Extism plug-ins with Go, Rust, AssemblyScript, Zig, C, Haskell and JavaScript


[^1]: Wazero is a project from **[Tetrate](https://tetrate.io/)**
[^2]: Extism is a project from **[Dylibso](https://dylibso.com/)**
