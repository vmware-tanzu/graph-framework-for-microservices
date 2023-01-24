#!/usr/bin/env bash

set -x

sleep 5
ORIGIN=$(curl -X OPTIONS 'http://localhost:5001/api/v1/namespaces/' -H 'Origin: http://domain' -H "Access-Control-Request-Method: GET" -s --head | grep 'Access-Control-Allow-Origin'| uniq | wc -l)

if [ $ORIGIN -eq 1 ]; then
    echo "CORS Header added"
    exit 0
else
    exit 1
fi

