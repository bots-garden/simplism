#!/bin/bash

# ----------------------------
#  Spin framework
# ----------------------------
cd hello-spin
spin build
spin up --listen 0.0.0.0:8082 &
sleep 5
hey -n 3000 -c 1000 -m GET http://0.0.0.0:8082 > ../spin-report.txt
#hey -n 3000 -c 1000 -m POST \
#-H "content-type: application/json; charset=utf-8" \
#-d '{"firstName":"Bob","lastName":"Morane"}' \
#http://0.0.0.0:8082
pkill -f spin
cd ..



