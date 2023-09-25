package user

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
	tenantconfig "github.com/vmware-tanzu/graph-framework-for-microservices/api/config/tenant"
)

type User struct {
	nexus.Node

	Username  string `json:"username" yaml:"username"`
	Mail      string `json:"email,omitempty" yaml:"mail,omitempty"`
	FirstName string `json:"firstName,omitempty" yaml:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty" yaml:"lastName,omitempty"`
	Password  string `json:"password" yaml:"password"`
	TenantId  string `json:"tenantId" yaml:"tenantId"`
	Realm     string `json:"realm" yaml:"realm,omitempty"`

	Tenant tenantconfig.Tenant `nexus:"link"`
}
