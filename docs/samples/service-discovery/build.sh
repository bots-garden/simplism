#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o service-discovery.wasm \
-target wasi main.go
