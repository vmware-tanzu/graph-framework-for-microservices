for i in {71..100} 
do
var1='{
    "server": {
        "url": "http://nexus-api-gw.nexus3",
        "zipkin": "http://zipkin:9411",
        "tsdb": "postgres://postgres:dXmrYXVfwgD2JZvl@tsdb.timescale:5432/testdb"
    },
    "tests": [
        {
            "name": "write_n_objects",
            "concurrency": 1,
            "ops_count": 1000,
            "sample_rate": 1,
            "rest": [
                "put_operations'$i'"
            ]
        }
    ]
}'

echo curl --location --request POST 'localhost:8000/tests/t'$i \
--header \'Content-Type: application/json\' \
--data-raw \'$var1\'
echo sleep 1
done
