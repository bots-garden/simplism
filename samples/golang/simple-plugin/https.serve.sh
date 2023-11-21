#!/bin/bash
clear
bat $0 --line-range 5:
echo ""
go run ../../../main.go listen \
simple.wasm say_hello --http-port 443 \
  --cert-file simplism.bots.garden.crt \
  --key-file simplism.bots.garden.key

