#!/bin/bash
tinygo build -o index.wasm  -target wasi ./index.go

ls -lh *.wasm