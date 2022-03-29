package ifc

type BaseNodeLinkManagerIteratorInterface interface {
	GetNodeId() string
	GetNodeType() string
	GetNodeKeyValue() string
	GetNodeParentNodeKey() string
	GetNodeParentNodeType() string
	GetParentNodeIterator() BaseNodeLinkManagerIteratorInterface
	ForEachCachedRLinkWithNodeType(fn func(rLinkHdlr BaseNodeLinkManagerIteratorInterface) interface{},
		nodeType string)
	ForEachCachedRLink(fn func(rLinkHdlr BaseNodeLinkManagerIteratorInterface) interface{})
	GetNodeProperties() PropertyType
	GetCachedRLinkHdlr(nodeType, nodeKey string) BaseNodeLinkManagerIteratorInterface
	IsRLinkFeatureFlagEnabled() bool
	GetLinkProperties() (PropertyType, bool)
	ForEachCachedChildLinkWithNodeType(fn func(rLinkHdlr BaseNodeLinkManagerIteratorInterface) interface{},
		nodeType string)
	ForEachCachedChildLink(fn func(rLinkHdlr BaseNodeLinkManagerIteratorInterface) interface{})
	GetCachedChildLinkHdlr(nodeType, nodeKey string) BaseNodeLinkManagerIteratorInterface
}
