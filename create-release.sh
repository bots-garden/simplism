#!/bin/bash
set -o allexport; source .release.env; set +o allexport
# You need to create a .release.env file with these variables
# - GITHUB_TOKEN
goreleaser release --snapshot
