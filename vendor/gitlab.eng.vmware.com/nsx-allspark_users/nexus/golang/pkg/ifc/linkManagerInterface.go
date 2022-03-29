package ifc

type BaseNodeLinkManagerInterface interface {
	UpsertType(nodeType string) map[string]*GLink
	UpsertRLinkType(nodeType string) map[string]*GLink
	HasType(nodeType string) bool
	HasRLinkType(nodeType string) bool
	Has(nodeType, childKey string) bool
	HasRLink(nodeType, childKey string) bool
	Get(nodeType, childKey string) (*GLink, bool)
	SetLinkProperty(nodeType, childKey string, prop PropertyType) bool
	Find(nodeType string, fn func(lnk *GLink) bool) (string, bool)
	Add(nodeType, nodeKey string, lnk *GLink)
	AddRLink(nodeType, nodeKey string, lnk *GLink)
	Delete(nodeType, nodeKey string)
	DeleteRLink(nodeType, nodeKey string)
	ForEach(fn func(ntype, nkey string, lnk *GLink))
	ForEachType(ntype string, fn func(ntype, nkey string, lnk *GLink))
	ForEachRLink(fn func(ntype, nkey string, lnk *GLink))
	ForEachRLinkType(ntype string, fn func(ntype, nkey string, lnk *GLink))
	GetNextKey(nodeType, key string) string
	GetRLink(nodeType, childKey string) (*GLink, bool)         // rlinks
	GetRLinksForNodeType(nodeType string) ([]*GLink, bool)     // rlinks
	GetRLinks() ([]*GLink, bool)                               // rlinks
	GetRLinkId(nodeType string, nValue string) (string, bool)  // rlinks
	GetRLinkIdsForNodeType(nodeType string) ([]string, bool)   // rlinks
	GetRLinkIds() ([]string, bool)                             //rlinks
	GetChildLinksForNodeType(nodeType string) ([]*GLink, bool) // hlinks
	GetChildLinks() ([]*GLink, bool)                           // hlinks
	GetChildLink(nodeType, childKey string) (*GLink, bool)     // hlinks
}
