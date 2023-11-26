#!/bin/bash

for i in {1..100}
do
    curl http://localhost:8001/move
    curl http://localhost:8002/move
    curl http://localhost:8003/move

    curl http://localhost:9000/raiders
done
