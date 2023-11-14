#!/bin/bash
clear
bat $0 --line-range 5:
echo ""
go run main.go \
./samples/go-plugin/simple.wasm \
say_hello \
8080

