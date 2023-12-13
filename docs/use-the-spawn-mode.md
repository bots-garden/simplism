# Spawn mode: start Simplism process remotely with an API

You can use a Simplism service to spawn new Simplism services remotely.

## Start the service with the "spawn" mode activated

To activate the "spawn" mode, use the `--spawn-mode` flag and the `--admin-spawn-token` flag:

```bash
rm service-discovery/*.db
simplism listen \
../service-discovery/service-discovery.wasm handle \
--http-port 9000 \
--log-level info \
--service-discovery true \
--admin-discovery-token people-are-strange \
--information "👋 I'm the spawner service" \
--spawn-mode true \
--admin-spawn-token michael-burnham-rocks
```

Now, you can use the `http://localhost:9000/spawn` endpoint to create a new Simplism service.

## Spawn new Simplism services with the API

```bash
curl -X POST \
http://localhost:9000/spawn \
-H 'admin-spawn-token:michael-burnham-rocks' \
-H 'Content-Type: application/json; charset=utf-8' \
--data-binary @- << EOF
{
    "wasm-file":"../hello-people/hello-people.wasm", 
    "wasm-function":"handle", 
    "http-port":"9091", 
    "discovery-endpoint":"http://localhost:9000/discovery", 
    "admin-discovery-token":"people-are-strange",
    "admin-spawn-token":"michael-burnham-rocks",
    "information": "✋ I'm the hello-people service",
    "service-name": "hello-people"
}
EOF

curl -X POST \
http://localhost:9000/spawn \
-H 'admin-spawn-token:michael-burnham-rocks' \
-H 'Content-Type: application/json; charset=utf-8' \
--data-binary @- << EOF
{
    "wasm-file":"../hello/hello.wasm", 
    "wasm-function":"handle", 
    "http-port":"9092", 
    "discovery-endpoint":"http://localhost:9000/discovery", 
    "admin-discovery-token":"people-are-strange",
    "admin-spawn-token":"michael-burnham-rocks",
    "information": "✋ I'm the hello service",
    "service-name": "hello"
}
EOF
```

## Get the service list

```bash
curl http://localhost:9000/discovery \
-H 'admin-discovery-token:people-are-strange'
```

## Query the services by service name through the discovery service

```bash
curl http://localhost:9000/service/hello-people \
-H 'content-type: application/json; charset=utf-8' \
-d '{"firstName":"Bob","lastName":"Morane"}'

curl http://localhost:9000/service/hello \
-d 'Bob Morane'
```
