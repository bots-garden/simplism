#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o hello-world.wasm \
-target wasi main.go
