package cache

import (
	"encoding/json"
	"strings"
	"time"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/internal/cache/manager"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/internal/utils"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/base"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/ifc"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/logging"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/scheduler"
)

type DataModelCache struct {
	dm           ifc.DataModelInterface
	cacheManager ifc.DataModelCacheManagerInterface
	sch          *scheduler.Scheduler
	running      bool
	msgDelay     uint32
	name         string
	featureFlags []string
}

func NewDataModelCache(name string, dm ifc.DataModelInterface, dmFeatureFlags []string) *DataModelCache {
	e := &DataModelCache{
		dm:           dm,
		name:         name,
		sch:          scheduler.NewScheduler(),
		running:      false,
		featureFlags: dmFeatureFlags,
		msgDelay:     0}
	e.cacheManager = manager.NewCacheManager(name, e)
	return e
}
func (c *DataModelCache) SetMessagingDelay(n uint32) {
	c.msgDelay = n
}
func (c *DataModelCache) SetName(n string) {

}
func (c *DataModelCache) isRLinkFeatureFlagEnabled() bool {
	for _, v := range c.featureFlags {
		if strings.Compare(v, common.Rlink_EnableFeatureFlag) == 0 {
			return true
		}
	}
	return false
}
func (c *DataModelCache) Init(dmbus ifc.DMBusInterface) {
	dmbus.RegisterCB(func(n *ifc.Notification) {
		c.msgCBPre(n)
	})
	c.running = true
}
func (c *DataModelCache) Stop() {
	c.running = false
}
func (c *DataModelCache) msgCBPre(msg *ifc.Notification) {
	if c.running {
		if c.msgDelay != 0 {
			// this is used for testing
			time.Sleep(time.Duration(c.msgDelay) * time.Millisecond)
		}
		c.MsgCB(msg)
	}
}
func (c *DataModelCache) IsNodeInCache(n string) bool {
	return c.cacheManager.IsNodeInCache(n)
}
func (c *DataModelCache) PurgeCacheTree(n string) {
	c.cacheManager.PurgeCacheTree(n)
}
func (c *DataModelCache) AddSubscription(path ifc.NodePathList, depth uint32) {
	c.addSubscriptionInternal(path, depth, nil, nil, false, "")
}
func (c *DataModelCache) DelSubscription(path ifc.NodePathList) {
	c.cacheManager.DelSubscription(path)
}
func (c *DataModelCache) RegisterNodeCB(path ifc.NodePathList, cbfn *ifc.CallbackFuncNode) {
	c.addSubscriptionInternal(path, 0, cbfn, nil, false, "")
}
func (c *DataModelCache) RegisterLinkCB(path ifc.NodePathList, cbfn *ifc.CallbackFuncLink, destNodeType string) {
	c.addSubscriptionInternal(path, 0, nil, cbfn, true, destNodeType)
}
func (c *DataModelCache) Sync(node ifc.BaseNodeInterface, forced bool, syncSubTrees bool) {
	c.syncInternal(node, forced, syncSubTrees)
}
func (c *DataModelCache) GetNode(parentId, nodeId, nodeType string,
	sync bool, syncSubTrees bool) (ifc.BaseNodeInterface, bool) {
	logging.Debugf("GetNode(): parentId: %s, nodeId: %s\n",
		parentId, nodeId)
	if valid := c.fetch(parentId, nodeId, sync, syncSubTrees); !valid {
		if !valid && nodeId == "root" {
			valid = c.fetch(parentId, "/Root/"+common.NodeFixedProp_NodeSingletonKeyValue, sync,
				syncSubTrees)
			if valid {
				return c.cacheManager.GetCachedNode("/Root/" + common.NodeFixedProp_NodeSingletonKeyValue)
			}
		}
		if !valid {
			return nil, false
		}
	}
	return c.cacheManager.GetCachedNode(nodeId)
}

func (c *DataModelCache) GetCachedNode(nodeId string, sync bool) (ifc.BaseNodeInterface, bool) {
	logging.Debugf("%s id %s %t %t", c.name, nodeId, c.cacheManager.IsNodePresentAndValid(nodeId), sync)
	if c.cacheManager.IsNodePresentAndValid(nodeId) {
		if sync {
			if gcn, gcnOk := c.cacheManager.GetCachedNode(nodeId); gcnOk {
				parent := gcn.GetBaseParent()
				parentId := ""
				if parent != nil {
					parentId = parent.GetId()
				}
				hasParentNode := c.fetch(parentId, nodeId, sync, false)
				if !hasParentNode {
					return nil, false
				}
			}
		}
		return c.cacheManager.GetCachedNode(nodeId)
	} else if nodeId == "root" {
		return c.GetCachedNode("/Root/"+common.NodeFixedProp_NodeSingletonKeyValue, sync)
	}
	return nil, false
}

