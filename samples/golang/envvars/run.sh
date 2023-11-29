#!/bin/bash
export MESSAGE="👋 Hello World 🌍"
export ABOUT="🥰 Simplism has a cute mascot 🤗"
simplism listen \
envvars.wasm handle \
--http-port 8080 \
--log-level info \
--env '["MESSAGE","ABOUT"]'

