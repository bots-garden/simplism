#!/bin/bash
simplism listen \
store.wasm handle \
--http-port 8080 \
--log-level info \
--store-mode true \
--admin-store-token morrison-hotel \
--information "👋 I'm the store service" \
--service-name store

