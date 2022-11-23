package nexus_compare

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Compare lib tests", func() {

	It("should return false when no changes", func() {
		ans, _, err := CompareFiles([]byte(baseSpec), []byte(baseSpec))
		Expect(err).NotTo(HaveOccurred())
		Expect(ans).To(BeFalse())
	})
	It("should return true and report change type", func() {
		ans, text, err := CompareFiles([]byte(baseSpec), []byte(changeType))

		Expect(err).NotTo(HaveOccurred())
		Expect(ans).To(BeTrue())
		changeCheck := []string{"/spec/versions/name=v1/schema/openAPIV3Schema/properties/spec/properties/name/type", "value change", "- string", "+ int"}
		for _, v := range changeCheck {
			Expect(text.String()).Should(ContainSubstring(v))
		}
	})
	It("should return true and report field deletion", func() {
		ans, text, err := CompareFiles([]byte(baseSpec), []byte(fieldDeletion))
		Expect(err).NotTo(HaveOccurred())
		Expect(ans).To(BeTrue())
		changeCheck := []string{"/spec/versions/name=v1/schema/openAPIV3Schema/properties/status/properties/nexus/properties", "one field removed", "remoteGeneration"}
		for _, v := range changeCheck {
			Expect(text.String()).Should(ContainSubstring(v))
		}
	})
	It("should return false when adding field", func() {
		ans, text, err := CompareFiles([]byte(baseSpec), []byte(addedField))
		Expect(err).NotTo(HaveOccurred())
		Expect(ans).To(BeFalse())
		Expect(text.String()).To(BeEmpty())
	})
	It("should report change in nexus annotation", func() {
		ans, text, err := CompareFiles([]byte(baseSpec), []byte(changeAnnotation))
		Expect(err).NotTo(HaveOccurred())
		Expect(ans).To(BeTrue())
		changeCheck := []string{"nexus annotation changes", "/is_singleton", "value change"}
		for _, v := range changeCheck {
			Expect(text.String()).Should(ContainSubstring(v))
		}
	})
	It("should report no changes in nexus annotation with added field, singleton to false and change in api", func() {
		ans, _, err := CompareFiles([]byte(other), []byte(other2))
		Expect(err).NotTo(HaveOccurred())
		Expect(ans).To(BeFalse())
	})
	It("should report changes in nexus annotation with deleted field, singleton to true", func() {
		ans, text, err := CompareFiles([]byte(other2), []byte(other))
		Expect(err).NotTo(HaveOccurred())
		Expect(ans).To(BeTrue())
		changeCheck := []string{"nexus annotation changes", "/children", "accesscontrolpolicies.policypkg.tsm.tanzu.vmware.com:", "/is_singleton", "value change"}
		for _, v := range changeCheck {
			Expect(text.String()).Should(ContainSubstring(v))
		}
	})
})