func (c *DataModelCache) GetCachedNodeSync(nodeId string) (ifc.BaseNodeInterface, bool) {
	if c.cacheManager.IsNodePresentAndValid(nodeId) {
		return c.cacheManager.GetCachedNode(nodeId)
	} else if nodeId == "root" {
		return c.GetCachedNodeSync("/Root/" + common.NodeFixedProp_NodeSingletonKeyValue)
	}
	return nil, false
}

/*
   Link updates will have the linkid and the node path for the destination node and the nodeid for the parent node.
   when link  updates are received, check if the parent node is in cache, (Sync the parent node if present)
   if so then check if the destination node path is satisfying the subscription
   if so then fetch the destinaiton node and add to cache.
   Go thorugh each link of the destination node and check if they need to be fetched.
   DO this recursively ...

   whenever syncing ...
   any addition of link/node will trigger sub check and further loading from graph.

*/
// handle node update notification
func (c *DataModelCache) msgCBWithValueHandleNode(msg *ifc.Notification) {
	sid := msg.UpdatedObjId
	objPath := msg.ObjectPath
	updateType := common.UpdateTypeConvFromString(msg.UpdateType)
	rev := msg.Revision
	// serialize on sid
	logging.Debugf("%s:: LOCK key=%s", c.name, sid)
	lck := c.sch.Wait(sid)
	defer c.sch.Done(lck)
	defer logging.Debugf("%s:: UnLOCK key=%s", c.name, sid)

	if c.cacheManager.CheckPath(objPath) {
		nodeId := msg.UpdatedObjId
		// if this node is already present update it.
		// if this node is not present and it's parent node is present add the node to cache.
		if n, nok := c.cacheManager.GetCachedNode(nodeId); nok {
			logging.Debugf("%s:got node update for existing node in cache %s will sync", c.name, nodeId)
			eRevProp := n.GetImmediateBaseProperties()[common.NodeFixedProp_Revision]
			eRev := int64(0)
			if eRevProp != nil {
				eRev = eRevProp.(int64)
			}
			if updateType != common.UpdateType_NodeDelete {
				if rev < eRev {
					//this.log.debug(` Ignore node update ${nodeId} :: received ${rev} and local ${eRev}`);
					return
				}
				newProp := msg.Value
				if _, ok := newProp[common.NodeFixedProp_createdBy]; !ok {
					newProp[common.NodeFixedProp_createdBy] = n.GetImmediateBaseProperties()[common.NodeFixedProp_createdBy]
					newProp[common.NodeFixedProp_creationTime] =
						n.GetImmediateBaseProperties()[common.NodeFixedProp_creationTime]
				}
				n.SetImmediateBaseProperties(newProp)
			} else {
				if rev < eRev {
					// this.log.debug(` Ignore node Delete ${nodeId} :: received ${rev} and local ${eRev}`);
					return
				}
				c.cacheManager.PurgeCacheTree(nodeId)
			}
		} else {
			logging.Debugf("%s:got node update for non-existant node id:%s, will skip and wait for an link update", c.name, nodeId)
			// nop .. will wait for the link update to pull this node.
		}
	}
}

// handle notification for link update when it's an existing link in database
func (c *DataModelCache) msgCBWithValueExistingLink(msg *ifc.Notification, parentNode ifc.BaseNodeInterface, del bool) string {
	dstNodeType := msg.UpdatedObjType
	dstNodeKey := msg.UpdatedObjKey
	rev := msg.Revision
	dnodeId := ""
	if del {
		if parentNodeLinkTo, ok := parentNode.GetLinks().Get(dstNodeType, dstNodeKey); ok {
			lrev := parentNodeLinkTo.Properties[common.LinkFixedProp_Revision]
			if lrev != nil {
				if rev != 0 && rev < lrev.(int64) {
					logging.Debugf("%s:Ignore Link Delete Update %d < %d", c.name, rev, lrev)
					return dnodeId
				}
			}
			parentNode.DeleteImmediateLink(dstNodeType, dstNodeKey)
		}
	} else {
		if parentNodeLinkTo, ok := parentNode.GetLinks().Get(dstNodeType, dstNodeKey); ok {
			lrev := parentNodeLinkTo.Properties[common.LinkFixedProp_Revision]
			if lrev != nil {
				if rev != 0 && rev < lrev.(int64) {
					logging.Debugf("%s: Ignore Link Update %d < %d", c.name, rev, lrev)
					return dnodeId
				}
			}
			nlink := &ifc.GLink{
				Id:                parentNodeLinkTo.Id,
				LinkType:          parentNodeLinkTo.LinkType,
				Properties:        msg.Value,
				SourceNodeId:      parentNodeLinkTo.SourceNodeId,
				DestinationNodeId: parentNodeLinkTo.DestinationNodeId,
			}
			if _, ok := nlink.Properties[common.LinkFixedProp_createdBy]; !ok {
				nlink.Properties[common.LinkFixedProp_createdBy] = parentNodeLinkTo.Properties[common.LinkFixedProp_createdBy]
				nlink.Properties[common.LinkFixedProp_creationTime] = parentNodeLinkTo.Properties[common.LinkFixedProp_creationTime]
			}
			parentNode.UpsertImmediateLink(dstNodeType, dstNodeKey, nlink)
			olinkProp := parentNodeLinkTo.Properties
			dnodeId = (olinkProp[common.LinkFixedProp_destNodeId]).(string)
		}
	}
	return dnodeId
}

