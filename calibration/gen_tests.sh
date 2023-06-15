#!/bin/sh
for i in {1..100}
do
        echo '          {
            "name": "put_operations'$i'",
            "method": "PUT",
            "path": "/root/default/leader/default/mgr/m'$i'/operations/{{random}}",
            "data": "{ \"employeeID\": 0,\"name\": \"string\"}"  
        },'
done
