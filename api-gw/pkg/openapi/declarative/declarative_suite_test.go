package declarative_test

import (
	"testing"

	log "github.com/sirupsen/logrus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	Uri         = "/v1alpha1/project/{projectId}/global-namespaces"
	ResourceUri = "/v1alpha1/project/{projectId}/global-namespaces/{id}"
	ListUri     = "/v1alpha1/global-namespaces/test"
)

var (
	spec = []byte(`openapi: 3.0.0
info:
  version: 1.0.0
  title: NSX-SM <Tenant/Operator> APIs
  description: <Tenant/Operator> APIs for NSX service mesh.
  termsOfService: 'http://nsxservicemesh.vmware.com/terms/'
  contact:
    name: VMware NSX-ServiceMesh Team
    email: support@nsxservicemesh.vmware.com
    url: 'http://nsxservicemesh.vmware.com/'
  license:
    name: VMWare
    url: 'https://nsxservicemesh.vmware.com/licenses/LICENSE.html'
servers:
  - url: 'http://127.0.0.1:3000'
basePath: /v1
paths:
  '/v1alpha1/project/{projectId}/global-namespaces/{id}':
    put:
      x-controller-name: GlobalNamespaceControllerV1Alpha1
      x-operation-name: putGlobalNamespaceV1
      tags:
        - Global Namespaces (v1alpha1)
      description: Create the global namespace
      x-nexus-kind-name: GlobalNamespace
      x-nexus-group-name: gns.vmware.org
      x-nexus-identifier: id
      x-nexus-short-name: gns
      responses:
        '200':
          description: 'global namespace updated '
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/GlobalNamespaceConfig'
        '201':
          description: 'global namespace created '
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/GlobalNamespaceConfig'
        default:
          description: unexpected error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ApiHttpError'
      parameters:
        - name: id
          in: path
          schema:
            type: string
          required: true
      requestBody:
        description: Global namespace Config
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/GlobalNamespaceConfig'
        x-parameter-index: 1
      operationId: GlobalNamespaceControllerV1Alpha1.putGlobalNamespaceV1
    get:
      x-controller-name: GlobalNamespaceControllerV1Alpha1
      x-operation-name: getGlobalNamespaceV1
      x-nexus-identifier: id
      tags:
        - Global Namespaces (v1alpha1)
      description: Return the config for a global namespace
      x-nexus-kind-name: GlobalNamespace
      x-nexus-group-name: gns.vmware.org
      x-nexus-short-name: gns
      responses:
        '200':
          description: global namespace config
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/GlobalNamespaceConfig'
        default:
          description: unexpected error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ApiHttpError'
      parameters:
        - name: id
          in: path
          schema:
            type: string
          required: true
      operationId: GlobalNamespaceControllerV1Alpha1.getGlobalNamespaceV1
    delete:
      x-controller-name: GlobalNamespaceControllerV1Alpha1
      x-operation-name: deleteGlobalNamespaceV1
      tags:
        - Global Namespaces (v1alpha1)
      description: Delete the global namespace
      x-nexus-kind-name: GlobalNamespace
      x-nexus-group-name: gns.vmware.org
      x-nexus-identifier: id
      x-nexus-short-name: gns
      responses:
        '200':
          description: 'global namespace delete '
          content:
            application/json:
              schema:
                type: string
        default:
          description: unexpected error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ApiHttpError'
      parameters:
        - name: id
          in: path
          schema:
            type: string
          required: true
      operationId: GlobalNamespaceControllerV1Alpha1.deleteGlobalNamespaceV1
  /v1alpha1/project/{projectId}/global-namespaces:
    post:
      x-controller-name: GlobalNamespaceControllerV1Alpha1
      x-operation-name: postGlobalNamespaceV1
      tags:
        - Global Namespaces (v1alpha1)
      description: Create the global namespace
      responses:
        '200':
          description: 'global namespace created '
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/GlobalNamespaceConfig'
        default:
          description: unexpected error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ApiHttpError'
      requestBody:
        description: Global namespace Config
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/GlobalNamespaceConfig'
      operationId: GlobalNamespaceControllerV1Alpha1.postGlobalNamespaceV1
    get:
      x-controller-name: GlobalNamespaceControllerV1Alpha1
      x-operation-name: getGNSListV1
      tags:
        - Global Namespaces (v1alpha1)
      description: Get a list of GNS IDs that are defined
      x-nexus-kind-name: GlobalNamespace
      x-nexus-group-name: gns.vmware.org
      x-nexus-short-name: gns
      responses:
        '200':
          description: list of gns defined
          content:
            application/json:
              schema:
                type: array
                items:
                  type: string
        default:
          description: unexpected error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ApiHttpError'
      operationId: GlobalNamespaceControllerV1Alpha1.getGNSListV1
  /v1alpha1/global-namespaces/test:
    get:
      x-controller-name: GlobalNamespaceControllerV1Alpha1
      x-operation-name: getGNSListV1
      tags:
        - Global Namespaces (v1alpha1)
      description: Get a list of GNS IDs that are defined
      x-nexus-kind-name: GlobalNamespaceList
      x-nexus-group-name: gns.vmware.org
      x-nexus-list-endpoint: true
      x-nexus-short-name: gns
      responses:
        '200':
          description: list of gns defined
          content:
            application/json:
              schema:
                type: array
                items:
                  type: string
        default:
          description: unexpected error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ApiHttpError'
      operationId: GlobalNamespaceControllerV1Alpha1.getGNSListV1
components:
  schemas:
    ApiHttpError:
      title: ApiHttpError
      properties:
        code:
          type: number
        message:
          type: string
      required:
        - code
        - message
      additionalProperties: false
    Condition:
      title: Condition
      properties:
        type:
          type: string
          description: START_WITH | EXACT
        match:
          type: string
      additionalProperties: false
    MatchCondition:
      title: MatchCondition
      properties:
        namespace:
          $ref: '#/components/schemas/Condition'
        cluster:
          $ref: '#/components/schemas/Condition'
        service: 
          $ref: '#/components/schemas/MatchCondition'
      required:
        - namespace
      additionalProperties: false
    GlobalNamespaceConfig:
      title: GlobalNamespaceConfig
      properties:
        name:
          type: string
          pattern: '^[a-z0-9][a-z0-9-.]*[a-z0-9]$'
          minLength: 2
          maxLength: 253
        display_name:
          type: string
        domain_name:
          type: string
        use_shared_gateway:
          type: boolean
        mtls_enforced:
          type: boolean
        ca_type:
          type: string
          enum:
            - PreExistingCA
            - self-signed
        ca:
          type: string
        description:
          type: string
        color:
          type: string
        version:
          type: string
        match_conditions:
          type: array
          items:
            $ref: '#/components/schemas/MatchCondition'
        api_discovery_enabled:
          type: boolean
      required:
        - name
        - domain_name
        - match_conditions
      additionalProperties: false`)
)

func TestDeclarative(t *testing.T) {
	log.StandardLogger().ExitFunc = nil
	RegisterFailHandler(Fail)
	RunSpecs(t, "Declarative Suite")
}
