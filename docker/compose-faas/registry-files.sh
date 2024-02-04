#!/bin/bash
curl http://localhost:9090/registry/discover \
-H 'content-type:text/plain; charset=UTF-8' \
-H 'private-registry-token: people-are-strange'