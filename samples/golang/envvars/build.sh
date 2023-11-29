#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o envvars.wasm \
-target wasi main.go
