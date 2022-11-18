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
