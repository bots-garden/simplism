#!/bin/bash
clear
bat $0 --line-range 5:
echo ""
go run ../../../main.go listen \
simple.wasm say_hello --http-port 8080 --log-level info
