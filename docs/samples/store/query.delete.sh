#!/bin/bash
curl -X "DELETE" http://localhost:8080/store?key=002 \
-H 'admin-store-token: morrison-hotel'
