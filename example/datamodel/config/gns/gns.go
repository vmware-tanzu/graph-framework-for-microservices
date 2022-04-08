package gns

import (
	"net/http"

	service_group "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/config/gns/service-group"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/config/policy"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/nexus"
)

var GNSRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestUris{
		{
			Uri:     "/v1alpha2/projects/{project}/global-namespace/{gnses.gns.tsm.tanzu.vmware.com}",
			Methods: []string{http.MethodGet, http.MethodPut, http.MethodDelete},
		},
		{
			Uri:     "/v1alpha2/projects/{project}/global-namespaces",
			Methods: []string{http.MethodGet},
		},
	},
}

var DNSRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestUris{
		{
			Uri:     "/v1alpha2/projects/{project}/dns/{dnses.gns.tsm.tanzu.vmware.com}",
			Methods: []string{http.MethodGet, http.MethodPut, http.MethodDelete},
		},
		{
			Uri:     "/v1alpha2/projects/{project}/dnses",
			Methods: []string{http.MethodGet},
		},
	},
}

// Gns struct.
// nexus-rest-api-gen:GNSRestAPISpec
// specification of GNS.
type Gns struct {
	nexus.Node
	Domain                 string
	UseSharedGateway       bool
	Description            Description
	GnsServiceGroups       map[string]service_group.SvcGroup `nexus:"child"`
	GnsAccessControlPolicy policy.AccessControlPolicy        `nexus:"child"`
	Dns                    Dns                               `nexus:"link"`
	State                  GnsState                          `nexus:"status"`
}

// This is Description struct.
type Description struct {
	Color     string
	Version   string
	ProjectId string
}

// nexus-rest-api-gen:DNSRestAPISpec
type Dns struct {
	nexus.Node
}

type GnsState struct {
	Working     bool
	Temperature int
}
