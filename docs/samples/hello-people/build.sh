#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o hello-people.wasm \
-target wasi main.go
