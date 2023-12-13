#!/bin/bash
simplism listen \
service-discovery.wasm handle \
--http-port 9000 \
--log-level info \
--service-discovery true \
--admin-discovery-token people-are-strange

