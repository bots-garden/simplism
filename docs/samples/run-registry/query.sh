#!/bin/bash
# POST
curl http://localhost:9090/registry \
-H 'content-type: application/json; charset=utf-8' \
-H 'admin-registry-token: morrison-hotel' \
-d '{"firstName":"Bob","lastName":"Morane"}'

# GET
curl http://localhost:9090/registry \
-H 'content-type: application/json; charset=utf-8' \
-H 'admin-registry-token: morrison-hotel' \
