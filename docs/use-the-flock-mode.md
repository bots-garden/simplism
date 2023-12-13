# Start all the Slimplism services from a config file with the Flock mode

> prerequisite: read [Use the service discovery feature](use-service-discovery.md)

You can start several Slimplism services from a config file with the Flock mode. First, create a config file with the following content:

> config.yml
```yaml
service-discovery:
  wasm-file: ../service-discovery/service-discovery.wasm
  wasm-function: handle
  http-port: 9000
  log-level: info
  service-discovery: true
  admin-discovery-token: people-are-strange

hello-people:
  wasm-file: ../hello-people/hello-people.wasm
  wasm-function: handle
  http-port: 8081
  log-level: info
  service-name: hello-people
  admin-discovery-token: people-are-strange
  discovery-endpoint: http://localhost:9000/discovery

hello: 
  wasm-file: ../hello/hello.wasm
  wasm-function: handle
  http-port: 8082
  log-level: info
  service-name: hello
  admin-discovery-token: people-are-strange
  discovery-endpoint: http://localhost:9000/discovery
```

## Start all the Slimplism services

To start all the Slimplism services from a config file with the Flock mode, use the following command:

```bash
simplism flock ./config.yml
```

## Get the service list

```bash
curl http://localhost:9000/discovery \
-H 'admin-discovery-token:people-are-strange'
```

## Query the services by service name through the discovery service

As soon as service discovery is enabled, you can use it directly to query services by name:

```bash
curl http://localhost:9000/service/hello-people \
-H 'content-type: application/json; charset=utf-8' \
-d '{"firstName":"Bob","lastName":"Morane"}'

curl http://localhost:9000/service/hello \
-d 'Bob Morane'
```
