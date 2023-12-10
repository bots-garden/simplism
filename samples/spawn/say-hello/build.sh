#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o say-hello.wasm \
-target wasi main.go
