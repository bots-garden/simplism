#!/bin/bash
#rm ../service-discovery/*.db
simplism listen \
? ? \
--http-port 9000 \
--log-level info \
--service-discovery true \
--admin-discovery-token people-are-strange \
--information "👋 I'm the spawner service" \
--spawn-mode true \
--http-port-auto true \
--admin-spawn-token michael-burnham-rocks
#--recovery-mode false
# ../service-discovery/service-discovery.wasm handle \
