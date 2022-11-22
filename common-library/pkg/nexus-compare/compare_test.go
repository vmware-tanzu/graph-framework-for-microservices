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
		changeCheck := []string{"/spec/versions/name=v1/schema/openAPIV3Schema/properties/status/properties/nexus/properties", "one map entry removed", "remoteGeneration"}
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
