#!/bin/bash
#curl http://localhost:9000/discovery \
#-H 'admin-discovery-token:people-are-strange'

curl -X DELETE http://localhost:9000/spawn?simplismid=8807 \
-H 'admin-spawn-token: michael-burnham-rocks'
