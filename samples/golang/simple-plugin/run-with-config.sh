#!/bin/bash
clear
bat $0 --line-range 5:
echo ""
go run ../../../main.go config ./config.yml simple-plugin

