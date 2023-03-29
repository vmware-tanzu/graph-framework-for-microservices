package policypkg

import (
	servicegroup "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel/config/gns/service-group"
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

type AccessControlPolicy struct {
	nexus.Node
	PolicyConfigs ACPConfig `nexus:"children"`
}

// ACPConfig is a configuration of AccessControl Policy
type ACPConfig struct {
	nexus.Node
	DisplayName  string
	Gns          string
	Description  string
	Tags         []string
	ProjectId    string
	DestGroups   ResourceGroupIDs `nexus:"@jsonencoded(file:'./root/config/policy/policy-config/policy-config.ts', gofile:'policy-config.go', name: 'ResourceGroupIDs')"`
	SourceGroups ResourceGroupIDs `nexus:"@jsonencoded(file:'./root/config/policy/policy-config/policy-config.ts', gofile:'policy-config.go', name: 'ResourceGroupIDs')"`

	DestSvcGroups   servicegroup.SvcGroup `nexus:"links"` // support named children/links as map or `links` annotations
	SourceSvcGroups servicegroup.SvcGroup `nexus:"links"` // support named children/links as map or `links` annotations
	Conditions      []string
	Action          PolicyCfgActions `nexus:"@jsonencoded(file:'./root/config/policy/policy-config/policy-config.ts', gofile:'policy-config.go', name: 'PolicyCfgActions')"`
	Status          ACPStatus        `nexus:"status"`
}

type ACPStatus struct {
	StatusABC int
	StatusXYZ int
}

type ResourceGroupRef struct {
	Name string
	Type string
}

type ACPSvcGroupLinkInfo struct {
	ServiceName string
	ServiceType string
}

type PolicyActionType string

const (
	PolicyActionType_Allow  PolicyActionType = "ALLOW"
	PolicyActionType_Deny   PolicyActionType = "DENY"
	PolicyActionType_Log    PolicyActionType = "LOG"
	PolicyActionType_Mirror PolicyActionType = "MIRROR"
)

type PolicyCfgAction struct {
	Action PolicyActionType `json:"action" mapstructure:"action"`
}

type PolicyCfgActions []PolicyCfgAction

type ResourceGroupID struct {
	Name string `json:"name" mapstruction:"name"`
	Type string `json:"type" mapstruction:"type"`
}

type ResourceGroupIDs []ResourceGroupID

// nexus-graphql-query:gns.CloudEndpointGraphQLQuerySpec
type VMpolicy struct {
	nexus.Node
}
