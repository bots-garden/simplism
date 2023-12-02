#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o service-two.wasm \
-target wasi main.go
