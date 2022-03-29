package datamodel

import (
	"bytes"
	"errors"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/internal/cache"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/internal/dmbus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/internal/graphdb"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/internal/utils"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/dmcache"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/ifc"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/logging"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/scheduler"

	"github.com/google/uuid"
)

type DataModel struct {
	graph ifc.GraphDBInterface
	dmbus ifc.DMBusInterface
	// config
	cache        ifc.DataModelCacheInterface
	id           string
	scheduler    *scheduler.Scheduler
	running      bool
	ccnt         uint64
	name         string
	featureFlags []string
}

var dataModelCreateCounter uint64 = 1

func (dm *DataModel) processDMFeatureFlags(featureFlags string) {
	flagList := strings.Split(featureFlags, ";")
	if len(flagList) == 0 {
		return
	}
	for _, v := range flagList {
		dm.featureFlags = append(dm.featureFlags, v)
	}
}

func (dm *DataModel) IsRLinkFeatureFlagEnabled() bool {
	for _, v := range dm.featureFlags {
		if strings.Compare(v, common.Rlink_EnableFeatureFlag) == 0 {
			return true
		}
	}
	return false
}

func (dm *DataModel) GetLatestRevision() int64 {
	return dm.graph.GetLatestRevision()
}

func (dm *DataModel) CompactRevision(revision, compactCtxtTimeout int64, compactPhysical bool) (*ifc.EtcdCompactResponse, error) {
	return dm.graph.CompactRevision(revision, compactCtxtTimeout, compactPhysical)
}

func NewDataModel(name, etcdLoc, featureFlags string) *DataModel {
	dm := &DataModel{
		graph:        nil,
		name:         name,
		id:           uuid.New().String(),
		cache:        nil,
		featureFlags: []string{},
		// readOnlyMode: false,
		scheduler: scheduler.NewScheduler(),
		running:   false}
	dm.processDMFeatureFlags(featureFlags)
	dm.ccnt = atomic.AddUint64(&dataModelCreateCounter, 1)
	dm.cache = cache.NewDataModelCache(name, dm, dm.featureFlags)
	dm.graph = graphdb.New(name, etcdLoc, dm.featureFlags)
	dm.dmbus = dmbus.New(name, dm.id)
	runtime.SetFinalizer(dm, func(dt *DataModel) {
		logging.Infof("DataModel %s is Garbage Collected\n", dt.GetName())
	})

	logging.Debugf("CREATE: New Data Model %s[%+v] ccnt=%d Done.\n", name, *dm, dm.ccnt)
	return dm
}

func (dm *DataModel) Init() {
	/*
	   this.graph = await createGraphDB(this.config);
	   this.dmbus = new DMBus(this.config, dt);
	   await this.dmbus.init(this.graph);
	   this.cache.init(dt, this.dmbus);
	   this.running = true;
	   logger.debug(`dm ${dt} starting with id ${this.id}`);
	*/
	dm.cache.Init(dm.dmbus)
	dm.dmbus.Init(dm.graph)
	dm.running = true

	// Check for debug server
	hasEnabled, port := dm.checkForDMDebugServeEnableAndPort()
	if hasEnabled {
		if !dmcache.DebugServerStarted {
			dmcache.DebugServerStarted = true
			go dmcache.StartDebugServer(port)
		}
		dm.SetDMDebugCache()
	}
}

func (dm *DataModel) checkForDMDebugServeEnableAndPort() (bool, string) {
	for _, v := range dm.featureFlags {
		items := strings.Split(v, ":")
		if strings.Compare(items[0], common.DM_Debug_Server_Enable) == 0 {
			logging.Debugf("DM Debug Sever flag (%s) enabled", items[0])
			if len(items) == 2 {
				return true, items[1] //true, debug_server_port
			}
			return true, ""
		}
	}
	return false, ""
}

func (dm *DataModel) SetDMDebugCache() {
	dmcache.GlobalDMCacheMutex.Lock()
	defer dmcache.GlobalDMCacheMutex.Unlock()

	dmcache.GlobalDMCacheMap[dm.GetId()] = dm
}

func (dm *DataModel) SetGraph(graph ifc.GraphDBInterface) {
	dm.graph = graph
}
func (b *DataModel) GetCCnt() uint64 {
	return b.ccnt
}
func (b *DataModel) GetName() string {
	return b.name
}
func (dm *DataModel) SetMessagingDelay(n uint32) {
	dm.dmbus.AddMessageDelay(n)
}
func (dm *DataModel) IsRunning() bool {
	return dm.running
}
func (dm *DataModel) GetId() string {
	return dm.id
}
func (dm *DataModel) SetId(id string) {
	ns := uuid.MustParse("55e5b027-3214-1234-2222-7e58898f2b31")
	dm.id = uuid.NewMD5(ns, []byte(id)).String()
}