var baseSpec = `
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    nexus: |
      {"name":"gns.IgnoreChild","hierarchy":["roots.root.tsm.tanzu.vmware.com","configs.config.tsm.tanzu.vmware.com","gnses.gns.tsm.tanzu.vmware.com"],"is_singleton":false,"nexus-rest-api-gen":{"uris":null}}
  creationTimestamp: null
  name: ignorechilds.gns.tsm.tanzu.vmware.com
spec:
  conversion:
    strategy: None
  group: gns.tsm.tanzu.vmware.com
  versions:
    - name: v1
      schema:
        openAPIV3Schema:
          properties:
            spec:
              properties:
                name:
                  type: string
              required:
                - name
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
`
var changeType = `
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    nexus: |
      {"name":"gns.IgnoreChild","hierarchy":["roots.root.tsm.tanzu.vmware.com","configs.config.tsm.tanzu.vmware.com","gnses.gns.tsm.tanzu.vmware.com"],"is_singleton":false,"nexus-rest-api-gen":{"uris":null}}
  creationTimestamp: null
  name: ignorechilds.gns.tsm.tanzu.vmware.com
spec:
  conversion:
    strategy: None
  group: gns.tsm.tanzu.vmware.com
  versions:
    - name: v1
      schema:
        openAPIV3Schema:
          properties:
            spec:
              properties:
                name:
                  type: int
              required:
                - name
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
`
var fieldDeletion = `
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    nexus: |
      {"name":"gns.IgnoreChild","hierarchy":["roots.root.tsm.tanzu.vmware.com","configs.config.tsm.tanzu.vmware.com","gnses.gns.tsm.tanzu.vmware.com"],"is_singleton":false,"nexus-rest-api-gen":{"uris":null}}
  creationTimestamp: null
  name: ignorechilds.gns.tsm.tanzu.vmware.com
spec:
  conversion:
    strategy: None
  group: gns.tsm.tanzu.vmware.com
  versions:
    - name: v1
      schema:
        openAPIV3Schema:
          properties:
            spec:
              properties:
                name:
                  type: string
              required:
                - name
              type: object
            status:
              properties:
                nexus:
                  properties:
                    sourceGeneration:
                      format: int64
                      type: integer
                  required:
                    - sourceGeneration
                    - remoteGeneration
                  type: object
              type: object
          type: object
`
var addedField = `
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    nexus: |
      {"name":"gns.IgnoreChild","hierarchy":["roots.root.tsm.tanzu.vmware.com","configs.config.tsm.tanzu.vmware.com","gnses.gns.tsm.tanzu.vmware.com"],"is_singleton":false,"nexus-rest-api-gen":{"uris":null}}
  creationTimestamp: null
  name: ignorechilds.gns.tsm.tanzu.vmware.com
spec:
  conversion:
    strategy: None
  group: gns.tsm.tanzu.vmware.com
  versions:
    - name: v1
      schema:
        openAPIV3Schema:
          properties:
            spec:
              properties:
                name:
                  type: string
              required:
                - name
              type: object
            status:
              properties:
                nexus:
                  properties:
                    remoteGeneration:
                      format: int64
                      type: integer
                    addedField:
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
`
var changeAnnotation = `
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    nexus: |
      {"name":"gns.IgnoreChild","hierarchy":["roots.root.tsm.tanzu.vmware.com","configs.config.tsm.tanzu.vmware.com","gnses.gns.tsm.tanzu.vmware.com"],"is_singleton":true,"nexus-rest-api-gen":{"uris":null}}
  creationTimestamp: null
  name: ignorechilds.gns.tsm.tanzu.vmware.com
spec:
  conversion:
    strategy: None
  group: gns.tsm.tanzu.vmware.com
  versions:
    - name: v1
      schema:
        openAPIV3Schema:
          properties:
            spec:
              properties:
                name:
                  type: string
              required:
                - name
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
`

