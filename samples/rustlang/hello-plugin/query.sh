#!/bin/bash
clear
bat $0 --line-range 5:
echo ""
curl -v -X POST \
http://localhost:8080/hello/world \
-H 'content-type: application/json; charset=utf-8' \
-d '{"firstName":"Bob","lastName":"Morane"}'

