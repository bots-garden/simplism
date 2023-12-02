#!/bin/bash

for i in {1..100}
do
    # call the discovery service
    curl http://localhost:9000

    # move the raiders
    curl http://localhost:8001/move
    curl http://localhost:8002/move
    curl http://localhost:8003/move

    # ask the basestar for the raiders list
    curl http://localhost:8010/raiders
done
