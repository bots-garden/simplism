#!/bin/bash

#curl -X DELETE http://localhost:9000/spawn/pid/87277 \
#-H 'admin-spawn-token: michael-burnham-rocks'

curl -X DELETE http://localhost:9000/spawn/name/hello \
-H 'admin-spawn-token: michael-burnham-rocks'

#curl -X DELETE http://localhost:9000/spawn/name/hello-people \
#-H 'admin-spawn-token: michael-burnham-rocks'