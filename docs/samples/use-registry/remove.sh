#!/bin/bash
curl -X DELETE http://localhost:9090/registry/remove/hello.wasm \
-H 'admin-registry-token: morrison-hotel'