func (dm *DataModel) DeleteDMFromDebugCache() {
	dmcache.GlobalDMCacheMutex.Lock()
	defer dmcache.GlobalDMCacheMutex.Unlock()

	delete(dmcache.GlobalDMCacheMap, dm.GetId())
}

func (dm *DataModel) Shutdown() {
	if dm.running {
		dm.running = false
		dm.cache.Stop()
		dm.dmbus.Shutdown()
		dm.graph.Shutdown()
		logging.Infof("SHUTDOWN: Done Stopping Data Model %s[%p] Shutdown called running=%t\n", dm.name, dm, dm.running)
	}

	hasEnabled, _ := dm.checkForDMDebugServeEnableAndPort()
	if hasEnabled {
		dm.DeleteDMFromDebugCache()
	}
}

func (dm *DataModel) Sync(node ifc.BaseNodeInterface, forceSync bool) {
	lck := dm.scheduler.Wait(node.GetId())
	defer dm.scheduler.Done(lck)
	dm.cache.Sync(node, forceSync, false)
}

func (dm *DataModel) Delete(node ifc.BaseNodeInterface) {
	dm.DeleteId(node.GetId(), false)
}

func (dm *DataModel) DeleteFSync(node ifc.BaseNodeInterface) {
	dm.DeleteId(node.GetId(), true)
}

func (dm *DataModel) DeleteId(nodeId string,
	forceDelete bool) {
	dm.deleteTree(nodeId, forceDelete, 0)
}

func (dm *DataModel) GetRootNode(nodeId string, dbSync bool) (ifc.BaseNodeInterface, bool) {
	return dm.cache.GetCachedNodeSync(nodeId)
}

func (dm *DataModel) GetNode(parentId, nodeType, nodeId string, dbSync bool, syncSubTrees bool) (ifc.BaseNodeInterface, bool) {
	logging.Debugf("getNode before scheduler for nodeId=%s", nodeId)
	lck := dm.scheduler.Wait(nodeId)
	defer dm.scheduler.Done(lck)
	logging.Debugf("getNode after scheduler for nodeId=%s", nodeId)
	return dm.cache.GetNode(parentId, nodeId, nodeType, dbSync, syncSubTrees)
}

func (dm *DataModel) GetCachedNode(nodeId string) (ifc.BaseNodeInterface, bool) {
	return dm.cache.GetCachedNodeSync(nodeId)
}
func (dm *DataModel) IsNodeInCache(n string) bool {
	return dm.cache.IsNodeInCache(n)
}
func (dm *DataModel) GetNodeFromCache(id string) (ifc.BaseNodeInterface, bool) {
	return dm.cache.GetCachedNode(id, false)
}

func (dm *DataModel) PopulatePathAndFetchNodes(path ifc.NodePathList, forceSync bool) []ifc.BaseNodeInterface {
	spath := utils.JsonMarshal(path)
	logging.Debugf("->populatePathAndFetchNode called for %s [%t]", spath, forceSync)
	// from top of the path list
	// fetch every node starting at the root
	// and linking the childrent
	// till the bottom node is fetched
	// then return the bottom and the parent node
	retNodes := []ifc.BaseNodeInterface{}
	root := true
	var nd ifc.BaseNodeInterface
	var ok bool
	for _, itm := range path {
		ntype := itm[ifc.NodePathName_nodeType]
		nkey := itm[ifc.NodePathName_nodeIdentifier]
		if root {
			// root is always in cache.
			nd, ok = dm.cache.GetCachedNode("root", false)
			if !ok {
				nd, ok = dm.cache.GetCachedNode("/Root/root", false)
			}
			if ok {
				retNodes = append(retNodes, nd)
			} else {
				return retNodes
				//panic(errors.New("Undefined Root Node!!! name=" + nkey))
			}
			root = false
		} else {
			// next node is from parent
			if ok && nd.GetBaseParent() != nil && !nd.GetLinks().Has(ntype, nkey) {
				// try to sync the node if the edge is missing
				dm.cache.GetNode(nd.GetBaseParent().GetId(), nd.GetId(), nd.GetType(),
					forceSync, false)
			}
			if nd != nil && nd.GetLinks().HasType(ntype) {
				if lnk, ok := nd.GetLinks().Get(ntype, nkey); ok {
					cid := lnk.DestinationNodeId
					var ndNext ifc.BaseNodeInterface
					if ndNext, ok = dm.cache.GetCachedNode(cid, forceSync); !ok {
						nodeType, _ := lnk.Properties[common.LinkFixedProp_NodeType]
						ndNext, _ = dm.cache.GetNode(nd.GetId(), cid, nodeType.(string),
							forceSync, false)
					}
					nd = ndNext
				} else {
					return []ifc.BaseNodeInterface{}
				}
				if nd != nil {
					retNodes = append(retNodes, nd)
				} else {
					// one of the nodes in the path is undefined
					logging.Errorf(
						"Unable to populatePathAndFetchNode with %s when processing %s.%s, one of the node in path is undefined",
						spath, ntype, nkey)
					return []ifc.BaseNodeInterface{}
				}
			} else {
				// could not follow path
				logging.Debugf(
					"Unable to populatePathAndFetchNode with %s when processing %s.%s",
					spath, ntype, nkey)
				return []ifc.BaseNodeInterface{}
			}
		}
	}
	// reached the end
	return retNodes
}