// handle notification for link update when it's a new link and not in cache
func (c *DataModelCache) msgCBWithValueNewLink(msg *ifc.Notification, parentNode ifc.BaseNodeInterface) string {
	// it's anew link
	dnodeId := ""
	linkId := msg.UpdatedObjId
	dstNodeType := msg.UpdatedObjType
	dstNodeKey := msg.UpdatedObjKey
	parentId := msg.UpdatedObjParentId
	// just get it from the db
	// while we are waiting there may be a  sub added
	// which will do a link add call back,
	// after the wait this function should do an link update instead of add.
	if _, ok := msg.Value[common.LinkFixedProp_createdBy]; ok {
		// if created by is not present in the value it will
		// need to be fetched from the database.
		srcSplit := strings.Split(linkId, "/")
		srcId := strings.Join(srcSplit[:len(srcSplit)-3], "/")
		lnk := &ifc.GLink{
			Id:                linkId,
			LinkType:          msg.Value[common.LinkFixedProp_NodeType].(string), // ["type"].(string),
			Properties:        msg.Value,
			DestinationNodeId: msg.Value[common.LinkFixedProp_destNodeId].(string),
			SourceNodeId:      srcId,
		}
		lnk.Properties[common.LinkFixedProp_Revision] = msg.Revision
		parentNode.UpsertImmediateLink(dstNodeType, dstNodeKey, lnk)
		dnodeId = lnk.Properties[common.LinkFixedProp_destNodeId].(string)
	} else {
		var lnk *ifc.GLink = &ifc.GLink{
			Id:                linkId,
			LinkType:          "??",
			Properties:        nil,
			SourceNodeId:      parentId,
			DestinationNodeId: ""}
		lnk.Properties = msg.Value
		lnk.DestinationNodeId = (lnk.Properties[common.LinkFixedProp_destNodeId]).(string)
		parentNode.UpsertImmediateLink(dstNodeType, dstNodeKey, lnk)
		dnodeId = (lnk.Properties[common.LinkFixedProp_destNodeId]).(string)
	}
	return dnodeId
}

