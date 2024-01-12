#!/bin/bash

curl -X POST \
http://localhost:9000/spawn \
-H 'admin-spawn-token:michael-burnham-rocks' \
-H 'Content-Type: application/json; charset=utf-8' \
--data-binary @- << EOF
{
    "wasm-file":"../hello/hello.wasm", 
    "wasm-function":"handle", 
    "discovery-endpoint":"http://localhost:9000/discovery", 
    "admin-discovery-token":"people-are-strange",
    "admin-spawn-token":"michael-burnham-rocks",
    "information": "✋ I'm the hello service",
    "service-name": "hello"
}
EOF
echo ""
