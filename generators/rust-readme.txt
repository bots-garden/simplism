# <name>

> Build the wasm plug-in:
```bash
cargo clean
cargo build --release --target wasm32-wasi
```

> Serve the wasm plug-in with Simplism:
```bash
simplism listen \
./target/wasm32-wasi/release/<name>.wasm handle --http-port 8080 --log-level info
```

> Query the wasm plug-in:
```bash
curl http://localhost:8080/hello/world \
-H 'content-type: application/json; charset=utf-8' \
-d '{"firstName":"Bob","lastName":"Morane"}'
```
