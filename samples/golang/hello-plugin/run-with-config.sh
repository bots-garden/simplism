#!/bin/bash
clear
bat $0 --line-range 5:
echo ""
go run ../../../main.go config ./config.yml hello-plugin &

#go run ../../../main.go config ./config.yml hello-plugin-1 &

#go run ../../../main.go config ./config.yml hello-plugin-2 &


