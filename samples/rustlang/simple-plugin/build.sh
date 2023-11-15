#!/bin/bash
clear
bat $0 --line-range 5:
echo ""
cargo clean
cargo build --release --target wasm32-wasi #--offline
ls -lh ./target/wasm32-wasi/release/*.wasm
