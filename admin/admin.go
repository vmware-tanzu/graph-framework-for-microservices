package admin

import "golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"

// Upstream defines an address to proxy the request to
type Upstream struct {
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
	Port   uint32 `json:"port"`
}

// MatchType configuration that determines where the admin gateway has to look for a match condition
type MatchType string

var (
	Jwt    MatchType = "jwt"
	Header MatchType = "header"
)

type MatchCondition struct {
	Type  MatchType `json:"type"`
	Key   string    `json:"key"`
	Value string    `json:"value"`
}

type ProxyRule struct {
	nexus.Node

	// Information about what part of the request must be matched
	MatchCondition MatchCondition `json:"matchCondition"`

	// If the match condition is satisfied, the namespace of the tenant api-gw we will proxy to.
	Upstream Upstream `json:"upstream"`
}
