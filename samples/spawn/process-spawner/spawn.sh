#!/bin/bash
#curl http://localhost:8080/spawn \
#-H 'admin-spawn-token:michael-burnham-rocks'
#echo ""

curl -X POST \
http://localhost:8080/spawn \
-H 'admin-spawn-token:michael-burnham-rocks' \
-H 'Content-Type: application/json; charset=utf-8' \
-d '{"wasm-file":"../say-hello/say-hello.wasm", "wasm-function":"handle", "http-port":"9091", "discovery-endpoint":"http://localhost:8080/discovery", "admin-discovery-token":"michael-burnham-rocks"}'
echo ""

curl -X POST \
http://localhost:8080/spawn \
-H 'admin-spawn-token:michael-burnham-rocks' \
-H 'Content-Type: application/json; charset=utf-8' \
-d '{"wasm-file":"../say-hello/say-hello.wasm", "wasm-function":"handle", "http-port":"9092", "discovery-endpoint":"http://localhost:8080/discovery", "admin-discovery-token":"michael-burnham-rocks"}'
echo ""

curl -X POST \
http://localhost:8080/spawn \
-H 'admin-spawn-token:michael-burnham-rocks' \
-H 'Content-Type: application/json; charset=utf-8' \
--data-binary @- << EOF
{
    "wasm-file":"../say-hello/say-hello.wasm", 
    "wasm-function":"handle", 
    "http-port":"9093", 
    "discovery-endpoint":"http://localhost:8080/discovery", 
    "admin-discovery-token":"michael-burnham-rocks"
}
EOF
echo ""


