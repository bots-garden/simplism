#!/bin/bash
curl http://localhost:9000/discovery \
-H 'content-type:application/json; charset=UTF-8' \
-H 'admin-discovery-token:people-are-strange'

echo ""
echo ""

curl http://localhost:9000/discovery \
-H 'content-type:text/plain; charset=UTF-8' \
-H 'admin-discovery-token:people-are-strange'