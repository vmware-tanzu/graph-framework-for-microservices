package policy

import (
	service_group "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/config/gns/service-group"
        "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
)

type AccessControlPolicy struct {
	nexus.Node
	PolicyConfigs map[string]ACPConfig `nexus:"child"`
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

	DestSvcGroups   map[string]service_group.SvcGroup `nexus:"link"`
	SourceSvcGroups map[string]service_group.SvcGroup `nexus:"link"`
	Conditions      []string
	Action          PolicyCfgActions `nexus:"@jsonencoded(file:'./root/config/policy/policy-config/policy-config.ts', gofile:'policy-config.go', name: 'PolicyCfgActions')"`
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