// handle notification that provide the value with the change on key
func (c *DataModelCache) msgCBWithValue(msg *ifc.Notification) {
	// the update for a node will need to be serialized with the sync notificaiton for the same
	// if (!msg.value) throw Error('no value in msg notification');
	updateType := common.UpdateTypeConvFromString(msg.UpdateType)
	nupdate := (updateType == common.UpdateType_NodeAdd ||
		updateType == common.UpdateType_NodeDelete ||
		updateType == common.UpdateType_NodeUpdate)
	if nupdate {
		/*
			Node deletion op generates an update event when ToBeDeleted is set to true in the node properties.
			To stub out an event for a node update in case of a deletion op, need changes here to check for FixedNodeProp_ToBeDeleted in msg.Value and comparewith the one in existing set of properties of the node.
			Would not recommend this since we'd want the deleted prop to be reflected in basenode and cache node objects.
		*/
		c.msgCBWithValueHandleNode(msg)
		return
	}
	sid := msg.UpdatedObjParentId
	// serialize on sid
	logging.Debugf("%s:: got link update Before Lock path = %s parent lock=%s", c.name, msg.UpdatedObjId, sid)
	logging.Debugf("%s:: LOCK key=%s", c.name, sid)
	lck := c.sch.Wait(sid)
	defer c.sch.Done(lck)
	defer logging.Debugf("%s:: UnLOCK key=%s", c.name, sid)

	logging.Debugf("%s:: got link update After Lock path = %s parent lock=%s", c.name, msg.UpdatedObjId, sid)
	parentId := msg.UpdatedObjParentId
	dstNodeType := msg.UpdatedObjType
	dstNodeKey := msg.UpdatedObjKey
	linkId := msg.UpdatedObjId
	if pn, pnOk := c.cacheManager.GetCachedNode(parentId); pnOk {
		dnodeId := ""
		lnkOk := pn.GetLinks().Has(dstNodeType, dstNodeKey)
		logging.Debugf("%s::   Linkupdate found parent in cache will update. parentId=%s link present=%t", c.name, parentId, lnkOk)

		if lnkOk {
			dnodeId = c.msgCBWithValueExistingLink(msg, pn, updateType == common.UpdateType_LinkDelete)
		} else if updateType != common.UpdateType_LinkDelete {
			// it's anew link
			dnodeId = c.msgCBWithValueNewLink(msg, pn)
		} else {
			logging.Debugf("%s:: Ignore update as this link %s is not in cache %d", c.name, linkId, updateType)
		}
		// do a sync on the destination node
		if dnodeId != "" && !c.cacheManager.IsNodeInCache(dnodeId) {
			c.GetNode(parentId, dnodeId, dstNodeType, false, false)
		}
		// just fetch parent and it will update all the child links
	} else {
		logging.Debugf("%s::  This parent is not in cache, will skip. parentid=%s", c.name, parentId)
	}
}
func (c *DataModelCache) MsgCB(msg *ifc.Notification) {
	//	if c.dm.GetId() != msg.Source {
	// logging.Debugf("MsgCB dmId = %s msg send Id = %s\n", c.dm.GetId(), msg.Source)
	c.msgCBWithValue(msg)
	//	}
}

func (c *DataModelCache) PopulateCacheAndReturnBaseNode(parentId string, n *ifc.GNode, keyName string) ifc.BaseNodeInterface {
	nodeId := n.Id
	if keyName == "" {
		keyName = common.NodeFixedProp_NodeDefaultKeyName // default keyname is 'name'
	}
	node, nodeOk := c.cacheManager.GetCachedNode(nodeId)
	// check if a prior node existed and do the callback's
	if !nodeOk {
		node = base.NewBaseNode(c.dm, n.Type, keyName, n.Id, n.Properties, false)
		c.cacheManager.WriteToCache(parentId, nodeId, node)
		logging.Debugf("%s:Node %s was not in cache Adding a new Node. %v",
			c.name, nodeId, n.Properties)
	} else {
		// Q: when base node is already there, is there no need to update the data in base node
		// from the graph node data in argument "n".
		// A: populateCacheLinks is misleading name, it updates the current base node with data from "n"
		logging.Debugf("%s node %s from cache %v", c.name, nodeId, n.Properties)
	}
	c.populateCacheLinks(node, n, false)
	return node
}

func (c *DataModelCache) addSubscriptionInternal(
	path ifc.NodePathList,
	depth uint32,
	cbfnNode *ifc.CallbackFuncNode,
	cbfnLink *ifc.CallbackFuncLink,
	isLinkCB bool,
	destNodeType string) {
	var syncList []string
	if isLinkCB && cbfnLink != nil {
		syncList = c.cacheManager.AddLinkSubscription(path, cbfnLink, destNodeType)
	} else if cbfnNode != nil {
		syncList = c.cacheManager.AddNodeSubscription(path, cbfnNode)
	} else {
		syncList = c.cacheManager.AddSubscription(path, depth)
	}
	// all the nodes in syncList need to be synced to database.
	for _, itm := range syncList {
		if nd, ndOk := c.cacheManager.GetCachedNode(itm); ndOk {
			go c.syncInternal(nd, true, false) // Maybe needs serialization? Not sure...
		} else {
			logging.Errorf(c.name + ": Sync list is with a node id ${itm} which is not in cache.")
		}
	}
	if depth > 1 {
		go c.forceRefreshNode(path, depth, true, nil) // Maybe needs serialization? Not sure...
	}
}

func (c *DataModelCache) SyncNodeProperty(node ifc.BaseNodeInterface, forced bool) {

	if !forced {
		if c.cacheManager.IsNodeSubscribed(node.GetId()) {
			return
		}
	}

	if prop := c.dm.GetGraph().GetNodeProperty(node.GetId()); prop != nil {
		dt := make(ifc.PropertyType)
		for k, v := range prop {
			dt[k] = v
		}
		// logging.Debugf("SyncNodeProperty: Nodeid %s prop %s", node.GetId(), dt)
		node.SetImmediateBaseProperties(dt)
	} else {
		logging.Debugf("%s returned no prop, marking %s as deleted.", c.name, node.GetId())
		c.cacheManager.PurgeCacheTree(node.GetId())
	}
}

