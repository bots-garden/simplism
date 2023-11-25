#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o basestar.wasm \
-target wasi main.go