var other = `
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    nexus: |
      {"name":"gns.Gns","hierarchy":["roots.root.tsm.tanzu.vmware.com","configs.config.tsm.tanzu.vmware.com"],"children":{"barchilds.gns.tsm.tanzu.vmware.com":{"fieldName":"FooChild","fieldNameGvk":"fooChildGvk","isNamed":false},"foos.gns.tsm.tanzu.vmware.com":{"fieldName":"Foo","fieldNameGvk":"fooGvk","isNamed":false},"ignorechilds.gns.tsm.tanzu.vmware.com":{"fieldName":"IgnoreChild","fieldNameGvk":"ignoreChildGvk","isNamed":false},"svcgroups.servicegroup.tsm.tanzu.vmware.com":{"fieldName":"GnsServiceGroups","fieldNameGvk":"gnsServiceGroupsGvk","isNamed":true}},"links":{"Dns":{"fieldName":"Dns","fieldNameGvk":"dnsGvk","isNamed":false}},"is_singleton":true,"nexus-rest-api-gen":{"uris":[{"uri":"/v1alpha2/global-namespace/{gns.Gns}","query_params":["config.Config"],"methods":{"DELETE":{"200":{"description":"OK"},"404":{"description":"Not Found"},"501":{"description":"Not Implemented"}},"GET":{"200":{"description":"OK"},"404":{"description":"Not Found"},"501":{"description":"Not Implemented"}},"PUT":{"200":{"description":"OK"},"201":{"description":"Created"},"501":{"description":"Not Implemented"}}}},{"uri":"/v1alpha2/global-namespaces","query_params":["config.Config"],"methods":{"LIST":{"200":{"description":"OK"},"404":{"description":"Not Found"},"501":{"description":"Not Implemented"}}}},{"uri":"/test-foo","query_params":["config.Config"],"methods":{"DELETE":{"200":{"description":"ok"},"404":{"description":"Not Found"},"501":{"description":"Not Implemented"}}}},{"uri":"/test-bar","query_params":["config.Config"],"methods":{"PATCH":{"400":{"description":"Bad Request"}}}}]},"description":"this is my awesome node"}
  creationTimestamp: null
  name: gnses.gns.tsm.tanzu.vmware.com
spec:
  conversion:
    strategy: None
  group: gns.tsm.tanzu.vmware.com
  names:
    kind: Gns
    listKind: GnsList
    plural: gnses
    shortNames:
    - gns
    singular: gns
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
              description:
                properties:
                  Color:
                    type: string
                  HostPort:
                    properties:
                      Host:
                        type: string
                      Port:
                        format: int32
                        type: integer
                    required:
                    - Host
                    - Port
                    type: object
                  Instance:
                    format: float
                    type: number
                  ProjectId:
                    type: string
                  TestAns:
                    items:
                      properties:
                        Name:
                          type: string
                      required:
                      - Name
                      type: object
                    type: array
                  Version:
                    type: string
                required:
                - Color
                - Version
                - ProjectId
                - TestAns
                - Instance
                - HostPort
                type: object
              differentSpec:
                type: object
                x-kubernetes-preserve-unknown-fields: true
              dnsGvk:
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
              domain:
                maxLength: 8
                minLength: 2
                pattern: abc
                type: string
              fooChildGvk:
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
              fooGvk:
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
              gnsAccessControlPolicyGvk:
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
              gnsServiceGroupsGvk:
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
              ignoreChildGvk:
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
              mapPointer:
                additionalProperties:
                  type: string
                type: object
              meta:
                type: string
              otherDescription:
                properties:
                  Color:
                    type: string
                  HostPort:
                    properties:
                      Host:
                        type: string
                      Port:
                        format: int32
                        type: integer
                    required:
                    - Host
                    - Port
                    type: object
                  Instance:
                    format: float
                    type: number
                  ProjectId:
                    type: string
                  TestAns:
                    items:
                      properties:
                        Name:
                          type: string
                      required:
                      - Name
                      type: object
                    type: array
                  Version:
                    type: string
                required:
                - Color
                - Version
                - ProjectId
                - TestAns
                - Instance
                - HostPort
                type: object
              port:
                format: int32
                type: integer
              serviceSegmentRef:
                properties:
                  Field1:
                    type: string
                  Field2:
                    type: string
                required:
                - Field1
                - Field2
                type: object
              serviceSegmentRefMap:
                additionalProperties:
                  properties:
                    Field1:
                      type: string
                    Field2:
                      type: string
                  required:
                  - Field1
                  - Field2
                  type: object
                type: object
              serviceSegmentRefPointer:
                properties:
                  Field1:
                    type: string
                  Field2:
                    type: string
                required:
                - Field1
                - Field2
                type: object
              serviceSegmentRefs:
                items:
                  properties:
                    Field1:
                      type: string
                    Field2:
                      type: string
                  required:
                  - Field1
                  - Field2
                  type: object
                type: array
              slicePointer:
                items:
                  type: string
                type: array
              useSharedGateway:
                type: boolean
              workloadSpec:
                type: object
                x-kubernetes-preserve-unknown-fields: true
            required:
            - domain
            - useSharedGateway
            - description
            - meta
            - port
            - otherDescription
            - mapPointer
            - slicePointer
            - workloadSpec
            - differentSpec
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
              state:
                properties:
                  Temperature:
                    format: int32
                    type: integer
                  Working:
                    type: boolean
                required:
                - Working
                - Temperature
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

var other2 = `
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    nexus: |
      {"name":"gns.Gns","hierarchy":["roots.root.tsm.tanzu.vmware.com","configs.config.tsm.tanzu.vmware.com"],"children":{"accesscontrolpolicies.policypkg.tsm.tanzu.vmware.com":{"fieldName":"GnsAccessControlPolicy","fieldNameGvk":"gnsAccessControlPolicyGvk","isNamed":false},"barchilds.gns.tsm.tanzu.vmware.com":{"fieldName":"FooChild","fieldNameGvk":"fooChildGvk","isNamed":false},"foos.gns.tsm.tanzu.vmware.com":{"fieldName":"Foo","fieldNameGvk":"fooGvk","isNamed":false},"ignorechilds.gns.tsm.tanzu.vmware.com":{"fieldName":"IgnoreChild","fieldNameGvk":"ignoreChildGvk","isNamed":false},"svcgroups.servicegroup.tsm.tanzu.vmware.com":{"fieldName":"GnsServiceGroups","fieldNameGvk":"gnsServiceGroupsGvk","isNamed":true}},"links":{"Dns":{"fieldName":"Dns","fieldNameGvk":"dnsGvk","isNamed":false}},"is_singleton":false,"nexus-rest-api-gen":{"uris":[{"uri":"/v1alpha2/global-namespace/{gns.Gns}","query_params":["config.Config"],"methods":{"DELETE":{"200":{"description":"OK"},"404":{"description":"Not Found"},"501":{"description":"Not Implemented"}},"GET":{"200":{"description":"OK"},"404":{"description":"Not Found"},"501":{"description":"Not Implemented"}},"PUT":{"200":{"description":"OK"},"201":{"description":"Created"},"501":{"description":"Not Implemented"}}}},{"uri":"/v1alpha2/global-namespaces","query_params":["config.Config"],"methods":{"LIST":{"200":{"description":"OK"},"404":{"description":"Not Found"},"501":{"description":"Not Implemented"}}}},{"uri":"/test-foo22","query_params":["config.Config"],"methods":{"DELETE":{"200":{"description":"ok"},"404":{"description":"Not Found"},"501":{"description":"Not Implemented"}}}},{"uri":"/test-bar","query_params":["config.Config"],"methods":{"PATCH":{"400":{"description":"Bad Request"}}}}]},"description":"this is my awesome node"}
  creationTimestamp: null
  name: gnses.gns.tsm.tanzu.vmware.com
