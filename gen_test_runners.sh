for i in {1..100} 
do
var1='{
    "server": {
        "url": "http://nexus-api-gw.scaletest9",
        "zipkin": "http://zipkin:9411",
        "tsdb": "postgres://postgres:Ordjdgz9zd9Hii07@tsdb.timescale:5432/testdb"
    },
    "tests": [
        {
            "name": "write_n_objects",
            "concurrency": 1,
            "ops_count": 1,
            "sample_rate": 1.0,
            "rest": [
                "put_operations10"
            ]
        }
    ]
}'
 
echo curl --location --request POST 'localhost:8000/tests/t'$i \
--header \'Content-Type: application/json\' \
--data-raw \'$var1\'
echo sleep 10
done
