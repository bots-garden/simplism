#!/bin/bash
rm ../service-discovery/*.db
simplism listen \
../service-discovery/service-discovery.wasm handle \
--http-port 9000 \
--log-level info \
--service-discovery true \
--admin-discovery-token people-are-strange \
--information "👋 I'm the spawner service" \
--spawn-mode true \
--admin-spawn-token michael-burnham-rocks