func (dm *DataModel) UpsertNode(
	parentId,
	linkType,
	linkKey string,
	linkProp ifc.PropertyType,
	nodeType string,
	nodeProp ifc.PropertyType) (ifc.BaseNodeInterface, *ifc.GLink) {
	var key string = common.NodeFixedProp_NodeDefaultKeyName
	if linkKey != "" {
		key = linkKey
	}
	cid := parentId + "/" + nodeType + "/" + nodeProp[key].(string)
	lck := dm.scheduler.Wait(cid)
	defer dm.scheduler.Done(lck)
	_, linkKeyInNP := nodeProp[linkKey]
	if linkKey == "" || !linkKeyInNP {
		panic(errors.New("linkKey " + linkKey + " is missing in node.properties"))
	}
	// what if graph has an older node that is populated
	newNode, newLink := dm.graph.UpsertChildNode(
		dm.id, parentId, linkType, linkKey, linkProp, nodeType, nodeProp)
	// create a base node from the data returned by graph access
	retNode := dm.cache.PopulateCacheAndReturnBaseNode(parentId, newNode, linkKey)

	pnode, pnodeOk := dm.cache.GetCachedNode(parentId, false)
	if pnodeOk {
		retNode.SetBaseParent(pnode)
		k := nodeProp[linkKey].(string)
		pnode.UpsertImmediateLink(nodeType, k, newLink)
	}
	return retNode, newLink
}

func (dm *DataModel) UpsertRootNode(nodeType, nodeName string, nodeProp ifc.PropertyType) ifc.BaseNodeInterface {
	if nodeName == "" {
		panic(errors.New("UpsertRootNode: nodeName is empty."))
	}
	logging.Debugf("upsertRootNode %s.%s", nodeType, nodeName)
	newNode := dm.graph.UpsertNode(dm.id, nodeType, nodeName, nodeProp)
	retNode := dm.cache.PopulateCacheAndReturnBaseNode("", newNode, common.NodeFixedProp_NodeDefaultKeyName)
	retNode.SetBaseParent(nil)
	return retNode
}
func (dm *DataModel) UpdateNodeAddProperties(nodeId string, nodeProp ifc.PropertyType) {
	logging.Debugf("UpdateNodeAddProperties id=%s prop=%s", nodeId, nodeProp)
	lck := dm.scheduler.Wait(nodeId)
	defer dm.scheduler.Done(lck)
	if len(nodeProp) == 0 {
		return
	}
	// get node from cache
	if node, ok := dm.cache.GetCachedNode(nodeId, false); ok {
		dm.graph.UpdateNodeAddProperties(dm.id, nodeId, nodeProp)
		spath := utils.JsonMarshal(node.GetFullPath())
		logging.Debugf("Update Node add Prop sending a message for path %s", spath)
		dm.SyncNodeProperty(node, true)
	} else {
		logging.Debugf("UpdateNodeAddProperties: Did not find the node, dropping update")
	}
}
func (dm *DataModel) UpdateNodeRemoveProperties(nodeId string, nodeProp []string) {
	logging.Debugf("UpdateNodeRemoveProperties id=%s prop=%s", nodeId, nodeProp)
	lck := dm.scheduler.Wait(nodeId)
	defer dm.scheduler.Done(lck)
	if len(nodeProp) == 0 {
		return
	}
	// get node from cache
	if node, ok := dm.cache.GetCachedNode(nodeId, false); ok {
		dm.graph.UpdateNodeRemoveProperties(dm.id, nodeId, nodeProp)
		dm.SyncNodeProperty(node, true)
	} else {
		logging.Debugf("UpdateNodeRemoveProperties:Did not find the node, dropping Delete")
	}
}

