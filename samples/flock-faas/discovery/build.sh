#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o discovery.wasm \
-target wasi main.go
