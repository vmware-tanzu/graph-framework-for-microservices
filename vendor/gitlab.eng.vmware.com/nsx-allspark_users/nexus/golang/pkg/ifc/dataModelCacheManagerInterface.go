package ifc

type DataModelCacheManagerInterface interface {
	// setLog(log)
	IsNodeInCache(n string) bool
	IsNodePresentAndValid(n string) bool
	IsNodeSubscribed(n string) bool
	IsNodeDeleted(n string) bool
	GetCachedNode(n string) (BaseNodeInterface, bool)

	// add anode to cache and perform the callbacks associated with this update
	WriteToCache(parentId, id string, base BaseNodeInterface)
	// missing function
	// update node
	// updateCachedNodeProperty(n: string, prop:  { [key: string]: any });
	// update link property

	// purge a node and all it's children that are not soft links
	// all the callbacks for this purge will be done as well.
	PurgeCacheTree(nodeId string)

	// delete subscription for a give path
	DelSubscription(path NodePathList)

	AddSubscription(
		path NodePathList,
		depth uint32) []string
	// add subscription for a path, with callback function
	AddNodeSubscription(
		path NodePathList,
		cbfn *CallbackFuncNode) []string
	// add subscription for a path, with callback function
	AddLinkSubscription(
		path NodePathList,
		cbfn *CallbackFuncLink,
		destNodeType string) []string

	// check if the given path is under subscription update
	CheckPath(path NodePathList) bool

	//    doNodeAddCallbacks(base: BaseNodeInterface, path: NodePathList);
	//    doNodeDeleteCallbacks(base: BaseNodeInterface, path: NodePathList);
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

	// need to update base node so it can do the callback's directly on update of property or delete of node
	// also link handleing should be done with the base node for link updates
	// the call back function will stay in the link cache manager
	// base node can call cache manager functions

	DumpCachedNodes() map[string]interface{}
}
