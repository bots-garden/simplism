# hello-world

> Add Extism dependency
```bash
go get github.com/extism/go-pdk
```

> Build the wasm plug-in:
```bash
tinygo build -scheduler=none --no-debug \
-o hello-world.wasm \
-target wasi main.go
```

> Serve the wasm plug-in with Simplism:
```bash
simplism listen \
hello-world.wasm handle --http-port 8080 --log-level info
```

> Query the wasm plug-in:
```bash
curl http://localhost:8080/hello/world \
-H 'content-type: application/json; charset=utf-8' \
-d '{"firstName":"Bob","lastName":"Morane"}'
```
