#!/bin/bash
cargo clean
cargo build --release --target wasm32-wasi
