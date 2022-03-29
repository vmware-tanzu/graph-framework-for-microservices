package ifc

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/common"
)

type DataModelInterface interface {
	GetName() string
	Init()
	IsRLinkFeatureFlagEnabled() bool
	GetCCnt() uint64
	SetGraph(graph GraphDBInterface)
	IsRunning() bool
	GetId() string
	SetId(id string)
	Shutdown()
	DeleteDMFromDebugCache()
	SetDMDebugCache()
	Sync(node BaseNodeInterface, forceSync bool)
	Delete(node BaseNodeInterface)
	DeleteId(nodeId string, fDelete bool)
	GetRootNode(nodeId string, dbsync bool) (BaseNodeInterface, bool)
	GetNode(paretnId string, nodeType, nodeId string, dbSync bool, syncSubTrees bool) (BaseNodeInterface, bool)
	GetCachedNode(nodeId string) (BaseNodeInterface, bool)
	GetNodeFromCache(id string) (BaseNodeInterface, bool)
	IsNodeInCache(n string) bool
	SetMessagingDelay(n uint32)
	// todo: Change to array of base from root.....
	PopulatePathAndFetchNodes(path NodePathList, forceSync bool) []BaseNodeInterface
	UpsertNode(
		parentId,
		linkType,
		linkKey string,
		linkProp PropertyType,
		nodeType string,
		nodeProp PropertyType) (BaseNodeInterface, *GLink)
	UpsertRootNode(nodeType, nodeName string, nodeProp PropertyType) BaseNodeInterface
	UpdateNodeAddProperties(nodeId string, nodeProp PropertyType)
	UpdateNodeRemoveProperties(nodeId string, nodeProp []string)
	UpsertLink(lnkObj common.UpsertLinkOpObj) *GLink

	/*
		UpsertLink(linkType, srcNodeId, dstNodeId, srcNodeType, dstNodeType, srcKeyValue, dstKeyValue string,
			srcNodeProp, dstNodeProp, linkProp PropertyType,
			isSingleton bool) *GLink
	*/
	DeleteLink(srcNodeId, dstNodeType, dstKeyValue string)
	DeleteRLink(srcNodeId, dstNodeType, dstKeyValue string)
	Subscribe(pattern NodePathList, depth uint32)
	Unsubscribe(pattern NodePathList)

	GetDBStats() common.GraphDBStats
	GetMsgStats() NotificationStats
	RegisterCB(pattern NodePathList, cbfn CallbackFuncNode)
	RegisterLinkCB(pattern NodePathList, destNodeType string, cbfn CallbackFuncLink)
	ClaimOwnership(srcPattern NodePathList, timeDelay uint32, prefix string)
	ClaimLinkOwnership(srcPattern NodePathList, destNode string, timeDelay uint32)

	UpdateLinkAddProperty(
		srcNodeId,
		dstNodeType,
		dstKeyValue string,
		linkProp PropertyType) bool
	UpdateLinkRemoveProperty(
		srcNodeId,
		dstNodeType,
		dstKeyValue string,
		linkPropKeys []string) bool
	GetNextChildKeyList(
		parentId, nodeType, childStartingKey string, batchSize uint32) []string
	GetNextRLinkChildKeyList(
		parentId, nodeType, childStartingKey string, batchSize uint32) []string
	SyncNodeProperty(node BaseNodeInterface, forced bool)
	SyncChild(pnode BaseNodeInterface, childNodeType, childNodeKey string,
		forced bool)
	GetGraph() GraphDBInterface
	DoNodeUpdateCallbacks(
		base BaseNodeInterface,
		curProp PropertyType,
		newProp PropertyType)
	DoLinkDeleteCallbacks(
		base BaseNodeInterface,
		destType string,
		dstKey string,
		oldLink PropertyType)
	DoLinkUpdateCallbacks(
		base BaseNodeInterface,
		dstKey string,
		destType string,
		curLink *GLink,
		newLink *GLink)
	GetLatestRevision() int64
	CompactRevision(int64, int64, bool) (*EtcdCompactResponse, error)
	CollectDMCachedData() map[string]interface{}
}
