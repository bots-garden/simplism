#!/bin/bash
curl http://localhost:8001 \
-H 'content-type: application/json; charset=utf-8' \
-d '{"messge":"hello"}'

curl http://localhost:8002 \
-H 'content-type: application/json; charset=utf-8' \
-d '{"messge":"hello"}'
