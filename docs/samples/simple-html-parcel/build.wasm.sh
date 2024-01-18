#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o index.wasm \
-target wasi main.go

ls -lh *.wasm