func (dm *DataModel) UpsertLink(lnk common.UpsertLinkOpObj) *ifc.GLink {
	logging.Debugf("Upsertlink(): %s--->%s  Acquiring lock on link dest node to add link.\n",
		lnk.SrcNodeObj.NodeId, lnk.DestNodeObj.NodeId)
	lck := dm.scheduler.Wait(lnk.SrcNodeObj.NodeId)
	newLink, newRLink := dm.graph.UpsertLink(
		dm.id, lnk)
	if pnode, pok := dm.cache.GetCachedNode(lnk.SrcNodeObj.NodeId, true); pok {
		pnode.UpsertImmediateLink(lnk.DestNodeObj.NodeType, lnk.DestNodeObj.NodeKeyValue, newLink)
		dm.scheduler.Done(lck)
		logging.Debugf("Upsertlink(): %s--->%s  Release lock on link dest node to add link.\n",
			lnk.SrcNodeObj.NodeId, lnk.DestNodeObj.NodeId)

		/*
			Create reverse links '_rlinks' path in the db from dst to src node via Upsert
			Wrap this op within a feature flag check
		*/
		if dm.IsRLinkFeatureFlagEnabled() {
			dstNodeId := lnk.DestNodeObj.NodeId
			logging.Debugf("Upsertlink(): %s--->%s  Acquiring lock on link dest node to add rlink.\n",
				lnk.SrcNodeObj.NodeId, lnk.DestNodeObj.NodeId)
			lck1 := dm.scheduler.Wait(lnk.DestNodeObj.NodeId)
			if dstNode, ok := dm.cache.GetCachedNode(lnk.DestNodeObj.NodeId, false); ok {
				if dstNode != nil {
					dstNode.UpsertImmediateRLink(lnk.SrcNodeObj.NodeType, lnk.SrcNodeObj.NodeKeyValue, newRLink)
					logging.Debugf("Reverse link added from %s-->%s",
						dstNodeId, lnk.SrcNodeObj.NodeId)
				}
			}
			logging.Debugf("Upsertlink(): %s--->%s  Release lock on link dest node to add rlink.\n",
				lnk.SrcNodeObj.NodeId, lnk.DestNodeObj.NodeId)
			dm.scheduler.Done(lck1)
		}
		return newLink
	}
	dm.scheduler.Done(lck)
	return newLink
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
func (dm *DataModel) DeleteLink(srcNodeId, dstNodeType, dstKeyValue string) {
	var destNodeId string
	srcNodeType := ""
	srcNodeKey := ""
	destLinkId := ""
	lck := dm.scheduler.Wait(srcNodeId)
	if pnode, pok := dm.cache.GetCachedNode(srcNodeId, false); pok {
		if link, linkOk := pnode.GetLinks().Get(dstNodeType, dstKeyValue); linkOk {
			srcNodeIdParts := strings.Split(srcNodeId, "/")
			srcNodeIdPartsLen := len(srcNodeIdParts)
			if srcNodeIdPartsLen > 1 {
				srcNodeType = srcNodeIdParts[srcNodeIdPartsLen-2]
				srcNodeKey = srcNodeIdParts[srcNodeIdPartsLen-1]
			}
			destNodeId = link.DestinationNodeId
			if destNodeId != "" && srcNodeType != "" && srcNodeKey != "" && dm.IsRLinkFeatureFlagEnabled() {
				destLinkId = destNodeId + "/_rlinks/" + srcNodeType + "/" + srcNodeKey
				logging.Debugf("DeleteSoftLinkWithRSoftLink(): lnk, rlnk: %s, %s\n",
					link.Id, destLinkId)
				dm.graph.DeleteSoftLinkWithRSoftLink(link.Id, destLinkId)
				pnode.DeleteImmediateLink(dstNodeType, dstKeyValue)
				dm.scheduler.Done(lck)
				logging.Debugf("deleteLink(): %s--->%s  Acquiring lock on link dest node to delete rlink.\n",
					srcNodeId, destNodeId)
				lck1 := dm.scheduler.Wait(destNodeId)
				if dnode, dok := dm.cache.GetCachedNode(destNodeId, false); dok {
					if dnode != nil {
						if _, dlinkOk := dnode.GetLinks().GetRLink(srcNodeType, srcNodeKey); dlinkOk {
							dnode.DeleteImmediateRLink(srcNodeType, srcNodeKey)
						}
					}
				}
				dm.scheduler.Done(lck1)
				logging.Debugf("deleteLink(): %s--->%s  Release lock on link dest node to delete rlink.\n",
					srcNodeId, destNodeId)
			} else {
				dm.graph.DeleteLink(link.Id)
				pnode.DeleteImmediateLink(dstNodeType, dstKeyValue)
				dm.scheduler.Done(lck)
			}
			return
		}
	}
	dm.scheduler.Done(lck)
}

/*
Creating a new API here to preserve backward compatibility until the existing API's are migrated
*/

func (dm *DataModel) DeleteRLink(srcNodeId, dstNodeType, dstKeyValue string) {
	lck := dm.scheduler.Wait(srcNodeId)
	defer dm.scheduler.Done(lck)
	if pnode, pok := dm.cache.GetCachedNode(srcNodeId, false); pok {
		if pnode != nil {
			if link, linkOk := pnode.GetLinks().GetRLink(dstNodeType, dstKeyValue); linkOk {
				dm.graph.DeleteLink(link.Id)
				pnode.DeleteImmediateRLink(dstNodeType, dstKeyValue)
			}
		}
	}
}

func (dm *DataModel) Subscribe(pattern ifc.NodePathList, depth uint32) {
	dm.cache.AddSubscription(pattern, depth)
}
func (dm *DataModel) Unsubscribe(pattern ifc.NodePathList) {
	dm.cache.DelSubscription(pattern)
}

func (dm *DataModel) GetDBStats() common.GraphDBStats {
	return dm.graph.GetStats()
}
func (dm *DataModel) GetMsgStats() ifc.NotificationStats {
	return dm.dmbus.GetStats()
}
func (dm *DataModel) RegisterCB(pattern ifc.NodePathList, cbfn ifc.CallbackFuncNode) {
	dm.cache.RegisterNodeCB(pattern, &cbfn)
}

func (dm *DataModel) RegisterLinkCB(pattern ifc.NodePathList, destNodeType string,
	cbfn ifc.CallbackFuncLink) {
	dm.cache.RegisterLinkCB(pattern, &cbfn, destNodeType)
}

func (dm *DataModel) ClaimOwnership(nodePath ifc.NodePathList, timeDelay uint32, prefix string) {
	logging.Debugf("claimOwnership Requested for %s timeDelay=%s prefix=%s.", nodePath, timeDelay, prefix)
	time.Sleep(time.Duration(timeDelay) * time.Millisecond)
	nodePathStr := utils.JsonMarshal(nodePath)
	logging.Debugf("claimOwnership Triggered for %s Register for CB.", nodePathStr)
	dm.RegisterCB(
		nodePath,
		func(base ifc.BaseNodeInterface,
			updateType int,
			oldData, newData ifc.PropertyType) {
			if updateType == common.UpdateType_NodeDelete {
				return
			}
			// check if the owner id and delete if not matching dm id
			nodeUpdatedByT, nodeUpdatedByOk := newData[common.NodeFixedProp_updatedBy]
			nodeUpdatedBy := ""
			if nodeUpdatedByT != nil {
				nodeUpdatedBy = nodeUpdatedByT.(string)
			}
			nodeCreatedBy := newData[common.NodeFixedProp_createdBy].(string)
			logging.Debugf("claimOwnership CB node:%s %d owner id=%s nodeUpdate/Created=%s/%s",
				base.GetId(), updateType, dm.GetId(), nodeUpdatedBy, nodeCreatedBy)
			oid := nodeUpdatedBy
			if !nodeUpdatedByOk {
				oid = nodeCreatedBy
			}
			if oid != "" && oid != dm.GetId() && (prefix == "" || (strings.HasPrefix(base.GetKeyValue(), prefix))) {
				spath := utils.JsonMarshal(base.GetFullPath())
				logging.Infof("GC:Deleting Node %s due to ownership change OldId: %s, dmId: %s, %s name: %s, prefix: %s",
					spath, oid, dm.GetId(), nodeUpdatedBy, base.GetKeyValue(), prefix)
				base.Delete()
			}
		})
}

func (dm *DataModel) ClaimLinkOwnership(srcNodePath ifc.NodePathList, destNodeType string, timeDelay uint32) {
	srcNodePathStr := utils.JsonMarshal(srcNodePath)
	logging.Debugf("ClaimLinkOwnership Requested for %s dest=%s timeDelay=%s", srcNodePathStr, destNodeType, timeDelay)
	time.Sleep(time.Duration(timeDelay) * time.Millisecond)
	logging.Debugf("ClaimLinkOwnership Triggered for %s Register for CB.", srcNodePathStr)
	dm.RegisterLinkCB(
		srcNodePath,
		destNodeType,
		func(baseSrc ifc.BaseNodeInterface,
			updateType int, key string, oldData, newData ifc.PropertyType) {
			if updateType != common.UpdateType_LinkDelete {
				linkUpdateByT, linkUpdatedByOk := newData[common.LinkFixedProp_updatedBy]
				linkUpdatedBy := linkUpdateByT.(string)
				linkCreatedBy := newData[common.LinkFixedProp_createdBy].(string)
				oid := linkUpdatedBy
				if !linkUpdatedByOk {
					oid = linkCreatedBy
				}
				if oid != "" && oid != dm.GetId() {
					baseSrcPathstr := utils.JsonMarshal(baseSrc.GetFullPath())
					logging.Infof("GC:Deleting Link from Nod e%s Link to %s.%s due to ownership change OldId: %s, newId: %s updatedby/createdby= %s/%s",
						baseSrcPathstr, destNodeType, key, oid, dm.GetId(), linkUpdatedBy, linkCreatedBy)
					baseSrc.DeleteLink(key, destNodeType)
				}
			}
		})
}

func (dm *DataModel) UpdateLinkAddProperty(
	srcNodeId,
	dstNodeType,
	dstKeyValue string,
	linkProp ifc.PropertyType) bool {
	lck := dm.scheduler.Wait(srcNodeId)
	defer dm.scheduler.Done(lck)
	if pnode, pnodeOk := dm.cache.GetCachedNode(srcNodeId, true); pnodeOk {
		if link, linkOk := pnode.GetLinks().Get(dstNodeType, dstKeyValue); linkOk {
			rev := dm.graph.UpdateLinkAddProperties(dm.id, link.Id, linkProp)
			// need the link again as it gets re-written by graph update.
			nlink := ifc.GLink{
				Id:                link.Id,
				LinkType:          link.LinkType,
				SourceNodeId:      link.SourceNodeId,
				DestinationNodeId: link.DestinationNodeId,
				Properties:        make(ifc.PropertyType)}
			for key, val := range link.Properties {
				nlink.Properties[key] = val
			}
			for key, val := range linkProp {
				nlink.Properties[key] = val
			}
			nlink.Properties[common.LinkFixedProp_Revision] = rev
			pnode.UpsertImmediateLink(dstNodeType, dstKeyValue, &nlink)
			return true
		}
	}
	return false
}

func (dm *DataModel) UpdateLinkRemoveProperty(
	srcNodeId,
	dstNodeType,
	dstKeyValue string,
	linkPropKeys []string) bool {
	lck := dm.scheduler.Wait(srcNodeId)
	defer dm.scheduler.Done(lck)
	if pnode, pnodeOk := dm.cache.GetCachedNode(srcNodeId, true); pnodeOk {
		if link, linkOk := pnode.GetLinks().Get(dstNodeType, dstKeyValue); linkOk {
			rev := dm.graph.UpdateLinkRemoveProperties(dm.id, link.Id, linkPropKeys)
			// it needs to be refreshed after db read
			nlink := ifc.GLink{
				Id:                link.Id,
				LinkType:          link.LinkType,
				SourceNodeId:      link.SourceNodeId,
				DestinationNodeId: link.DestinationNodeId,
				Properties:        make(ifc.PropertyType)}
			for k, v := range link.Properties {
				nlink.Properties[k] = v
			}
			for _, key := range linkPropKeys {
				delete(nlink.Properties, key)
			}
			nlink.Properties[common.LinkFixedProp_Revision] = rev
			pnode.UpsertImmediateLink(dstNodeType, dstKeyValue, &nlink)
			return true
		}
	}
	return false
}

func (dm *DataModel) GetNextChildKeyList(
	parentId, nodeType, childStartingKey string, batchSize uint32) []string {
	return dm.graph.GetNextLinkKey(parentId, nodeType, childStartingKey, int(batchSize))
}

func (dm *DataModel) GetNextRLinkChildKeyList(
	parentId, nodeType, childStartingKey string, batchSize uint32) []string {
	return dm.graph.GetNextRLinkKey(parentId, nodeType, childStartingKey, int(batchSize))
}

func (dm *DataModel) SyncNodeProperty(node ifc.BaseNodeInterface, forced bool) {
	dm.cache.SyncNodeProperty(node, forced)
}

func (dm *DataModel) SyncChild(pnode ifc.BaseNodeInterface, childNodeType, childNodeKey string, forced bool) {
	dm.cache.SyncChild(pnode, childNodeType, childNodeKey, forced)
}

func (dm *DataModel) GetGraph() ifc.GraphDBInterface {
	return dm.graph
}

func (dm *DataModel) deleteTree(nodeId string, forceDelete bool, depth int) ifc.NodePathList {
	lck := dm.scheduler.Wait(nodeId)
	var rlinks []*ifc.GLink
	if nodeId == "" {
		panic(errors.New("deleteTree called with nodeid Set to " + nodeId))
	}
	node, nodeOk := dm.cache.GetCachedNode(nodeId, false)
	nDeleted := false
	if nodeOk {
		nDeleted = node.IsDeleted()
	}
	logging.Debugf("%s: Deleting node %s got node from nodeOk=%t deleted=%t", dm.name,
		nodeId, nodeOk, nDeleted)
	/*
		Need to include a check for node.isToBeDeleted()
		Not required now to preserve backward compatibility
	*/
	if !nodeOk || node.IsDeleted() {
		dm.scheduler.Done(lck)
		return ifc.NodePathList{}
	}

	if node.IsAnyParentDeleted() {
		logging.Debugf("%s: One of the Parent is deleted so another delete is in progress, will skip this delete for %s", dm.name, nodeId)
	}

	/*
		Set delete property in the node to be true.
		This is required to handle deletions when a service/consumer executing delete op
		crashes and resumes to process the node following bootstrap phase.
		The logic to process this property needs to be handled in all consumer API's of a node, to ensure
		that upserts and creates are not executed if a node is marked to be deleted.
		The scope of this behavior extends beyond current commit.
		Current commit aims to introduce a 'NodeFixedProp_ToBeDeleted' in the Node property object.
	*/
	dm.scheduler.Done(lck)
	dm.UpdateNodeAddProperties(nodeId, map[string]interface{}{common.NodeFixedProp_ToBeDeleted: true})
	lck1 := dm.scheduler.Wait(nodeId)

	nodePath := node.GetFullPath()
	node.SetDeleted()
	// remove all the children
	node.GetLinks().ForEach(func(ntype, childKey string, lnk *ifc.GLink) {
		if lnk.Properties[common.LinkFixedProp_HardLink] == "true" {
			dm.deleteTree(lnk.DestinationNodeId, forceDelete, depth+1)
		}

	})

	/*
		Delete forward and reverse soft links for the given node in the database
		Wrap this op in a feature flag check
	*/
	if dm.IsRLinkFeatureFlagEnabled() {
		rlinks = dm.graph.DeleteLinks(nodeId)

		/*
			We need the list of rlinks to be returned from graphDB API
			to remove the links from basenode object
			Wrap this op within rlink feature flag check
		*/
		for _, rlink := range rlinks {
			if rlink.DestinationNodeId != "" {
				rDestNodeIdSplit := strings.Split(rlink.DestinationNodeId, "/")
				rDestNodeIdSplitLen := len(rDestNodeIdSplit)
				if rDestNodeIdSplitLen > 1 {
					rDestNodeType := rDestNodeIdSplit[rDestNodeIdSplitLen-2]
					rDestNodeKey := rDestNodeIdSplit[rDestNodeIdSplitLen-1]
					node.DeleteImmediateRLink(rDestNodeType, rDestNodeKey)
				}
			}
		}
	}
	/*
		Not required, but preserving for backward compatibility,
		until new workflow is deemed to be fit
	*/
	node.ReverseLinkIterate(func(tlinkkey, tnodeId, linkId string) {
		logging.Debugf("%s: Sushil 4 %s deleting reverse link \n", dm.name, linkId)
		// delete any reverse links for this node.
		if linkId != "" {
			dm.graph.DeleteLink(linkId)
		} else {
			dm.DeleteLink(tnodeId, node.GetType(), tlinkkey)
		}
	})

	logging.Debugf("%s: seleting Node in Graph %s", dm.name, nodeId)
	dm.graph.DeleteNode(nodeId)
	dm.cache.PurgeCacheTree(nodeId)
	dm.scheduler.Done(lck1)
	// Check if the rlink Dest node is present in the cache
	// If yes, delete links member for the rlink dest node
	if dm.IsRLinkFeatureFlagEnabled() {
		for _, rlink := range rlinks {
			if rlink.DestinationNodeId != "" {
				logging.Debugf("deletetree(): nodeId:%s, Acquiring lock on rlink dest node %s\n",
					nodeId, rlink.DestinationNodeId)
				// Get depth of rlink node and current node
				nodeIdSplit := strings.Split(nodeId, "/")
				nodeIdSplitLen := len(nodeIdSplit) - 1
				nodeIdDepth := nodeIdSplitLen / 2
				rnodeIdSplit := strings.Split(rlink.DestinationNodeId, "/")
				rnodeIdSplitLen := len(rnodeIdSplit) - 1
				rnodeIdDepth := rnodeIdSplitLen / 2
				lckAcquire := true
				if depth >= nodeIdDepth-rnodeIdDepth && strings.Contains(nodeId, rlink.DestinationNodeId) == true {
					lckAcquire = false
				}
				// Important: We need this check to prevent deadlocks in case of node deletions where
				// one of the ancestor nodes has a soft link to the current node Id
				// for ex: nodeId:/Root/root/Config/default/UserFolder/default/SoftServe/default---> rlink ---> /Root/root/Config/default
				// /Root/root/Config/default ---> soft link --> /Root/root/Config/default/UserFolder/default/SoftServe/default
				// In this case, a node delete triggered for Config node would hit 1st level recursion of userfolder node and
				// 2nd level recursion of softserve node. Now softserve has a rlink to config node and a lock on config node
				// could be acquired here. However, since config node is already locked since the delete op started on the
				// config node, we hit a deadlock. Hence, a depth factor is introduced, which measures the relative depth
				// of the node Id being deleted with the depth of the rlink node ID.
				// In case the delete op for current node is a part of a recursion chain triggered via a delete op for an
				// ancestor node that is either equal to the rlink dest or an ancestor of the rlink dest, we do not
				// acquire the lock and we know that the lock has been acquired in the current recursive call chain.
				var lck3 scheduler.SchedulerHandler
				if lckAcquire {
					lck3 = dm.scheduler.Wait(rlink.DestinationNodeId)
				}
				logging.Debugf("deletetree(): nodeId:%s, lock acquired on rlink dest node %s.\n",
					nodeId, rlink.DestinationNodeId)
				if rdNode, rdok := dm.cache.GetCachedNode(rlink.DestinationNodeId, false); rdok {
					if rdNode != nil {
						sNodeType := node.GetType()
						sNodeKey := node.GetKeyValue()
						if _, rdLinkOk := rdNode.GetLinks().Get(sNodeType, sNodeKey); rdLinkOk {
							rdNode.DeleteImmediateLink(sNodeType, sNodeKey)
						}
					}
				}
				if lckAcquire {
					dm.scheduler.Done(lck3)
				}
				logging.Debugf("deletetree(): nodeId:%s, Release lock on rlink dest node %s.\n",
					nodeId, rlink.DestinationNodeId)
			}
		}
	}
	return nodePath
}
func (c *DataModel) DoNodeUpdateCallbacks(
	base ifc.BaseNodeInterface,
	curProp ifc.PropertyType,
	newProp ifc.PropertyType) {
	c.cache.DoNodeUpdateCallbacks(base, curProp, newProp)
}
func (c *DataModel) DoLinkDeleteCallbacks(
	base ifc.BaseNodeInterface,
	destType string,
	dstKey string,
	oldLink ifc.PropertyType) {
	c.cache.DoLinkDeleteCallbacks(base, destType, dstKey, oldLink)
}
func (c *DataModel) DoLinkUpdateCallbacks(
	base ifc.BaseNodeInterface,
	dstKey string,
	destType string,
	curLink *ifc.GLink,
	newLink *ifc.GLink) {
	c.cache.DoLinkUpdateCallbacks(base, dstKey, destType, curLink, newLink)
}

func (dm *DataModel) CollectDMCachedData() map[string]interface{} {
	return dm.cache.DumpCachedNodes()
}
