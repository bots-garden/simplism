#!/bin/bash

for i in {1..100}
do
    curl http://localhost:8001/move
    curl http://localhost:8002/move
    curl http://localhost:8003/move

    curl http://localhost:9000/raiders
done


#hey -n 3000 -c 1000 -m GET http://localhost:8001/move &
#hey -n 3000 -c 1000 -m GET http://localhost:8002/move &
#hey -n 3000 -c 1000 -m GET http://localhost:8003/move &
