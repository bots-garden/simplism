This file is embedded in the simplism binary. If you need to rebuild it, go to the `./scratch` folder.

## Explanation

When you start Simplism with the:

- Registry mode (`--registry-mode`)
- Store mode (`--store-mode`)
- Discovery mode (`--service-discovery`)

You can use the `--wasm-file` flag  with this value: `?` instead of specify the path to the WASM file. Then, Simplism will extract and use the `scratch.wasm` file.

> examples:

```bash
simplism listen ? ? \
  --http-port 9090 \
  --log-level info \
  --registry-mode true \
  --registry-path ./wasm-files
```

```yaml
tiny-registry-config:
  wasm-file: "?"
  wasm-function: "?"
  http-port: 9090
  log-level: info
  registry-mode: true
  registry-path: ./wasm-files
```
