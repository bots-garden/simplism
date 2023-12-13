#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o service-one.wasm \
-target wasi main.go
