#!/bin/bash

#curl -X DELETE http://localhost:9000/spawn/kill/pid/87277 \
#-H 'admin-spawn-token: michael-burnham-rocks'

curl -X DELETE http://localhost:9000/spawn/kill/name/hello \
-H 'admin-spawn-token: michael-burnham-rocks'

#curl -X DELETE http://localhost:9000/spawn/kill/name/hello-people \
#-H 'admin-spawn-token: michael-burnham-rocks'