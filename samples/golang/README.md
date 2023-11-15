# Golang plug-ins

## simple-plugin

> Build
```bash
cd simple-plugin
tinygo build -scheduler=none --no-debug \
  -o simple.wasm \
  -target wasi main.go
```

> Run
```bash
go run ../../../main.go listen \
simple.wasm say_hello --http-port 8080 
```

> Query
```bash
curl http://localhost:8080
```

## hello-plugin

> Build
```bash
cd hello-plugin
tinygo build -scheduler=none --no-debug \
  -o simple.wasm \
  -target wasi main.go
```

> Run
```bash
go run ../../../main.go listen \
simple.wasm say_hello --http-port 8080 --log-level info
```

> Query
```bash
curl -v -X POST \
http://localhost:8080/hello/world \
-H 'content-type: application/json; charset=utf-8' \
-d '{"firstName":"Bob","lastName":"Morane"}'
```
