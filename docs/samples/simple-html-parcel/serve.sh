#!/bin/bash
echo "open 🌍 http://localhost:8080/"
simplism listen \
index.wasm handle --http-port 8080 --log-level info
