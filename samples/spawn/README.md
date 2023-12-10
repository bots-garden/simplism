# Spawn mode

## Build the Simplism "Spawner" service

```bash
cd process-spawner
tinygo build -scheduler=none --no-debug \
-o process-spawner.wasm \
-target wasi main.go
```

## Run the Simplism "Spawner" service

```bash
cd process-spawner
rm process-spawner.wasm.db
simplism config ./config.yml process-spawner-01
```

You should get:
```bash
🤖 this service is a service discovery
🔎 discovery mode activated: /discovery  ( 8080 )
🚀 this service can spawn other services
🌍 [process-spawner-01] http(s) server is listening on: 8080
```

## Build a Simplism service

```bash
cd say-hello
tinygo build -scheduler=none --no-debug \
-o say-hello.wasm \
-target wasi main.go
```

## "Spawn" 3 Simplism services

```bash
curl -X POST \
http://localhost:8080/spawn \
-H 'admin-spawn-token:michael-burnham-rocks' \
-H 'Content-Type: application/json; charset=utf-8' \
--data-binary @- << EOF
{
    "wasm-file":"../say-hello/say-hello.wasm", 
    "wasm-function":"handle", 
    "http-port":"9091", 
    "discovery-endpoint":"http://localhost:8080/discovery", 
    "admin-discovery-token":"michael-burnham-rocks",
    "information": "✋ I'm listening on port 9091",
    "service-name": "say-hello_9091"
}
EOF
echo ""

curl -X POST \
http://localhost:8080/spawn \
-H 'admin-spawn-token:michael-burnham-rocks' \
-H 'Content-Type: application/json; charset=utf-8' \
--data-binary @- << EOF
{
    "wasm-file":"../say-hello/say-hello.wasm", 
    "wasm-function":"handle", 
    "http-port":"9092", 
    "discovery-endpoint":"http://localhost:8080/discovery", 
    "admin-discovery-token":"michael-burnham-rocks",
    "information": "🖖 I'm listening on port 9092",
    "service-name": "say-hello_9092"
}
EOF
echo ""

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
    "admin-discovery-token":"michael-burnham-rocks",
    "information": "👋 I'm listening on port 9093",
    "service-name": "say-hello_9093"
}
EOF
echo ""
```

You should get:
```bash
👋 this service is discoverable
🌍 http server is listening on: 9091
👋 this service is discoverable
🌍 http server is listening on: 9092
👋 this service is discoverable
🌍 http server is listening on: 9093
```

## Get the list of the Simplism services

```bash
curl http://localhost:8080/discovery \
-H 'admin-discovery-token:michael-burnham-rocks'
```

You should get:
```json
{"39253":{"pid":39253,"functionName":"handle","filePath":"../say-hello/say-hello.wasm","recordTime":"2023-12-10T06:21:29.006957886Z","startTime":"2023-12-10T06:19:28.967587761Z","stopTime":"0001-01-01T00:00:00Z","httpPort":"9091","information":"✋ I'm listening on port 9091","serviceName":"say-hello_9091"},"39260":{"pid":39260,"functionName":"handle","filePath":"../say-hello/say-hello.wasm","recordTime":"2023-12-10T06:21:29.018214886Z","startTime":"2023-12-10T06:19:28.971749844Z","stopTime":"0001-01-01T00:00:00Z","httpPort":"9092","information":"🖖 I'm listening on port 9092","serviceName":"say-hello_9092"},"39268":{"pid":39268,"functionName":"handle","filePath":"../say-hello/say-hello.wasm","recordTime":"2023-12-10T06:21:29.019056094Z","startTime":"2023-12-10T06:19:28.977488636Z","stopTime":"0001-01-01T00:00:00Z","httpPort":"9093","information":"👋 I'm listening on port 9093","serviceName":"say-hello_9093"}}
```

## Query the Simplism services

```bash
curl http://localhost:9093 &
curl http://localhost:9092 &
curl http://localhost:9091
```

You should get:
```bash
🤗 Hello 👋
🤗 Hello 👋
🤗 Hello 👋
```

## Kill a Simplism service

```bash
curl -X DELETE \
http://localhost:8080/spawn?simplismid=39801 \
-H 'admin-spawn-token:michael-burnham-rocks'
```

and try:
```bash
curl http://localhost:9093
```

You should get:
```bash
curl: (7) Failed to connect to localhost port 9093 after 0 ms: Connection refused
```
