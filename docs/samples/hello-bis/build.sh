#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o hello-bis.wasm \
-target wasi main.go
