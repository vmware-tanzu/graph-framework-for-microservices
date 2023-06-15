package validate_test

import (
	"fmt"
)

func getRootCRDDef(isSingleton bool) string {
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

func getEmployeeCRDDef() string {
	return `
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    nexus: |
      {"name":"role.Employee","hierarchy":["roots.orgchart.vmware.org"],"is_singleton":false,"nexus-rest-api-gen":{"uris":[{"uri":"/root/{orgchart.Root}/employee/{role.Employee}","methods":{"DELETE":{"200":{"description":"OK"},"501":{"description":"Not Implemented"}},"GET":{"200":{"description":"OK"},"404":{"description":"Not Found"},"501":{"description":"Not Implemented"}},"PUT":{"200":{"description":"OK"},"201":{"description":"Created"},"501":{"description":"Not Implemented"}}}},{"uri":"/employees","methods":{"LIST":{"200":{"description":"OK"},"404":{"description":"Not Found"},"501":{"description":"Not Implemented"}}}}]}}
  creationTimestamp: null
  name: employees.role.vmware.org
spec:
  conversion:
    strategy: None
  group: role.vmware.org
  names:
    kind: Employee
    listKind: EmployeeList
    plural: employees
    shortNames:
    - employee
    singular: employee
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
          status:
            properties:
              nexus:
                properties:
                  remoteGeneration:
                    format: int64
                    type: integer
                  sourceGeneration:
                    format: int64
                    type: integer
                required:
                - sourceGeneration
                - remoteGeneration
                type: object
            type: object
        type: object
    served: true
    storage: true
    subresources:
      status: {}
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: null
  storedVersions:
  - v1


`
}

func getRootCRDObject(displayName string) string {
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

func getEmployeeCRDObject(displayName string, extraLabels map[string]string) string {
	var extraLabelsString string
	for k, v := range extraLabels {
		extraLabelsString += fmt.Sprintf("      %s: %s\n", k, v)
	}

	return `
apiVersion: role.vmware.org/v1
kind: Employee
metadata:
   name: someHashedName
   labels:
      nexus/is_name_hashed: "true"
      nexus/display_name: ` + displayName + `
` + extraLabelsString + `
`
}

func getAloneCRDObject() string {
	return `
apiVersion: orgchart.vmware.org/v1
kind: Alone
metadata:
   name: someHashedName
   labels:
      nexus/is_name_hashed: "true"
      nexus/display_name: "foo"
`
}
