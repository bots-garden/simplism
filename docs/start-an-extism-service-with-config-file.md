# Start an Simplism service from a configuration file

You can use a configuration file to start an Extism service, and you can define multiple configurations:

> config.yml
```yaml
config-one:
  wasm-file: ./hello-people.wasm
  wasm-function: handle
  http-port: 8080

config-two:
  wasm-file: ./hello-people.wasm
  wasm-function: handle
  http-port: 9090
```

> Start one or more configurations:
```shell
simplism config ./config.yml config-one
simplism config ./config.yml config-two
```
