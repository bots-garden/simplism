#!/bin/bash
curl http://localhost:9090/registry/push \
-H 'admin-registry-token: morrison-hotel' \
-F 'file=@hello.wasm;filename=hello.0.0.1.wasm'

curl http://localhost:9090/registry/push \
-H 'admin-registry-token: morrison-hotel' \
-F 'file=@hello-people.wasm;filename=hello-people.0.0.0.wasm'
