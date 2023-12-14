#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o use-store.wasm \
-target wasi main.go
