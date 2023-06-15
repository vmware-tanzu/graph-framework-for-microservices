#!/bin/sh
#
for i in {1..70}
do
curl -X 'GET' \
  'http://localhost:45192/operations?management.Mgr=m'$i \
  -H 'accept: application/json' >> 
done
