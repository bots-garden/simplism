#!/bin/bash

# ----------------------------
#  Spin framework
# ----------------------------
cd hello-spin
spin build
spin up --listen 0.0.0.0:8082 &
sleep 5
hey -n 3000 -c 1000 -m GET http://0.0.0.0:8082 > ../spin-report.txt
pkill -f spin
cd ..

# ----------------------------
#  Wasm Workers Server
# ----------------------------
cd hello-wws
tinygo build -o worker.wasm -target wasi main.go
wws --host 0.0.0.0 --port 8081 . &
sleep 10
hey -n 3000 -c 1000 -m GET http://0.0.0.0:8081/worker > ../wws-report.txt
pkill -f wws
cd ..

# ----------------------------
#  Simplism Golang Plugin
# ----------------------------
cd simple-go-plugin
tinygo build -scheduler=none --no-debug -o simple.wasm -target wasi main.go
go run ../../main.go listen simple.wasm say_hello --http-port 8080 &
sleep 5
hey -n 3000 -c 1000 -m GET http://0.0.0.0:8080 > ../simplism-go-report.txt
pkill -f go
cd ..

# ----------------------------
#  Simplism Rustlang Plugin
# ----------------------------
cd simple-rust-plugin
cargo clean
cargo build --release --target wasm32-wasi
go run ../../main.go listen ./target/wasm32-wasi/release/simple_plugin.wasm hello --http-port 8080 &
sleep 5
hey -n 3000 -c 1000 -m GET http://0.0.0.0:8080 > ../simplism-rust-report.txt
pkill -f go
cd ..

