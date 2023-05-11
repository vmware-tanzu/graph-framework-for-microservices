# authz-controller

Main responsibility of authz-controller is to monitor:

1. ResourceRole objects created in the cluster and creates the equivalent k8s ClusterRole.
2. ResourceRoleBinding objects created in the cluster and creates the equivalent k8s ClusterRoleBinding.
3. User objects created in the cluster and creates the certificate signed by ca and store it in UserCertificate CRD.
4. CRD created in the cluster and verify if the CRD type of interest, simply add it to the role object.

## Sample Nexus RBAC Object

```
apiVersion: authorization.nexus.org/v1
kind: ResourceRole
metadata:
  name: root-admin-role
spec:
  rules:
  - Hierarchical: false
    Resource:
      group: apix.mazinger.com
      kind: Root
    Verbs:
    - get
    - list
  - Hierarchical: true
    Resource:
      group: config.mazinger.com
      kind: ApiCollaborationSpace
    Verbs:
    - get
    - list
---

apiVersion: authorization.nexus.org/v1
kind: User
metadata:
  name: bob
---

apiVersion: authorization.nexus.org/v1
kind: ResourceRoleBinding
metadata:
  name: root-admin-role-binding
spec:
  roleGvk:
    group: authorization.nexus.org
    kind:  ResourceRole
    name:  11db4dfc940481cd1030e7aa1aaf6284b63be65b
  usersGvk:
    bob:
      group: authorization.nexus.org
      kind: User
      name: "bob"
```
