#!/bin/bash
curl http://localhost:9090/registry/pull/hello.wasm -o hello.wasm \
-H 'admin-registry-token: morrison-hotel'
