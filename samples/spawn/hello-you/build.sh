#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o hello-you.wasm \
-target wasi main.go
