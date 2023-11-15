# Rustlang plug-ins

## simple-plugin

> Build
```bash
cd simple-plugin
cargo clean
cargo build --release --target wasm32-wasi
```

> Run
```bash
go run ../../../main.go listen \
./target/wasm32-wasi/release/simple_plugin.wasm hello --http-port 8080
```

> Query
```bash
curl http://localhost:8080
```

## hello-plugin

> Build
```bash
cd hello-plugin
cargo clean
cargo build --release --target wasm32-wasi
```

> Run
```bash
go run ../../../main.go listen \
./target/wasm32-wasi/release/hello_plugin.wasm hello --http-port 8080 --log-level info
```

> Query
```bash
curl -v -X POST \
http://localhost:8080/hello/world \
-H 'content-type: application/json; charset=utf-8' \
-d '{"firstName":"Bob","lastName":"Morane"}'
```

