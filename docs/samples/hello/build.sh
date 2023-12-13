#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o hello.wasm \
-target wasi main.go
