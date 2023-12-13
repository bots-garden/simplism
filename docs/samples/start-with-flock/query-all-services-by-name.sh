#!/bin/bash
curl http://localhost:9000/service/hello-people \
-H 'content-type: application/json; charset=utf-8' \
-d '{"firstName":"Bob","lastName":"Morane"}'

curl http://localhost:9000/service/hello \
-d 'Bob Morane'
