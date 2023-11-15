#!/bin/bash
clear
bat $0 --line-range 5:
echo ""
go run ../../../main.go listen \
./target/wasm32-wasi/release/simple_plugin.wasm hello --http-port 8080 --log-level info
