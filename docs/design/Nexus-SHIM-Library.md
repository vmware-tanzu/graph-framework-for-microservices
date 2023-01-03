# Nexus Shim Library

Nexus Shim library provides hierarchical and relational view when interacting programmatically with Datamodel Runtime.

Datamodel is hierarchial and is defined as such in the DSL. Datamodel DSL provides syntax to specify the following, among other things:

* A definition/spec for a specific nexus node in datmodel
* The position in the hierarchy for a nexus node
* Relationships between nexus nodes

![ShimLayer](.content/images/ShimLayerUsecase.png)
## Nexus Node

Nexus Node is a first class entity in Nexus Datamodel, with following attributes:

* It has its own lifecycle
* It is a first class construct
* It is stored and mantained as a coherent data structure
* It can have relationships to other nodes
* It will expose apis/methods for clients to interact with it.
* It is a "Type" and as such objects of this Type will be created at runtime.

Nexus node, being a first class entity, applications and clients can create as many objects of this type as needed.

