package validate_test

import "fmt"

func getCRDDef(isSingleton bool) string {
	isSingletonString := fmt.Sprintf("%v", isSingleton)
	return `
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    nexus: |
      {"name":"Root.orgchart","is_singleton":` + isSingletonString + `,"children":{"employees.role.vmware.org":{"fieldName":"EmployeeRole","fieldNameGvk":"employeeRoleGvk","isNamed":false},"executives.role.vmware.org":{"fieldName":"ExecutiveRole","fieldNameGvk":"executiveRoleGvk","isNamed":false},"leaders.management.vmware.org":{"fieldName":"CEO","fieldNameGvk":"cEOGvk","isNamed":false}},"nexus-rest-api-gen":{"uris":[{"uri":"/root/{Root.orgchart}","methods":{"DELETE":{"200":{"description":"OK"},"501":{"description":"Not Implemented"}},"GET":{"200":{"description":"OK"},"404":{"description":"Not Found"},"501":{"description":"Not Implemented"}},"PUT":{"200":{"description":"OK"},"201":{"description":"Created"},"501":{"description":"Not Implemented"}}}},{"uri":"/roots","methods":{"GET":{"200":{"description":"OK"},"404":{"description":"Not Found"},"501":{"description":"Not Implemented"}}}}]}}
  creationTimestamp: null
  name: roots.orgchart.vmware.org
spec:
  conversion:
    strategy: None
  group: orgchart.vmware.org
  names:
    kind: Root
    listKind: RootList
    plural: roots
    shortNames:
    - root
    singular: root
  scope: Cluster
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            type: object
          spec:
            properties:
              cEOGvk:
                properties:
                  group:
                    type: string
                  kind:
                    type: string
                  name:
                    type: string
                required:
                - group
                - kind
                - name
                type: object
              employeeRoleGvk:
                properties:
                  group:
                    type: string
                  kind:
                    type: string
                  name:
                    type: string
                required:
                - group
                - kind
                - name
                type: object
              executiveRoleGvk:
                properties:
                  group:
                    type: string
                  kind:
                    type: string
                  name:
                    type: string
                required:
                - group
                - kind
                - name
                type: object
            type: object
        type: object
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: null
  storedVersions:
  - v1

`
}

func getCRDObject(displayName string) string {
	return `
apiVersion: orgchart.vmware.org/v1
kind: Root
metadata:
   name: someHashedName
   labels:
      nexus/is_name_hashed: "true"
      nexus/display_name: ` + displayName + `
`
}
