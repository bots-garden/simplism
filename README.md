# Simplism: a tiny HTTP server for Extism Plug-ins

![image](imgs/simplism-small-logo.jpeg)

## What is Simplism?

**Simplism** is a tiny HTTP server to serve [Extism](https://extism.org/) WebAssembly plug-ins and execute/call a single WebAssembly function.

> It's like the official [Extism CLI](https://github.com/extism/cli), but **Simplism** is "serving" the Extism WebAssembly plug-in instead of running it, and call a function at every HTTP request.

## 🚀 Getting started

### Install Simplism

```bash
SIMPLISM_DISTRO="Linux_arm64" # 👀 https://github.com/bots-garden/simplism/releases
VERSION="0.1.0"
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
              Languages: golang (or go), rustlang (or rust), javascript (or js)
              Ex: simplism generate js hello samples/javascript (it will create samples/javascript/hello/)
  config      Serve an Extism plug-in function using a yaml configuration file
              Arguments: [yaml file path] [config key]
  flock       Serve several Extism plug-in functions using a yaml configuration file
              Arguments: [yaml file path] [config key]

Flags for listen command:
  --http-port              int      HTTP port of the Simplism server (default: 8080)
  --log-level              string   Log level to print message
                                    Possible values: error, warn, info, debug, trace
  --allow-hosts            string   Hosts for HTTP request (json array) 
                                    Default: ["*"]
  --allow-paths            string   Allowed paths to write and read files (json string) 
                                    Default: {}
  --config                 string   Configuration data (json string)
                                    Default: {}
  --env                    string   Environment variables to forward to the wasm plug-in
                                    Default: []
  --wasi                   bool     Default: true
  --wasm-url               string   Url to download the wasm file
  --wasm-url-auth-header   string   Authentication header to download the wasm file, ex: "PRIVATE-TOKEN=IlovePandas"
                                    Or use this environment variable: WASM_URL_AUTH_HEADER='PRIVATE-TOKEN=IlovePandas'
  --cert-file              string   Path to certificate file (https)
  --key-file               string   Path to key file (https)
  --admin-reload-token     string   Admin token to be authorized to reload the wasm-plugin
                                    Or use this environment variable: ADMIN_RELOAD_TOKEN
                                    Use the /reload endpoint to reload the wasm-plugin
  --service-discovery      bool     The current Simplism server is a service discovery server
                                    Default: false
  --discovery-endpoint     string   The endpoint of the service discovery server
                                    It always ends with /discovery
                                    Example: http://localhost:9000/discovery
  --admin-discovery-token  string   Admin token to be authorized to post information to the service discovery server
                                    Or use this environment variable: ADMIN_DISCOVERY_TOKEN
                                    Use the /discovery endpoint to post information to the service discovery server
  --service-name           string   Name of the service (it can be useful with the service discovery mode)
  --information            string   Information about the service (it can be useful with the service discovery mode)
  --spawn-mode             bool     The current Simplism server is in spawn mode (it can create new simplism servers with the /spawn endpoint)
                                    Default: false
  --admin-spawn-token      string   Admin token to be authorized to spawn a new Simplism server
                                    Or use this environment variable: ADMIN_SPAWN_TOKEN
                                    Use the /spawn endpoint to spawn a new Simplism server
  --http-port-auto         bool     Automatically find an available port (only works in spawn mode)
                                    Default: false
  --recovery-path          string   Path of the recovery file (only works in spawn mode)
                                    Default: "recovery.yaml"
  --recovery-mode          bool     The current Simplism server is in recovery mode
                                    Default: true
  --store-mode             bool     The current Simplism server exposes a store api to save data in a bbolt database
                                    Use the /store endpoint (see documentation)
                                    Default: false
  --store-path             string   File path of the store db file  
                                    Default: file path of the wasm file + file name + ".store.db"
  --admin-store-token      string   Admin token to be authorized to use the store API of a Simplism server
                                    Or use this environment variable: ADMIN_STORE_TOKEN
  --registry-mode          bool     The current Simplism server exposes a registry api to upload wasm files
                                    Use the /registry endpoint (see documentation)
                                    Default: false
  --registry-path          string   File path of the uploaded wasm files  
  --admin-registry-token   string   Admin token to be authorized to use the registry API: POST(`/push`) and DELETE(`/remove`)
                                    Or use this environment variable: ADMIN_REGISTRY_TOKEN
  --private-registry-token string   Private registry token to be authorized to use the registry API: GET(`/pull`) and GET(`/discover`)
                                    Or use this environment variable: PRIVATE_REGISTRY_TOKEN
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

Simplism comes with a service discovery feature. It can be used to discover the running Simplism servers.
- One of the servers (simplism service) can be a service discovery server. The service discovery server can be configured with the `--service-discovery` flag:

```bash
simplism listen discovery-service/discovery-service.wasm handle \
--http-port 9000 \
--log-level info \
--service-discovery true \
--admin-discovery-token people-are-strange
```
> `--admin-discovery-token` is not mandatory, but it's probably a good idea to set it.

- Then, the other services can be configured with the `--discovery-endpoint` flag:

```bash
simplism listen service-one/service-one.wasm handle \
--http-port 8001 \
--log-level info \
--discovery-endpoint http://localhost:9000/discovery \
--admin-discovery-token people-are-strange &

simplism listen service-two/service-two.wasm handle \
--http-port 8002 \
--log-level info \
--discovery-endpoint http://localhost:9000/discovery \
--admin-discovery-token people-are-strange &

simplism listen service-three/service-three.wasm handle \
--http-port 8003 \
--log-level info \
--discovery-endpoint http://localhost:9000/discovery \
--admin-discovery-token people-are-strange &
```
> the 3 services will be discovered by the service discovery server. Every services will regularly post information to the service discovery server.

- You can query the service discovery server with the `/discovery` endpoint to get the list of the running services:

```bash
curl http://localhost:9000/discovery \
-H 'admin-discovery-token:people-are-strange'
```

- You can use the flock mode jointly with the service discovery:

```yaml
service-discovery:
  wasm-file: ./discovery/discovery.wasm
  wasm-function: handle
  http-port: 9000
  log-level: info
  service-discovery: true
  admin-discovery-token: this-is-the-way

basestar-mother:
  wasm-file: ./basestar/basestar.wasm
  wasm-function: handle
  http-port: 8010
  log-level: info
  discovery-endpoint: http://localhost:9000/discovery
  admin-discovery-token: this-is-the-way

raider-1:
  wasm-file: ./raider/raider.wasm
  wasm-function: handle
  http-port: 8001
  log-level: info
  discovery-endpoint: http://localhost:9000/discovery
  admin-discovery-token: this-is-the-way
```

## Spawn mode

> 🚧 this is a work in progress

If you activate the `--spawn-mode` flag, the Simplism server will be able tospawn a new Simplism server with the `/spawn` endpoint:

```bash
simplism listen ./process-spawner.wasm handle \
--http-port 8000 \
--log-level info \
--spawn-mode true \
--admin-spawn-token michael-burnham-rocks
```

Then, to "spawn" a new Simplism server process, you can use the `/spawn` endpoint with a simple curl request:

```bash
curl -X POST \
http://localhost:8080/spawn \
-H 'admin-spawn-token:michael-burnham-rocks' \
-H 'Content-Type: application/json; charset=utf-8' \
--data-binary @- << EOF
{
    "wasm-file":"../say-hello/say-hello.wasm", 
    "wasm-function":"handle", 
    "http-port":"9093", 
    "discovery-endpoint":"http://localhost:8080/discovery", 
    "admin-discovery-token":"michael-burnham-rocks"
}
EOF
echo ""
```

## Expose and use the "store API"

> start a Simplism server with the `--store-mode` flag:
```bash
simplism listen \
store.wasm handle \
--http-port 8080 \
--log-level info \
--store-mode true \
--admin-store-token morrison-hotel \
--information "👋 I'm the store service"
```

### Query the "store API"

> add records to the store:
```bash
curl http://localhost:8080/store \
-H 'content-type: application/json; charset=utf-8' \
-H 'admin-store-token: morrison-hotel' \
-d '{"key":"hello","value":"hello world"}'

curl http://localhost:8080/store \
-H 'content-type: application/json; charset=utf-8' \
-H 'admin-store-token: morrison-hotel' \
-d '{"key":"hey","value":"hey people"}'

curl http://localhost:8080/store \
-H 'content-type: application/json; charset=utf-8' \
-H 'admin-store-token: morrison-hotel' \
-d '{"key":"001","value":"first"}'

curl http://localhost:8080/store \
-H 'content-type: application/json; charset=utf-8' \
-H 'admin-store-token: morrison-hotel' \
-d '{"key":"002","value":"second"}'

curl http://localhost:8080/store \
-H 'content-type: application/json; charset=utf-8' \
-H 'admin-store-token: morrison-hotel' \
-d '{"key":"003","value":"third"}'
```

> get all records from the store:
```bash
curl http://localhost:8080/store \
-H 'admin-store-token: morrison-hotel'
```

> get a specific record from the store:
```bash
curl http://localhost:8080/store?key=hey \
-H 'admin-store-token: morrison-hotel'
```

> get all records from the store with a key prefix:
```bash
curl http://localhost:8080/store?prefix=00 \
-H 'admin-store-token: morrison-hotel'
```

> delete a specific record from the store:
```bash
curl -X "DELETE" http://localhost:8080/store?key=002 \
-H 'admin-store-token: morrison-hotel'
```


## Generate Extism plug-in projects for Simplism

You can use **Simplism** to generate a project skeleton of an **Extism** plug-in with the following languages:

- Golang
- Rustlang
- JavaScript

### Generate a Golang project

```bash
simplism generate golang hello my-projects
```
This command will create this tree structure:
```bash
my-projects
└── hello
   ├── build.sh
   ├── Dockerfile
   ├── go.mod
   ├── main.go
   ├── query.sh
   ├── README.md
   └── run.sh
```

### Generate a Rustlang project

```bash
simplism generate rustlang hello my-projects
```
This command will create this tree structure:
```bash
my-projects
└── hello
   ├── build.sh
   ├── Cargo.toml
   ├── Dockerfile
   ├── query.sh
   ├── README.md
   ├── run.sh
   └── src
      └── lib.rs
```

### Generate a JavaScript project

```bash
simplism generate js hello my-projects
```
This command will create this tree structure:
```bash
my-projects
└── hello
   ├── build.sh
   ├── Dockerfile
   ├── index.d.ts
   ├── index.js
   ├── query.sh
   ├── README.md
   └── run.sh
```

✋ more languages to come

## How is Simplism developed?

Simplism is developed in Go with **[Wazero](https://wazero.io/)**[^1] as the Wasm runtime and **[Extism](https://extism.org/)**[^2], which offers a Wazero-based Go SDK and a Wasm plugin system.

### Prerequisites
> 🚧 work in progress

To develop on the Simplism project and/or create Extism plug-ins, look at `.docker/compose/Dockerfile`, you will find the list of the necessary softwares, libraries, tools...

### 👋 Or you can use ready to use environments

[🍊 Open it with Gitpod](https://gitpod.io/#https://github.com/bots-garden/simplism)

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
