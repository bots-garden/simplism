#!/bin/bash
set -o allexport; source .env; set +o allexport
k9s --all-namespaces


