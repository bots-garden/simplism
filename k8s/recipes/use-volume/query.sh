#!/bin/bash
set -o allexport; source ../.env; set +o allexport
set -o allexport; source .env; set +o allexport

curl http://${APPLICATION_NAME}.${DNS} -d '👋 Hello World 🌍 on Civo'

