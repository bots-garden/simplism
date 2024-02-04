#!/bin/bash
curl http://localhost:9000/service/hello \
-d 'Bob Morane'

echo

curl http://localhost:9000/service/hello-bis \
-d 'Bob Morane'
