#!/bin/bash
clear
bat $0 --line-range 5:
echo ""
minism call simple.wasm say_hello \
--input "Bob Morane" \
--log-level info \
--config '{"firstName":"Jane","lastName":"Doe"}'

