#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o raider.wasm \
-target wasi main.go
