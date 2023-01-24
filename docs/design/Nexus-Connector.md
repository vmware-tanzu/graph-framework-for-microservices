# Nexus Connector

Nexus Connector "connects" two Nexus Runtimes.

![NexusConnector1](.content/images/NexusConnector1.png)

"Connect" essentially means: Synchronization of state between the two connected Nexus Runtimes.

It includes:

* on-demand replication of state between them

* extend nexus runtime across cluster / network boundaries
* applications are given a uniform interface to platform, i.e applications always work with Nexus Runtime irrespective of where they run
* auto replication / synchornization / propagation of state, without application or external trigger
* guarentees applications on either side of the connector can run in headless mode. Loss of connectivity to remote endpoint is not catastrophic. Applications continue to run on the last consistent and correct state.
* auto-reconciliation when connection is establish.

What "connect" does NOT mean:

* Nexus connector is responsible for creating/managing network connection to either runtimes. Rather, nexus connector will run on the connection provided to it.

* there is synchornization of ALL state between two runtimes. The state to be synchornized should be configured at runtime, by the product / application business logic
* nexus connector will replicate/propagate datamodel type. Nexus connector will only propagate objects at runtime, on-demand. Propagation of Datamodel type is left to be handled by product installation operator/mechanism.
* any preprocessing or postprocessing of state being replicated.
* replication of state not known to Nexus Runtime.

## State

State in a Nexus Runtime, that can be replicated are:

* Nexus Datamodel Nodes - metadata, spec and status
* sub-graph of Nexus Datamodel Nodes along with relationships - all nodes in the sub-graph, along with their metadata, spec and status


## Use cases

### Syncing of cluster specific state between SaaS and Application cluster

![NexusConnector2](.content/images/NexusConnector2.png)


## Workflow

### Configure RBAC in Nexus Runtimes, as needed, for K8s connector

Access to nexus runtime is controlled by RBAC. Its the products reponsiblity to create a user account, role and rolebinding that will be used by nexus connector to successfully replicate state.

It is recommended that a unique user account is created for each instance of nexus connector. Sharing of user account across more than one nexus connector, if there are multiple instances of it running, is discouraged.

### Bootstrap/install Nexus Connector

Configure Nexus Connector with credentials to connect to both Nexus Runtimes.

Nexus Connector will need the following to successfully connect to each of the Nexus Runtime endpoints:

1. Endpoint / domain
2. Certficates
3. User account / token
4. Enterprise proxy configuration, if any.

### Notify Nexus Runtime of intent to replicate specific state from local to remote Runtime.

## Nexus Connector Implementation

K8s connector provides a 2 way


## To be done / discussed
1. Nexus Datamodel implementation