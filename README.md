# Simplism: a tiny HTTP server for Extism Plug-ins

![image](imgs/simplism-small-logo.jpeg)

## What is Simplism?

**Simplism** is a tiny HTTP server to serve [Extism](https://extism.org/) WebAssembly plug-ins and execute/call a single WebAssembly function.

> It's like the official [Extism CLI](https://github.com/extism/cli), but **Simplism** is "serving" the Extism WebAssembly plugin instead of running it, and call a function at every HTPP request.

## 🚀 Getting started

### Install Simplism

```bash
SIMPLISM_DISTRO="Linux_arm64.tar" # 👀 https://github.com/bots-garden/simplism/releases
VERSION="0.0.0"
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
              
Flags for listen command:
  --http-port         string   HTTP port of the Simplism server (default: 8080)
  --log-level         string   Log level to print message
                               Possible values: error, warn, info, debug, trace
  --allow-hosts       string   Hosts for HTTP request (json array) 
                               Default: ["*"]
  --allow-paths       string   Allowed paths to write and read files (json string) 
                               Default: {}
  --config            string   Configuration data (json string)
                               Default: {}
  --wasi              bool     Default: true
  --wasm-url          string   Url to download the wasm file
  --auth-header-name  string   Authentication header name, ex: PRIVATE-TOKEN
  --auth-header-value string   Value of the authentication header, ex: IlovePandas  
```

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

