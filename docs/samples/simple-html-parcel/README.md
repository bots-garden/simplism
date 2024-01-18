# Simple HTML with Parcel

## Setup

```bash
npm install --save-dev parcel
npm install lit-element
```

## Build

```bash
# build html asset
npm run build
# build wasm file
tinygo build -scheduler=none --no-debug \
-o index.wasm \
-target wasi main.go
```

## Run

```bash
simplism listen \
index.wasm handle --http-port 8080 --log-level info
```
