#!/bin/bash
rm discovery-service/discovery-service.wasm.db
# This simplism service will record all the other simplism services
# The other simplism services will post their information to this service
# on http://localhost:9000/discovery
simplism listen discovery-service/discovery-service.wasm handle \
--http-port 9000 \
--log-level info \
--service-discovery true \
--admin-discovery-token people-are-strange &

# A simplism service is discoverable since it uses
# --discovery-endpoint flag
simplism listen service-one/service-one.wasm handle \
--http-port 8001 \
--log-level info \
--discovery-endpoint http://localhost:9000/discovery \
--admin-discovery-token people-are-strange &

simplism listen service-two/service-two.wasm handle \
--http-port 8002 \
--log-level info \
--discovery-endpoint http://localhost:9000/discovery \
--admin-discovery-token people-are-strange &

simplism listen service-three/service-three.wasm handle \
--http-port 8003 \
--log-level info \
--discovery-endpoint http://localhost:9000/discovery \
--admin-discovery-token people-are-strange &

# curl http://localhost:9000
# curl http://localhost:9000/discovery
# pkill -f simplism
