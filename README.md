# Simplism: a tiny HTTP server for Extism Plug-ins

![image](imgs/simplism-small-logo.jpeg)

## What is Simplism?

**Simplism** is a tiny HTTP server to serve [Extism](https://extism.org/) WebAssembly plug-ins and execute/call a single WebAssembly function.

> It's like the official [Extism CLI](https://github.com/extism/cli), but **Simplism** is "serving" the Extism WebAssembly plugin instead of running it, and call a function at every HTPP request.

## How is Simplism developed?

Simplism is developed in Go with **[Wazero](https://wazero.io/)**[^1] as the Wasm runtime and **[Extism](https://extism.org/)**[^2], which offers a Wazero-based Go SDK and a Wasm plugin system.

### Write an Extism plug-in

Look at the official Extism documentation https://extism.org/docs/category/write-a-plug-in (and into the `samples` directory of this repository).

> ✋ **important**: you can write Extism plug-ins with Go, Rust, AssemblyScript, Zig, C, Haskell and JavaScript

### Ready to use environments

- [🍊 Open it with Gitpod](https://gitpod.io/#https://github.com/bots-garden/simplism)
- [🐳 Open it with Docker Dev Environment (**✋ arm version**)](https://open.docker.com/dashboard/dev-envs?url=https://github.com/bots-garden/simplism/tree/main)
  - Prerequisites:
    - https://docs.docker.com/desktop/dev-environments/create-dev-env/#prerequisites
    - [Visual Studio Code Remote Containers Extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)

## Build Simplism

```bash
go build
./simplism version
```

## Run Simplism

```text
Usage:
  simplism [command] [arguments]

Available Commands:
  listen      Serve an Extism plugin function
              Arguments: [wasm file path] [function name]
  version     Display the Minism version
              Arguments: nothing

Flags:
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
./simplism listen ./samples/golang/simple-plugin/simple.wasm say_hello
```

```bash
./simplism listen ./samples/golang/hello-plugin/simple.wasm say_hello \
  --http-port 9090 \
  --log-level info \
  --allow-hosts '["*","*.google.com"]' \
  --config '{"message":"👋 hello world 🌍"}' \
  --allow-paths '{"data":"/mnt"}'
```


[^1]: Wazero is a project from **[Tetrate](https://tetrate.io/)**
[^2]: Extism is a project from **[Dylibso](https://dylibso.com/)**