func (c *DataModelCache) SyncChild(pnode ifc.BaseNodeInterface, childNodeType, childNodeKey string, forced bool) {
	// check if child node exists in cache.
	// debug.PrintStack()
	if cnl, cnlOk := pnode.GetLinks().Get(childNodeType, childNodeKey); cnlOk {
		cn, cnOk := c.GetCachedNode(cnl.DestinationNodeId, false)
		logging.Debugf("%s called for %s->%s/%s forced=%t",
			c.name, pnode.GetId(), childNodeType, childNodeKey, forced)
		if cnOk && cnl.Id != "" {
			if forced || (!c.cacheManager.IsNodeSubscribed(cn.GetId())) {
				lnk, lnkOk := c.dm.GetGraph().DescribeLinkById(cnl.Id)
				logging.Debugf("%s %s->%s/%s lnk = %s", c.name, pnode.GetId(), childNodeType, childNodeKey, lnk)
				if lnkOk {
					// likely link is deleted.
					// fall back to full sync
					// sync link property
					c.populateSingleCacheLink(pnode, lnk, false)
					// sync child node
					c.SyncNodeProperty(cn, forced)
					return
				}
			} else {
				logging.Debugf("%s called for %s->%s/%s forced=%t skipping sync as we are already subscribed",
					c.name, pnode.GetId(), childNodeType, childNodeKey, forced)
				return
			}
		}
	}
	logging.Debugf("%s called for %s->%s/%s forced=%t Full Sync ", c.name, pnode.GetId(), childNodeType, childNodeKey, forced)
	// node/link is not populdated in cache do a full sync
	// TODO this can be still focused.
	c.Sync(pnode, forced, false)
}

func (c *DataModelCache) syncInternal(node ifc.BaseNodeInterface, forced bool, syncSubTrees bool) {
	logging.Debugf("%s:: LOCK key=%s", c.name, node.GetId())
	lck := c.sch.Wait(node.GetId())
	defer c.sch.Done(lck)
	defer logging.Debugf("%s:: UnLOCK key=%s", c.name, node.GetId())
	if !c.dm.IsRunning() {
		return
	}
	if node.GetId() == "" {
		logging.Errorf(c.name + ":Node ID is not defined for the node being synced on. (Likely this base node is virtual and for subscription)")
	}
	logging.Debugf("syncInternal(): %s: Syncing node with id=%s from Graph forced=%t", c.name, node.GetId(), forced)
	// no sync if it's already subscribed
	if !forced {
		if node.IsDeleteCompleted() {
			return
		}
		if c.cacheManager.IsNodeSubscribed(node.GetId()) &&
			!syncSubTrees {
			logging.Debugf("syncInternal(): Node %s: Skip sync due to sub",
				node.GetId())
			return
		}
	}
	n := c.dm.GetGraph().DescribeNode(node.GetId(), node.GetType())
	// TODO: Add a diff for all the links ... and the destination node
	// all the new nodes if they are subscribed then do a sync on them as well.
	// if diff and if any callback is registered for this node path trigger it as well.

	// what if this node was deleted
	if n != nil {
		logging.Debugf("syncInternal(): %s: Syncing node with id=%s got graph response.", c.name, node.GetId())
		c.populateCacheLinks(node, n, syncSubTrees)
	} else {
		// this node is deleted, let's remove it and all it's children
		logging.Debugf("syncInternal(): %s: Setting  node id=%s to be deleted ", c.name, node.GetId())
		c.cacheManager.PurgeCacheTree(node.GetId())
	}
}

