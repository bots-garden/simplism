#!/bin/bash
clear
bat $0 --line-range 5:
echo ""
hey -n 3000 -c 1000 -m GET \
-H "Content-Type: application/json" \
"http://localhost:8080/hello/world" #> go-extism-report.txt

#hey -n 300 -c 100 -m POST \
