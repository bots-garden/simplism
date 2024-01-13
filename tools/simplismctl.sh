#!/bin/bash
set -o allexport; source .simplismconfig.env; set +o allexport

: <<'COMMENT'
## Get list of the simplism processes

- ./simplismctl.sh discovery list table
- ./simplismctl.sh discovery list json

## Call a function (service)

`./simplismctl.sh service call <service_name> <data> [--json | --text]`

- ./simplismctl.sh service call hello bob --text
- ./simplismctl.sh service call hello bob --json

- ./simplismctl.sh service call hello-people '{"firstName":"Bob","lastName":"Morane"}' --json

COMMENT


if [[ "$1" == "version" ]]; then
  simplism version 
  exit 0
fi

# Service Discovery
if [[ "$1" == "discovery" ]]; then
  # services list
  if [[ "$2" == "list" ]]; then
    if [[ "$3" == "table" ]]; then
        curl ${DISCOVERY_SERVICE_URL} \
        -H 'content-type:text/plain; charset=UTF-8' \
        -H "admin-discovery-token:${ADMIN_DISCOVERY_TOKEN}"
        exit 0
    fi

    if [[ "$3" == "json" ]]; then
        curl ${DISCOVERY_SERVICE_URL} \
        -H 'content-type:application/json; charset=UTF-8' \
        -H "admin-discovery-token:${ADMIN_DISCOVERY_TOKEN}"
        exit 0
    fi
    
    echo "😡 command unknown"
    exit 1
  fi

  # ...
  #if [[ "$2" == "vm" ]]; then
  #  . $(dirname "$0")/start-vm.sh
  #  exit 0
  #fi

fi

# Call function
if [[ "$1" == "service" ]]; then

  if [[ "$2" == "call" ]]; then
    service_name="$3"
    data="$4"
    content_type="$5"
    
    header_content_type="text/plain; charset=UTF-8"

    if [[ "$content_type" == "--json" ]]; then
        header_content_type="application/json; charset=UTF-8"
    fi
    
    if [[ "$content_type" == "--text" ]]; then
        header_content_type="text/plain; charset=UTF-8"
    fi

    curl ${CALL_SERVICE_URL}/${service_name} \
    -H "content-type:${header_content_type}" \
    -d "${data}"
    
  fi
  #echo "😡 command unknown"
  #exit 1
fi