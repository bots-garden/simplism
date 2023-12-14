#!/bin/bash
curl http://localhost:8080/store?prefix=00 \
-H 'content-type: application/json; charset=utf-8' \
-H 'admin-store-token: morrison-hotel'

