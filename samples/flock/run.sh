#!/bin/bash
clear
bat $0 --line-range 5:
echo ""
rm discovery/discovery.wasm.db
simplism flock ./config.yml

