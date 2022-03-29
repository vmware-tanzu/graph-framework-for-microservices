package ifc

type DataModelCacheInterface interface {
	Init(dmbus DMBusInterface)
	// sync a basenode to match data in graphdb and update the cache as well
	Sync(node BaseNodeInterface, forced bool, syncSubTrees bool)

	// if (sync is false) or (node.subscribe is true) returns node from the cache,
	//       if not in cache it will be synced to db and then the node is returned.
	// if sync is true and subscribed is false it will force a sync to db even if the
	//      nodeid is in cache.
	GetNode(parentId, nodeId, nodeType string, sync bool, syncSubTrees bool) (BaseNodeInterface, bool)

	// same as GetNode, but a copy of this node must be in cache before calling this function
	GetCachedNode(nodeId string, sync bool) (BaseNodeInterface, bool)
	GetCachedNodeSync(nodeId string) (BaseNodeInterface, bool)

	// give a graph node (From graphdb library), create a base node and return it.
	// also populate the base node in cache.
	// this is needed when a new node is created in the graph.
	PopulateCacheAndReturnBaseNode(parentId string, n *GNode, keyName string) BaseNodeInterface

	// remove node and all it's child nodes from the cache.
	PurgeCacheTree(nodeId string)

	SyncNodeProperty(node BaseNodeInterface, forced bool)
	SyncChild(
		pnode BaseNodeInterface, childNodeType string,
		childNodeKey string, forced bool)

	// addLink(node: BaseNodeInterface, dstNodeType: string, dstKeyValue: string, link: GLink);
	// notificaiton for an event this will cause subscribed nodes to be updated....

	// linkUpdateNotification used in case of parent to child hard link update
	// linkUpdateNotification(
	//     pnode: BaseNodeInterface,
	//     destType: string,
	//     destKeyValue: string,
	//     oldLink: GLink | undefined,
	//     newLink: GLink | undefined
	// );
	// Also have to decide on callback's.
	// it's natural to have the callback's here as well since the tree is populated here

	// how about regex type of subscriptions ....
	// how to handle those.
	// may be deffer regex to another commit.
	AddSubscription(path NodePathList, depth uint32)
	DelSubscription(path NodePathList)
	// RegisterCB(path NodePathList, cbfn func(interface{}), isLinkCB bool, dstNodeType string)
	RegisterNodeCB(path NodePathList, cbfn *CallbackFuncNode)
	RegisterLinkCB(path NodePathList, cbfn *CallbackFuncLink, destNodeType string)

	SetName(n string)
	// MsgCB(msg *Notification)
	IsNodeInCache(n string) bool
	Stop()
	// Init(name string)
	// init(name string, mbus *DMBusInterface)
	SetMessagingDelay(n uint32)

	ReadyDestNodeBeforeCB(
		base BaseNodeInterface,
		destKey, destType string,
		linkProperty PropertyType)

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

	DumpCachedNodes() map[string]interface{}
}