// fetch a node into cache if its successful will return true
func (c *DataModelCache) fetch(parentId string, nodeId string, forcedRead bool,
	syncSubTrees bool) bool {
	if nodeId == "" {
		logging.Fatalf("%s:fetch called with nodeid Set to %s.", c.name, nodeId)
	}
	/*
	   if (this.cacheManager.isNodeInCache(parentId))
	       throw new Error(
	           `Fetch called with parent id ${parentId} which is not in cache`
	       );*/

	if nd, ndOk := c.cacheManager.GetCachedNode(nodeId); ndOk {
		if !forcedRead {
			logging.Debugf("fetch(): parentId: %s, nodeId: %s"+
				" Node %s found in cache.Skip sync",
				parentId, nodeId, nodeId)
			return true
		}
		logging.Debugf("fetch(): parentId: %s, nodeId:%s"+
			" Node %s not found in cache. Sync",
			parentId, nodeId, nodeId)
		c.Sync(nd, false, syncSubTrees)
		return !nd.IsDeleted()
	}

	logging.Debugf("%s:: LOCK key=%s", c.name, nodeId)
	lck := c.sch.Wait(nodeId)
	defer c.sch.Done(lck)
	defer logging.Debugf("%s:: UnLOCK key=%s", c.name, nodeId)
	nodeType := ""
	if nd, ndOk := c.cacheManager.GetCachedNode(nodeId); ndOk {
		nodeType = nd.GetType()
	}

	// if a copy of the node is not in the local data cache
	// fetch it even if it's subscribed node.
	n := c.dm.GetGraph().DescribeNode(nodeId, nodeType)
	// what if this node was deleted
	if n != nil {
		keyName := common.NodeFixedProp_NodeDefaultKeyName
		if k, kok := n.Properties[common.NodeFixedProp_NodeKeyName]; kok {
			keyName = k.(string)
		}
		// double check, the node may be in cache now
		if nd, ndOk := c.cacheManager.GetCachedNode(nodeId); !ndOk {
			node := base.NewBaseNode(
				c.dm, n.Type, keyName, n.Id, n.Properties, false)
			// node.properties = n.properties;
			// node.id = n.id;
			c.cacheManager.WriteToCache(parentId, nodeId, node)
			c.populateCacheLinks(node, n, syncSubTrees)
			return true
		} else {
			// if the node is in cache just update it.
			if !nd.IsDeleted() {
				c.populateCacheLinks(nd, n, syncSubTrees)
				return true
			}
		}
	}
	return false
}

//
func (c *DataModelCache) populateSingleCacheLink(node ifc.BaseNodeInterface,
	newLink *ifc.GLink, syncSubTrees bool) {
	basePath := node.GetFullPath()
	destType := newLink.Properties[common.LinkFixedProp_NodeType].(string)
	destKeyValue := newLink.Properties[common.LinkFixedProp_NodeKeyValue].(string)
	destPath := append(basePath, ifc.NodePath{destType, destKeyValue})
	softLink := (newLink.Properties[common.LinkFixedProp_HardLink]).(string) != "true"
	sdlp := newLink.Properties[common.LinkFixedProp_SoftLinkDestinationPath]
	logging.Debugf("populateSingleCacheLink(): %s: destType=%s destPath=%s softLink=%t sdlp=%s", c.name, destType, destPath, softLink, sdlp)

	node.UpsertImmediateLink(destType, destKeyValue, newLink)
	logging.Debugf("populateSingleCacheLink(): %s: populateCacheLinks: 1 checkpath(%v) = %t", c.name, destPath, c.cacheManager.CheckPath(destPath))
	if softLink {
		// fetch destination obj
		var destNodePath ifc.NodePathList
		e := json.Unmarshal([]byte(sdlp.(string)), &destNodePath)
		if e != nil {
			logging.Errorf("populateSingleCacheLink(): %s:while parsing sdlp link property %s %e", c.name, sdlp, e)
		}
		// add reverse link
		dnodes := c.dm.PopulatePathAndFetchNodes(destNodePath, false)
		if len(dnodes) != 0 {
			logging.Debugf("Paths identified for node %s and destination node path %s: \n %v", node.GetId(),
				destNodePath, dnodes)
		}
		for _, dnode := range dnodes {
			/*
					Looks like all the dest paths are found and rlinks are created for each
					We only want to create rlinks for the relevant src and dst node i.e. newlink.id = node.GetId()+destType+destKeyValue
					Include logic here to parse destType and destKeyVal from dnode.id. For ex: rSoftLinks are populated for
				    	/Root/root/ -> /Root/root/,  destKeyValue: User0, linkId:/Root/root/Config/default/Service/Svc0/_links/User/User0
					and similarly for
					   /Root/root/<subtreeroot>/keyvalue
					  /Root/root/<subtreeroot>/keyvalue/subtreeroot2/keyvalueroot2
			*/
			dnodeIdSegments := strings.Split(dnode.GetId(), "/")
			slen := len(dnodeIdSegments)
			if slen >= 2 {
				tNodeName := dnodeIdSegments[slen-1]
				tNodeType := dnodeIdSegments[slen-2]
				if tNodeName == destKeyValue &&
					tNodeType == destType {
					dnode.AddReverseLink(node.GetId(), destKeyValue, newLink.Id)

				}
			}
		}
	} else if c.cacheManager.CheckPath(destPath) {

		// TODO: this whole else clause may be able to be removed.
		logging.Debugf("populateSingleCacheLink(): %s:Sub Fetch:Adding1 %s to %v", c.name, newLink.DestinationNodeId, destPath)
		go c.GetNode(node.GetId(), newLink.DestinationNodeId, destType, false, false) // TODO: should be able to get rid of goroutine here.
	} else if syncSubTrees {
		// Traverse and sync subTrees
		logging.Debugf("populateSingleCacheLink(): %s:"+
			" Invoking recursive subTree traversal"+
			" for parent node:%s, destNodeId:%s\n",
			c.name, node.GetId(), newLink.DestinationNodeId)
		c.GetNode(node.GetId(), newLink.DestinationNodeId, destType, true, syncSubTrees)
	}
}

