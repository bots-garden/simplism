# <name>

> Build the wasm plugin:
```bash
cargo clean
cargo build --release --target wasm32-wasi
```

> Serve the wasm plugin with Simplism:
```bash
./simplism listen \
./target/wasm32-wasi/release/<name>.wasm handle --http-port 8080 --log-level info
```

> Query the wasm plugin:
```bash
curl http://localhost:8080/hello/world \
-H 'content-type: application/json; charset=utf-8' \
-d '{"firstName":"Bob","lastName":"Morane"}'
```
