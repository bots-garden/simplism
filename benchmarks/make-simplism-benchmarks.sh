#!/bin/bash

# ----------------------------
#  Simplism Golang Plugin
# ----------------------------
cd simple-go-plugin
tinygo build -scheduler=none --no-debug -o simple.wasm -target wasi main.go
simplism listen simple.wasm say_hello --http-port 8083 &
sleep 5
hey -n 3000 -c 1000 -m GET http://0.0.0.0:8083 > ../simplism-go-report.txt
#pkill -f simplism
cd ..

# ----------------------------
#  Simplism Rustlang Plugin
# ----------------------------
cd simple-rust-plugin
cargo clean
cargo build --release --target wasm32-wasi
simplism listen ./target/wasm32-wasi/release/simple_plugin.wasm hello --http-port 8084 &
sleep 5
hey -n 3000 -c 1000 -m GET http://0.0.0.0:8084 > ../simplism-rust-report.txt
#pkill -f simplism
cd ..

