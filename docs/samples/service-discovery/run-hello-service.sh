#!/bin/bash
simplism listen \
../hello/hello.wasm handle \
--http-port 8082 \
--log-level info \
--service-name hello \
--admin-discovery-token people-are-strange \
--discovery-endpoint http://localhost:9000/discovery

