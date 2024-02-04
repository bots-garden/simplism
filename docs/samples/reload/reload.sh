#!/bin/bash
# ✋ this won't update the service discovery
# this api is a facility

read -r -d '' params << EOM
{
  "wasm-url":"http://localhost:9090/registry/pull/hello-bis.wasm", 
  "wasm-file":"./hello-bis.wasm", 
  "wasm-function":"handle", 
  "wasm-url-auth-header":"private-registry-token: people-are-strange",
}
EOM


curl -v -X POST \
http://localhost:8080/reload \
-H 'content-type: application/json; charset=utf-8' \
-H 'admin-reload-token:1234567890' \
-d "$params"
