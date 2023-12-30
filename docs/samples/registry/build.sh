#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o tiny-registry.wasm \
-target wasi main.go
