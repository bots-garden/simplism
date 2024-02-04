#!/bin/bash
cd /projects/hello
go get github.com/extism/go-pdk
tinygo build -scheduler=none --no-debug \
-o hello.wasm \
-target wasi main.go
