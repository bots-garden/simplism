#!/bin/bash
#curl http://localhost:8080/spawn \
#-H 'admin-spawn-token:michael-burnham-rocks'
#echo ""

curl -X POST \
http://localhost:8080/spawn \
-H 'admin-spawn-token:michael-burnham-rocks' \
-H 'Content-Type: application/json; charset=utf-8' \
--data-binary @- << EOF
{
    "wasm-file":"../say-hello/say-hello.wasm", 
    "wasm-function":"handle", 
    "http-port":"9091", 
    "discovery-endpoint":"http://localhost:8080/discovery", 
    "admin-discovery-token":"michael-burnham-rocks",
    "information": "✋ I'm listening on port 9091",
    "service-name": "say-hello_9091"
}
EOF
echo ""

curl -X POST \
http://localhost:8080/spawn \
-H 'admin-spawn-token:michael-burnham-rocks' \
-H 'Content-Type: application/json; charset=utf-8' \
--data-binary @- << EOF
{
    "wasm-file":"../say-hello/say-hello.wasm", 
    "wasm-function":"handle", 
    "http-port":"9092", 
    "discovery-endpoint":"http://localhost:8080/discovery", 
    "admin-discovery-token":"michael-burnham-rocks",
    "information": "🖖 I'm listening on port 9092",
    "service-name": "say-hello_9092"
}
EOF
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
    "admin-discovery-token":"michael-burnham-rocks",
    "information": "👋 I'm listening on port 9093",
    "service-name": "say-hello_9093"
}
EOF
echo ""

curl -X POST \
http://localhost:8080/spawn \
-H 'admin-spawn-token:michael-burnham-rocks' \
-H 'Content-Type: application/json; charset=utf-8' \
--data-binary @- << EOF
{
    "wasm-file":"../hello-you/hello-you.wasm", 
    "wasm-function":"handle", 
    "http-port":"9094", 
    "discovery-endpoint":"http://localhost:8080/discovery", 
    "admin-discovery-token":"michael-burnham-rocks",
    "information": "👋 I'm listening on port 9093",
    "service-name": "hello-you"
}
EOF
echo ""

