#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o hey-two.wasm \
-target wasi main.go
