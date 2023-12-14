#!/bin/bash
curl http://localhost:8080/store \
-H 'content-type: application/json; charset=utf-8' \
-H 'admin-store-token: morrison-hotel' \
-d '{"key":"hello","value":"hello world"}'

curl http://localhost:8080/store \
-H 'content-type: application/json; charset=utf-8' \
-H 'admin-store-token: morrison-hotel' \
-d '{"key":"hey","value":"hey people"}'

curl http://localhost:8080/store \
-H 'content-type: application/json; charset=utf-8' \
-H 'admin-store-token: morrison-hotel' \
-d '{"key":"001","value":"first"}'

curl http://localhost:8080/store \
-H 'content-type: application/json; charset=utf-8' \
-H 'admin-store-token: morrison-hotel' \
-d '{"key":"002","value":"second"}'

curl http://localhost:8080/store \
-H 'content-type: application/json; charset=utf-8' \
-H 'admin-store-token: morrison-hotel' \
-d '{"key":"003","value":"third"}'
