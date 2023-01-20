# Nexus Scaling

This document captures the details about nexus scaletests and the max resource utilization and turn-around time in the event of any particualr service failures.

## Testcase1

### Org Chart Datamodel

```mermaid
graph TD;
  Root-->Leader;
  Root-->Employee;
  Root-->Executive;
 
  Leader-->Manager;
  Leader-->HumanResources;
  Leader -- link --> Executive;
 
  Manager --> Dev;
  Manager --> Operations;
  Manager -- link --> Employee;
 
  HumanResources-- link --> Employee;

  Dev -- link --> Employee;

  Operations -- link --> Employee;
 ```

```
Scale database to have 100K objects

100 Managers (parent)

1000 Ops (children)
```

### Resource Configuration Per Service

| Service | CPU | Memory | Replicas | Resource Usage at peak |
|---------|-----|--------|----------|------------------------|
|Nexus API Gateway|490m|512Mi  |1| ![](images/api-gw.png?raw=true)|
|Nexus Kube API Server|480m|No Limit|1|![](images/apiserver.png?raw=true)|
|Nexus Kube Ctrl Mgr|490m|512Mi|1| ![](images/kcm.png?raw=true)|
|Nexus ETCD|480m|No Limit|1|
|Nexus GraphQL|490m|2Gi|1| ![](images/graphql.png?raw=true)|
|Nexus Validation|490m|480Mi|1|
|Nexus Controller|480m|480Mi|1|

### Key Stats With 100K objects in the system

1. GraphQL responds within ~ 8 seconds to query 100K objects
2. On GraphQL Server restart, it takes ~ 4 Minutes to respond first successful query  
