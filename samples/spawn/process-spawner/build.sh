#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o process-spawner.wasm \
-target wasi main.go
