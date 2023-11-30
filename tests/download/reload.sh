#!/bin/bash
clear
bat $0 --line-range 5:
echo ""
curl -v -X POST \
http://localhost:8080/reload \
-H 'content-type: application/json; charset=utf-8' \
-d '{"wasm-url":"http://0.0.0.0:3333/hey-two/hey-two.wasm", "wasm-file": "./hey-two.wasm", "wasm-function": "handle"}'