spec:
  conversion:
    strategy: None
  group: gns.tsm.tanzu.vmware.com
  names:
    kind: Gns
    listKind: GnsList
    plural: gnses
    shortNames:
    - gns
    singular: gns
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
              description:
                properties:
                  Color:
                    type: string
                  HostPort:
                    properties:
                      Host:
                        type: string
                      Port:
                        format: int32
                        type: integer
                    required:
                    - Host
                    - Port
                    type: object
                  Instance:
                    format: float
                    type: number
                  ProjectId:
                    type: string
                  TestAns:
                    items:
                      properties:
                        Name:
                          type: string
                      required:
                      - Name
                      type: object
                    type: array
                  Version:
                    type: string
                required:
                - Color
                - Version
                - ProjectId
                - TestAns
                - Instance
                - HostPort
                type: object
              differentSpec:
                type: object
                x-kubernetes-preserve-unknown-fields: true
              dnsGvk:
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
              domain:
                maxLength: 8
                minLength: 2
                pattern: abc
                type: string
              fooChildGvk:
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
              fooGvk:
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
              gnsAccessControlPolicyGvk:
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
              gnsServiceGroupsGvk:
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
              ignoreChildGvk:
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
              mapPointer:
                additionalProperties:
                  type: string
                type: object
              meta:
                type: string
              otherDescription:
                properties:
                  Color:
                    type: string
                  HostPort:
                    properties:
                      Host:
                        type: string
                      Port:
                        format: int32
                        type: integer
                    required:
                    - Host
                    - Port
                    type: object
                  Instance:
                    format: float
                    type: number
                  ProjectId:
                    type: string
                  TestAns:
                    items:
                      properties:
                        Name:
                          type: string
                      required:
                      - Name
                      type: object
                    type: array
                  Version:
                    type: string
                required:
                - Color
                - Version
                - ProjectId
                - TestAns
                - Instance
                - HostPort
                type: object
              port:
                format: int32
                type: integer
              serviceSegmentRef:
                properties:
                  Field1:
                    type: string
                  Field2:
                    type: string
                required:
                - Field1
                - Field2
                type: object
              serviceSegmentRefMap:
                additionalProperties:
                  properties:
                    Field1:
                      type: string
                    Field2:
                      type: string
                  required:
                  - Field1
                  - Field2
                  type: object
                type: object
              serviceSegmentRefPointer:
                properties:
                  Field1:
                    type: string
                  Field2:
                    type: string
                required:
                - Field1
                - Field2
                type: object
              serviceSegmentRefs:
                items:
                  properties:
                    Field1:
                      type: string
                    Field2:
                      type: string
                  required:
                  - Field1
                  - Field2
                  type: object
                type: array
              slicePointer:
                items:
                  type: string
                type: array
              useSharedGateway:
                type: boolean
              workloadSpec:
                type: object
                x-kubernetes-preserve-unknown-fields: true
            required:
            - domain
            - useSharedGateway
            - description
            - meta
            - port
            - otherDescription
            - mapPointer
            - slicePointer
            - workloadSpec
            - differentSpec
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
              state:
                properties:
                  Temperature:
                    format: int32
                    type: integer
                  Working:
                    type: boolean
                required:
                - Working
                - Temperature
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
