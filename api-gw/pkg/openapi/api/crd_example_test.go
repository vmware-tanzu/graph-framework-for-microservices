package api_test

var crdExample = `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    nexus: |
      {"name":"management.Leader","hierarchy":["roots.orgchart.vmware.org"],"children":{"humanresourceses.hr.vmware.org":{"fieldName":"HR","fieldNameGvk":"hRGvk","isNamed":false},"mgrs.management.vmware.org":{"fieldName":"EngManagers","fieldNameGvk":"engManagersGvk","isNamed":true}},"links":{"Role":{"fieldName":"Role","fieldNameGvk":"roleGvk","isNamed":false}},"is_singleton":true,"nexus-rest-api-gen":{"uris":[{"uri":"/root/{orgchart.Root}/leader/{management.Leader}","methods":{"DELETE":{"200":{"description":"OK"},"404":{"description":"Not Found"},"501":{"description":"Not Implemented"}},"GET":{"200":{"description":"OK"},"404":{"description":"Not Found"},"501":{"description":"Not Implemented"}},"PUT":{"200":{"description":"OK"},"201":{"description":"Created"},"501":{"description":"Not Implemented"}}},"auth":false},{"uri":"/leader","methods":{"DELETE":{"200":{"description":"OK"},"404":{"description":"Not Found"},"501":{"description":"Not Implemented"}},"GET":{"200":{"description":"OK"},"404":{"description":"Not Found"},"501":{"description":"Not Implemented"}},"PUT":{"200":{"description":"OK"},"201":{"description":"Created"},"501":{"description":"Not Implemented"}}},"auth":false},{"uri":"/leaders","methods":{"LIST":{"200":{"description":"OK"},"404":{"description":"Not Found"},"501":{"description":"Not Implemented"}}},"auth":false}]}}
  creationTimestamp: null
  name: leaders.management.vmware.org
spec:
  conversion:
    strategy: None
  group: management.vmware.org
  names:
    kind: Leader
    listKind: LeaderList
    plural: leaders
    shortNames:
    - leader
    singular: leader
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
              designation:
                type: string
              employeeID:
                format: int32
                type: integer
              engManagersGvk:
                additionalProperties:
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
              hRGvk:
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
              name:
                type: string
              roleGvk:
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
            required:
            - designation
            - name
            - employeeID
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
              status:
                properties:
                  DaysLeftToEndOfVacations:
                    format: int32
                    type: integer
                  IsOnVacations:
                    type: boolean
                required:
                - IsOnVacations
                - DaysLeftToEndOfVacations
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
