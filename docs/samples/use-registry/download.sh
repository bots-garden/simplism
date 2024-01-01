#!/bin/bash
curl http://localhost:9090/registry/pull/hello.wasm -o hello-with-curl.wasm \
-H 'private-registry-token: people-are-strange'
