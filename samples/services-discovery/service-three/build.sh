#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o service-three.wasm \
-target wasi main.go
