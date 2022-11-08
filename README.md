# nexus-calibration
A tool to run load and calibrate Nexus runtime.


## Getting started
Build the docker image of the tool and push to a docker registry. Once uploaded the image can be used to run workload on any kubernetes cluster

### Build
Build docker image and publish it
```
make docker_build
docker push <image_name:version> 
```
### Run the workload
Use the configmap provided in the deployment folder to set the configuration to requirements
The available functions can be called by adding the keys to the configMap
The workload data is captured locally by calculating the time taken to run each call. Low, High and average times are provided ( rudimentary calculation ).
More information can be accessed by configuring with zipkin end point. ( Currently, it is mandatory to configure.) 
```
apiVersion: v1
kind: ConfigMap
metadata:
  name: nexus-calib-config
data:
  conf.yaml: |
    server:
      url: http://nexus-api-gw.nexus
      zipkin: http://zipkin:9411
    concurrency: 5
    timeout: 60
    rest:
    - put_employee
    - get_hr
    graphql:
    - get_managers
    - get_employee_role
```
- Apply the zipkin deployment file to add zipkin backend services. 
- Now apply the config map yaml and then the deployment yaml to run the workload

The workload can be re-run by creating a new job and reconfiguring the config map

## Architecture of the tool

