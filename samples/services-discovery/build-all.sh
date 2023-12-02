#!/bin/bash
cd discovery-service
./build.sh
cd ..

cd service-one
./build.sh
cd ..

cd service-two
./build.sh
cd ..

cd service-three
./build.sh
cd ..

