# nexus-calibration
A tool to run load and calibrate Nexus runtime.


## Getting started
Build the docker image of the tool and push to a docker registry. Once uploaded the image can be used to run workload on any kubernetes cluster

### Build
Build docker image and publish it
```
make docker_build
docker tag nexus-calib:latest <image_name:version>  
docker push <image_name:version> 
```

### Install
- Install tsdb to store test data (postgres db is used to store timescale data)
```
kubectl create ns timescale
helm install tsdb timescale/timescaledb-single -n timescale

PGPASSWORD_POSTGRES=$(kubectl get secret --namespace timescale "tsdb-credentials" -o jsonpath="{.data.PATRONI_SUPERUSER_PASSWORD}" | base64 --decode)

# login to psql and setup the database and tables with appropriate schemas
kubectl run -i --tty --rm psql --image=postgres \
      --env "PGPASSWORD=$PGPASSWORD_POSTGRES" \
      --command -- psql -U postgres \
      -h tsdb.timescale.svc.cluster.local postgres

postgres-# CREATE DATABASE testdb;

# connect to database 
postgres=# \c testdb;

# create table and setup the index
postgres-#  CREATE TABLE trace_data (
Timestamp  TIMESTAMPTZ NOT NULL,
Duration INT NOT NULL,
Name TEXT NOT NULL,
Error INT NOT NULL,
Trace_id TEXT NOT NULL,
Message TEXT NOT NULl
);


postgres-#  SELECT create_hypertable('trace_data','timestamp');
CREATE INDEX ix_symbol_time ON trace_data (name, timestamp DESC);
```
  The above commands creates the database and exposes the service at - postgres://tsdb.timescale:5432
- Install zipkin, backed by mysql database
```
kubectl apply -f deployment/zipkin-mysql.yaml
```
  Installs zipkin along with mysql, and exposes zipkin endpoint at - http://zipkin:9411

- Install grafana dashboard
```
kubectl apply -f deployment/grafana.yaml
```
- Install the nexus calibration tool
```
kubectl apply -f deployment/config-map.yaml
kubectl apply -f deployment/deployment.yaml
```
  Installs the tool and exposes the service at http://nexus-calibrate:8000


### Run the workload
Once the installation is complete, the tool can be accessed via port forwarding the service to locahost

#### Steps to add tests 
- To add rest test 
  Use postman or a similar http client
  Do a POST operation - localhost:8000/rest/tests
  with the below body. 
```
{
    "spec": [
        {
            "name": "put_manager",
            "method": "PUT",
            "path": "/root/default/leader/default/mgr/m1",
            "data": "{ \"employeeID\": 0,\"name\": \"string\"}"
        },
        {
            "name": "put_employee",
            "method": "PUT",
            "path": "/root/default/employee/{{random}}",
            "data": "{}"
        },
        {
            "name": "put_manager2",
            "method": "PUT",
            "path": "/root/default/leader/default/mgr/{{random}}",
            "data": "{ \"employeeID\": 0,\"name\": \"string\"}"
        },
        {
            "name": "get_one_manager",
            "method": "GET",
            "path": "/root/default/leader/default/mgr/m1"
        },
        {
            "name": "get_managers_rest",
            "method": "GET",
            "path": "/mgrs"
        }
    ]
}
```
- To add graphql query tests
  POST call to - localhost:8000/gql/tests
  This will add the below tests to the tool
```
{
    "spec": [
        {
            "name": "get_mgrs",
            "method": "{root {    Id, ParentLabels, CEO { EngManagers { Id }}}}"
    }, {
            "name": "get_leaders",
            "method": "{root { Id, ParentLabels, CEO {Id}}}"
    }
]
}
```
The tests that have been added can now be accessed/run by using the "name" keyword that was supplied while adding them

#### To run a particular test
Do a POST call to - localhost:8000/tests/t1
```
{
    "server": {
        "url": "http://nexus-api-gw.nexus5",
        "zipkin": "http://zipkin:9411",
        "tsdb": "postgres://postgres:dXmrYXVfwgD2JZvl@tsdb.timescale:5432/testdb"
    },
    "tests": [
        {
            "name": "write_n_objects",
            "concurrency": 10,
            "ops_count": 5,
            "sample_rate": 1.0,
            "graphql": [
                "get_mgrs",
                "get_leaders"
            ]
        }
    ]
}
```

The workload can be re-run by creating a new job

## Architecture of the tool

Below image shows the high level architecture
![Architecture image](images/nexus_calib_tool_image.jpeg?raw=true "Tool Architecture")
