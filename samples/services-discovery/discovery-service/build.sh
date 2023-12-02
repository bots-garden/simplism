#!/bin/bash
tinygo build -scheduler=none --no-debug \
-o discovery-service.wasm \
-target wasi main.go
