#!/bin/bash
curl -v \
http://localhost:9000/discovery \
-H 'content-type: application/json; charset=utf-8' \
-H 'admin-discovery-token:this-is-the-way'


