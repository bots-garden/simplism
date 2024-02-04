#!/bin/bash
curl -X POST \
http://localhost:9000/spawn \
-H 'admin-spawn-token:michael-burnham-rocks' \
-H 'Content-Type: application/json; charset=utf-8' \
--data-binary @- << EOF
{
    "wasm-url": "http://simplism-registry:9090/registry/pull/hello.wasm",
    "wasm-url-auth-header": "private-registry-token=people-are-strange",
    "wasm-file":"/tmp/hello.wasm", 
    "wasm-function":"handle", 
    "discovery-endpoint":"http://simplism-spawner:9000/discovery", 
    "admin-discovery-token":"people-are-strange",
    "admin-spawn-token":"michael-burnham-rocks",
    "information": "✋ I'm the hello service",
    "service-name": "hello"
}
EOF

curl -X POST \
http://localhost:9000/spawn \
-H 'admin-spawn-token:michael-burnham-rocks' \
-H 'Content-Type: application/json; charset=utf-8' \
--data-binary @- << EOF
{
    "wasm-url": "http://simplism-registry:9090/registry/pull/hello-bis.wasm",
    "wasm-url-auth-header": "private-registry-token=people-are-strange",
    "wasm-file":"/tmp/hello-bis.wasm", 
    "wasm-function":"handle", 
    "discovery-endpoint":"http://simplism-spawner:9000/discovery", 
    "admin-discovery-token":"people-are-strange",
    "admin-spawn-token":"michael-burnham-rocks",
    "information": "✋ I'm the hello-bis service",
    "service-name": "hello-bis"
}
EOF
