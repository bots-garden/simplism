#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o ../server/embedded/scratch.wasm \
-target wasi main.go
