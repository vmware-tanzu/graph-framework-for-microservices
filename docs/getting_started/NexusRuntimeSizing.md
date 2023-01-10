## Nexus Sizing

  This page helps to install nexus runtime with custom resource values (cpu/memory)
  
  User has to configure both cpu and memory for custom resource values. Please refer the examples as below.


### To configure all the runtime components with custom resource values

:bulb: The values shown below are examples only, please change them as you need for your runtime


```
nexus runtime install --namespace <Namespace> \
--cpuResources api_gateway=490m  \
--cpuResources kubeapiserver=480m \
--cpuResources kubecontrollermanager=490m \
--cpuResources etcd=480m \
--cpuResources validation=490m \
--cpuResources graphql=490m \
--cpuResources nexus_controller=480m \
--memoryResources api_gateway=256Mi \
--memoryResources kubeapiserver=480Mi \
--memoryResources kubecontrollermanager=256Mi \
--memoryResources etcd=480Mi \
--memoryResources validation=480Mi \
--memoryResources graphql=256Mi \
--memoryResources nexus_controller=480Mi
```

### To configure any one or more of the runtime component with custom resource values

```
nexus runtime install --namespace <Namespace> \
--cpuResources api_gateway=490m  \
--memoryResources api_gateway=256Mi 
```
