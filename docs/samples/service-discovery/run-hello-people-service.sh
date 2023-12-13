#!/bin/bash
simplism listen \
../hello-people/hello-people.wasm handle \
--http-port 8081 \
--log-level info \
--service-name hello-people \
--admin-discovery-token people-are-strange \
--discovery-endpoint http://localhost:9000/discovery





