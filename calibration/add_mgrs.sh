#!/bin/sh
for i in {1..100}
do
	echo 'http://localhost:45192/root/default/leader/default/mgr/m'$i 
	curl -X 'PUT' \
  'http://localhost:45192/root/default/leader/default/mgr/m'$i \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "employeeID": 1,
  "name": "string"
}'
done
