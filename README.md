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
kubectl apply -f deployment/deployment.yaml
```
  Installs the tool and exposes the service at http://nexus-calibrate:8000


### Run the workload
Use the configmap provided in the deployment folder to set the configuration to requirements
The available functions can be called by adding the keys to the configMap
The workload data is captured locally by calculating the time taken to run each call. Low, High and average times are provided ( rudimentary calculation ).
More information can be accessed by configuring with zipkin end point. ( Currently, it is mandatory to configure.) 
- Apply grafana dashboad yaml
- Apply the zipkin deployment file to add zipkin backend services. 
- create timescale ns and install tsdb.  ``` helm install tsdb timescale/timescaledb-single -n timescale ```
- Now apply the config map yaml and then the deployment yaml to run the workload

The workload can be re-run by creating a new job and reconfiguring the config map

## Architecture of the tool

Below image shows the high level architecture
![Architecture image](images/nexus_calib_tool_image.jpeg?raw=true "Tool Architecture")
