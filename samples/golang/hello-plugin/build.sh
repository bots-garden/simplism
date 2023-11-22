#!/bin/bash
clear
bat $0 --line-range 5:
echo ""
tinygo build -scheduler=none --no-debug \
  -o hello.wasm \
  -target wasi main.go

ls -lh *.wasm
