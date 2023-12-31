#!/bin/bash

hey -n 3000 -c 1000 -m POST \
-H "content-type: application/json; charset=utf-8" \
-d '{"firstName":"Bob","lastName":"Morane"}' \
http://0.0.0.0:8080 

#> ../simplism-hello-world-report.txt
