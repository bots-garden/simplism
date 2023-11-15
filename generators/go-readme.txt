# <name>

> Build the wasm plugin:
```bash
tinygo build -scheduler=none --no-debug \
  -o <name>.wasm \
  -target wasi main.go
```

> Serve the wasm plugin with Simplism:
```bash
./simplism listen \
<name>.wasm handle --http-port 8080 --log-level info
```

> Query the wasm plugin:
```bash
curl http://localhost:8080/hello/world \
-H 'content-type: application/json; charset=utf-8' \
-d '{"firstName":"Bob","lastName":"Morane"}'
```
