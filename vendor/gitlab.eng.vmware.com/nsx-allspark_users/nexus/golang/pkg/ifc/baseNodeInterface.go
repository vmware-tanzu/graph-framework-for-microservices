package ifc

import (
	"sync"
)

// BaseNodeInterface defines what a base node is and how it should behave.
// We require that any BaseNode type that follows this contract embed a mutex as synchronization is often needed.
type BaseNodeInterface interface {
	sync.Locker

	GetCCnt() uint64
	SetDeleted()
	SetDeleteCompleted()
	IsDeleteCompleted() bool
	GetId() string
	GetType() string
	SetType(t string)
	IsDeleted() bool
	GetKeyValue() string
	GetKeyName() string
	SetKeyName(n string)
	SetKeyValue(v string)
	IsTypeRefNode() bool
	AddReverseLink(nodeId, linkKeyValue, linkId string)
	RemoveReverseLink(nodeId string)
	ReverseLinkIterate(fn func(linkKeyValue string, nodeId string, linkId string))
	//	IsIDPath() bool
	// all the links and rlinks are stored
	GetLinks() BaseNodeLinkManagerInterface
	// immediate task, just update the local data structure. no update to db
	SetImmediateBaseProperties(pnew PropertyType)
	GetImmediateBaseProperties() PropertyType
	DeleteImmediateLink(nodeType, nodeKey string)
	UpsertImmediateLink(nodeType, nodeKey string, lnk *GLink)
	DeleteImmediateRLink(nodeType, nodeKey string)
	UpsertImmediateRLink(nodeType, nodeKey string, lnk *GLink)
	// update are form the db
	SetBaseProperties(pnew PropertyType, merge bool)
	DelBaseProperties(keys []string)
	GetBaseProperties() (PropertyType, bool)                                                     // promise
	GetChild(childKey, nodeType string, forcedSync bool) (BaseNodeInterface, PropertyType, bool) // promise
	//: Promise<[undefined | BaseNodeInterface, undefined | { [key: string]: any }]>;
	GetCachedChild(childKey, nodeType string) (BaseNodeInterface, PropertyType, bool)

	// returns the object and it's parent
	GetSoftLinkedChildLeg(childKey, nodeType string) []BaseNodeInterface // promise

	LinkAddProperty(childKey, nodeType string, linkProp PropertyType) bool   // promsie
	LinkRemoveProperty(childKey, nodeType string, linkPropKey []string) bool // promise
	Sync(forcedSync bool)                                                    // promise
	GetNextChildKey(nodeName, nodeType string) string
	GetNextChildKeyList(nodeName string, nodeType string, cnt uint32) []string      // promise
	GetNextRLinkChildKeyList(nodeName string, nodeType string, cnt uint32) []string // promise
	//    registerBaseCB(cbfn: Function);
	GetBaseParent() BaseNodeInterface
	SetBaseParent(n BaseNodeInterface)
	GetChildLink(childKey, nodeType string) (PropertyType, bool) // promise

	UpsertChild(nodeType, nodeKeyName string, nodeProp PropertyType, linkType string,
		linkProp PropertyType) BaseNodeInterface // promise
	UpsertLink(destNode BaseNodeInterface, linkType string, linkProp PropertyType,
		isSingleton bool, rootName string, useUUIDasKey bool) *GLink // promise
	DeleteLink(childKey, nodeType string) bool                         // promise
	DeleteLinkToNode(destNode BaseNodeInterface, linkType string) bool // promise

	Delete() // promise

	IsAnyParentDeleted(level ...uint32) bool
	GetFullPath() NodePathList
	// get/set creator ID
	GetCreatorId() string            // : Promise<string>;
	SetCreatorId(seed string) string // Promise<string>;
	// get the stats for data base access and message  bus access
	GetDBStatsString() string
	GetMSGStatsString() string
	// allows a process to indicate that it will own a certain part of the graph
	// as a result items created by other processes will be automatically removed.
	ClaimOwnership(path NodePathList, timeDelay uint32, prefix string)
	ClaimLinkOwnership(path NodePathList, destNodeType string, timeDelay uint32)
	// call back for node changes
	RegisterCB(pattern NodePathList, cbfn CallbackFuncNode)
	// call back for link changes
	RegisterLinkCB(pattern NodePathList, destNodeType string, cbfn CallbackFuncLink)
	Subscribe(pattern NodePathList, depth uint32)
	Unsubscribe(pattern NodePathList)
	// getRefNode()
	GetTypeRefChildNode(ntype, nkeyname, nkeyvalue string) BaseNodeInterface
	IsRLinkFeatureFlagEnabled() bool
	Checksub()
	// rlink iterator handler - cache
	GetCachedRLinkHdlr(nodeType, childKey string) BaseNodeLinkManagerIteratorInterface
	ForEachCachedRLinkWithNodeType(fn func(rLinkHdlr BaseNodeLinkManagerIteratorInterface) interface{},
		nodeType string)
	ForEachCachedRLink(fn func(rLinkHdlr BaseNodeLinkManagerIteratorInterface) interface{})
	// hlink iterator handler - cache
	GetCachedChildLinkHdlr(nodeType, childKey string) BaseNodeLinkManagerIteratorInterface
	ForEachCachedChildLinkWithNodeType(fn func(rLinkHdlr BaseNodeLinkManagerIteratorInterface) interface{},
		nodeType string)
	ForEachCachedChildLink(fn func(rLinkHdlr BaseNodeLinkManagerIteratorInterface) interface{})
	GetLatestRevision() int64
	CompactRevision(int64, int64, bool) (*EtcdCompactResponse, error)
}
