#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o store.wasm \
-target wasi main.go