func (c *DataModelCache) populateSingleCacheRLink(node ifc.BaseNodeInterface,
	newLink *ifc.GLink, syncSubTrees bool) {
	destType := newLink.Properties[common.LinkFixedProp_NodeType].(string)
	destKeyValue := newLink.Properties[common.LinkFixedProp_NodeKeyValue].(string)
	logging.Debugf("DM cache: %s: populateSingleCacheRLink(): destType=%s destNodeKey=%s"+
		" Upsert immediate rlink called for %s", c.name, destType, destKeyValue,
		newLink.Id)
	node.UpsertImmediateRLink(destType, destKeyValue, newLink)
}

// Mark and sweep links in node cache
func (c *DataModelCache) markAndSweepNodeCacheLinks(node ifc.BaseNodeInterface, nodeData *ifc.GNode, rLink bool) {
	// remove the links that do not exist in the data returned from the database.
	linkListInDB := make(map[string][]string)
	var lnkObj []*ifc.GLink
	if rLink {
		lnkObj = nodeData.RLinks
	} else {
		lnkObj = nodeData.Links
	}
	for _, itm := range lnkObj {
		// all the links that are returned from db.
		destKeyValueT, nkvOk := itm.Properties[common.LinkFixedProp_NodeKeyValue]
		destTypeT, ntOk := itm.Properties[common.LinkFixedProp_NodeType]
		destKeyValue := destKeyValueT.(string)
		destType := destTypeT.(string)
		if nkvOk && ntOk {
			// if link has the proper properties. get the destination
			jk := destType + "," + destKeyValue
			linkListInDB[jk] = []string{destType, destKeyValue}
		}
	}
	// now iterate through all the items in the current link and remove the one's that re not present
	removeList := [][]string{}
	if !rLink {
		node.GetLinks().ForEach(func(ntype, nkey string, l *ifc.GLink) {
			jk := ntype + "," + nkey
			if _, dtOk := linkListInDB[jk]; !dtOk {
				removeList = append(removeList, []string{ntype, nkey})
			}
		})
	} else {
		node.GetLinks().ForEachRLink(func(ntype, nkey string, l *ifc.GLink) {
			jk := ntype + "," + nkey
			if _, dtOk := linkListInDB[jk]; !dtOk {
				removeList = append(removeList, []string{ntype, nkey})
			}
		})
	}
	for _, itm := range removeList {
		nt := itm[0]
		nk := itm[1]
		if !rLink {
			node.DeleteImmediateLink(nt, nk)
		} else {
			node.DeleteImmediateRLink(nt, nk)
		}
	}
}

// populate data read from graph in 'n' into base node obj in 'node'
func (c *DataModelCache) populateCacheLinks(node ifc.BaseNodeInterface,
	nodeData *ifc.GNode, syncSubTrees bool) {
	// if the node is under delete just return with nop.
	if node.IsDeleted() {
		return
	}
	rLinkFeatureFlagEnabled := c.isRLinkFeatureFlagEnabled()
	// basePath := node.GetFullPath()
	// should this call be made before or after the properties are updated
	// for now keeping it after the properties are updated
	logging.Debugf("%s %s::%s %s", c.name, node.GetType(), node.GetKeyValue(), node.GetId())
	node.SetImmediateBaseProperties(nodeData.Properties)
	// let's do a diff and execute the callbacks'
	{
		c.markAndSweepNodeCacheLinks(node, nodeData, false)
		if rLinkFeatureFlagEnabled {
			c.markAndSweepNodeCacheLinks(node, nodeData, true)
		}
	}
	for _, itm := range nodeData.Links {
		logging.Debugf("populateCacheLinks(): %s====> Iterating over LINKS %s", c.name, itm.Properties)
		if !c.dm.IsRunning() {
			break
		}
		_, nkvOk := itm.Properties[common.LinkFixedProp_NodeKeyValue]
		_, ntOk := itm.Properties[common.LinkFixedProp_NodeType]
		if nkvOk && ntOk {
			c.populateSingleCacheLink(node, itm, syncSubTrees)
		} else {
			if c.dm.IsRunning() {
				logging.Errorf("populateCacheLinks(): %s:%s Link missing fixed properties %s", c.name, itm.Id, itm)
			}
		}
	}
	if rLinkFeatureFlagEnabled {
		for _, itm := range nodeData.RLinks {
			logging.Debugf("populateCacheLinks(): %s====> Iterating over RLINKS %s", c.name, itm.Properties)
			if !c.dm.IsRunning() {
				break
			}
			_, nkvOk := itm.Properties[common.LinkFixedProp_NodeKeyValue]
			_, ntOk := itm.Properties[common.LinkFixedProp_NodeType]
			if nkvOk && ntOk {
				c.populateSingleCacheRLink(node, itm, false)
			} else {
				if c.dm.IsRunning() {
					logging.Errorf("populateCacheLinks(): %s:%s RLink missing fixed properties %s", c.name, itm.Id, itm)
				}
			}
		}
	}
}

