package root

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/config"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/nexus"
)

// nexus-rest-api-gen:
//   uri: /v1alpha1/projects/$PID/global-namespace/$GID
//   methods:
//     - GET
//     - PUT
//     - DELETE
//   response:
//     200:
//       message: success message
//     400:
//       message: not found message
//     401:
//       message: unauthorized message
// nexus-api-validation-endpoint:
//   - service: service-name
//     endpoint: /foo/bar
//   - service: service-name-2
//     endpoint: /foo/bar
// nexus-version: v1
type Root struct {
	nexus.Node
	Config config.Config `nexus:"child"`
}
