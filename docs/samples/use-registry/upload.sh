#!/bin/bash
curl http://localhost:9090/registry/push \
-H 'admin-registry-token: morrison-hotel' \
-F 'file=@hello.wasm'

curl http://localhost:9090/registry/push \
-H 'admin-registry-token: morrison-hotel' \
-F 'file=@hello-people.wasm'