func (c *DataModelCache) DoNodeUpdateCallbacks(
	base ifc.BaseNodeInterface,
	curProp ifc.PropertyType,
	newProp ifc.PropertyType) {
	c.cacheManager.DoNodeUpdateCallbacks(base, curProp, newProp)
}
func (c *DataModelCache) DoLinkDeleteCallbacks(
	base ifc.BaseNodeInterface,
	destType string,
	dstKey string,
	oldLink ifc.PropertyType) {
	c.cacheManager.DoLinkDeleteCallbacks(base, destType, dstKey, oldLink)
}
func (c *DataModelCache) DoLinkUpdateCallbacks(
	base ifc.BaseNodeInterface,
	dstKey string,
	destType string,
	curLink *ifc.GLink,
	newLink *ifc.GLink) {

	{
		bpath := base.GetFullPath()
		if len(bpath) > 1 && base.GetBaseParent() == nil {
			logging.Warnf("%s : No Parent set during callback processing\n", base.GetId())
			return
		}
	}
	// adding go routing to handle callback, this was needed for soft link sync of destination node
	// if not in gorouting the call to ReadyDestNodeBeforeCB will cause deadlock.

	// TODO: Do this goroutine one level higher where the update type is decided.
	// Having a goroutine here can cause races with event types.
	go c.cacheManager.DoLinkUpdateCallbacks(base, dstKey, destType, curLink, newLink)
}

// this function is to force refresh of a part of the graph starting at the node path and up to
// the depth mentioned. This is used to setup the cache in case data is present in db before the subscribe
// is called.
func (c *DataModelCache) forceRefreshNode(
	path ifc.NodePathList, depth uint32, init bool, nd ifc.BaseNodeInterface) {
	if init && depth > 0 && len(path) > 1 {
		ndList := c.dm.PopulatePathAndFetchNodes(path, true)
		ndl := len(ndList)
		logging.Debugf("forceRefreshNode Init checking on path=%s Returning ndlist.length=%d", path, ndl)
		if len(path) == len(ndList) {
			// we got the node ... now need to recurse
			c.forceRefreshNode(path, depth, false, ndList[len(ndList)-1])
		} else {
			logging.Debugf("forceRefreshNode checking on path %s Returning ndlist.length=%d", path, ndl)
		}
	} else if depth > 0 && nd != nil {
		nd.GetLinks().ForEach(func(nt, nk string, lnk *ifc.GLink) {
			if prop, propOk := lnk.Properties[common.LinkFixedProp_HardLink]; propOk && prop == "true" {
				destNode, _, destNodeOk := nd.GetChild(nk, nt, true)
				logging.Debugf("forceRefreshNode[depth=%d] Trigger sync of %s got %t", depth, lnk.DestinationNodeId, destNode != nil)
				if destNodeOk {
					c.forceRefreshNode(path, depth-1, false, destNode)
				}
			}
		})
	}
}

func (c *DataModelCache) ReadyDestNodeBeforeCB(
	base ifc.BaseNodeInterface,
	destKey, destType string,
	linkProperty ifc.PropertyType) {
	if _, _, destNodeOk := base.GetChild(destKey, destType, false); !destNodeOk {
		if linkProperty[common.LinkFixedProp_HardLink] == "true" {
			c.dm.SyncChild(base, destType, destKey, true)
		} else {
			destNodePathStr := linkProperty[common.LinkFixedProp_SoftLinkDestinationPath]
			destNodePath := utils.NPathFromStr(destNodePathStr.(string))
			c.dm.PopulatePathAndFetchNodes(destNodePath, true)
		}
	}
}

func (c *DataModelCache) DumpCachedNodes() map[string]interface{} {
	return c.cacheManager.DumpCachedNodes()
}
