#!/bin/bash
clear
bat $0 --line-range 5:
echo ""
extism call simple.wasm say_hello \
  --input "Bob Morane" \
  --log-level info \
  --set-config '{"firstName":"Jane","lastName":"Doe"}' \
  --wasi true
echo ""
