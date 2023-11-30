#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o hey-one.wasm \
-target wasi main.go
