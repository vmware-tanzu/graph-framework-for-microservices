#!/usr/bin/env bash

set -x

expected=$1

RETRY=4
ATTEMPT=0
echo 'testing header-based routing with a valid header...'
while [ $ATTEMPT -ne $RETRY ]; do
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" localhost:10000/leaders?orgchart.Root=default -H "x-tenant:t-1")
  if [ "$STATUS" == "$expected" ] ; then
    echo "SUCCESS"
    break
  else
    ATTEMPT=$((ATTEMPT + 1))
    if [ $ATTEMPT -eq $RETRY ]; then
      echo "FAILED: $STATUS. expected $expected"
      exit 1
    else
      sleep 5
    fi
  fi
done

echo 'testing header-based routing with an invalid header...'
while [ $ATTEMPT -ne $RETRY ]; do
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" localhost:10000/leaders?orgchart.Root=default -H "x-tenant:asdf")
  if [ "$STATUS" == "400" ] ; then
    echo "SUCCESS"
    break
  else
    ATTEMPT=$((ATTEMPT + 1))
    if [ $ATTEMPT -eq $RETRY ]; then
      echo "FAILED: $STATUS. expected 400"
      exit 1
    else
      sleep 5
    fi
  fi
done
