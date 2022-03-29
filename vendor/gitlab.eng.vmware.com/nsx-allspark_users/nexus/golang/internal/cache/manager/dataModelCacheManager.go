package manager

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	. "gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/internal/cache/subnode"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/internal/utils"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/ifc"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/scheduler"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/logging"
)

type CacheManagerSubListMap struct {
	data sync.Map
}

func NewCacheManagerSubListMap() *CacheManagerSubListMap {
	return &CacheManagerSubListMap{}
}
func (c *CacheManagerSubListMap) Set(key string, dt uint32) {
	c.data.Store(key, dt)
}
func (c *CacheManagerSubListMap) Get(key string) (uint32, bool) {
	dt, ok := c.data.Load(key)
	if ok {
		return dt.(uint32), ok
	} else {
		return 0, ok
	}
}
func (c *CacheManagerSubListMap) Del(key string) {
	c.data.Delete(key)
}

type CacheManagerNodeCacheMap struct {
	data sync.Map
}

func NewCacheManagerNodeCacheMap() *CacheManagerNodeCacheMap {
	return &CacheManagerNodeCacheMap{}
}

func (c *CacheManagerNodeCacheMap) Set(key string, dt ifc.BaseNodeInterface) {
	c.data.Store(key, dt)
}
func (c *CacheManagerNodeCacheMap) Get(key string) (ifc.BaseNodeInterface, bool) {
	dt, ok := c.data.Load(key)
	if ok {
		return dt.(ifc.BaseNodeInterface), ok
	} else {
		return nil, ok
	}
}
func (c *CacheManagerNodeCacheMap) Del(key string) {
	c.data.Delete(key)
}

type CacheManagerSubTreeMap struct {
	data sync.Map
}

func NewCacheManagerSubTreeMap() *CacheManagerSubTreeMap {
	return &CacheManagerSubTreeMap{}
}
func (c *CacheManagerSubTreeMap) Set(key string, dt *SubNodeType) {
	c.data.Store(key, dt)
}
func (c *CacheManagerSubTreeMap) Get(key string) (*SubNodeType, bool) {
	dt, ok := c.data.Load(key)
	if ok {
		return dt.(*SubNodeType), ok
	} else {
		return nil, ok
	}
}
func (c *CacheManagerSubTreeMap) Del(key string) {
	c.data.Delete(key)
}

// DEBUG helper to print the full subscription state out.
func (c *CacheManagerSubTreeMap) printNodes(node *SubNodeType, ntype string, nvalue string, depth int, st *[]string) {
	lst := ""
	for i := 0; i < depth; i++ {
		lst = lst + "\t"
	}
	lst = lst + fmt.Sprintf("\tType.value: %s.%s, Weight: %d, Depth/sub: %d/%d cachedNodes: %v", ntype, nvalue, node.GetWeight(),
		node.GetDepth(), node.GetCbfnNodeLength()+node.CbfnLinkLength(), node.GetCachedNodeKeys())
	*st = append(*st, lst)
	node.ChildForEach(func(ntype, nvalue string, nd *SubNodeType) {
		c.printNodes(nd, ntype, nvalue, depth+1, st)
	})
}
func (c *CacheManagerSubTreeMap) Print() {
	logging.Debugf("Printing the current subscription Tree")
	st := []string{}
	c.data.Range(func(k, v interface{}) bool {
		c.printNodes(v.(*SubNodeType), k.(string), "", 2, &st)
		return true
	})
	logging.Debugf("%s", strings.Join(st, "\n"))
}

type CacheManager struct {
	fetchAll  bool
	subTree   *CacheManagerSubTreeMap   // map[string]*SubNodeType
	nodeCache *CacheManagerNodeCacheMap // map[string]ifc.BaseNodeInterface
	subList   *CacheManagerSubListMap   // map[string]uint32
	dmc       ifc.DataModelCacheInterface
	cbsch     *scheduler.Scheduler
	name      string
}

// func (cm *CacheManager) SetLog(log *js.Object) {
// 	cm.log = log
// }
func NewCacheManager(name string, dmc ifc.DataModelCacheInterface) *CacheManager {
	e := &CacheManager{
		name:      name,
		fetchAll:  false,
		subTree:   NewCacheManagerSubTreeMap(),   //make(map[string]*SubNodeType),
		nodeCache: NewCacheManagerNodeCacheMap(), // make(map[string]ifc.BaseNodeInterface),
		subList:   NewCacheManagerSubListMap(),   // make(map[string]uint32)}
		cbsch:     scheduler.NewScheduler(),
		dmc:       dmc,
	}
	return e
}

func (cm *CacheManager) IsNodeInCache(n string) bool {
	_, ok := cm.nodeCache.Get(n)
	return ok
}

func (cm *CacheManager) IsNodePresentAndValid(n string) bool {
	nd, ok := cm.nodeCache.Get(n)
	if ok {
		return !nd.IsDeleted()
	}
	return false
}

func (cm *CacheManager) IsNodeSubscribed(n string) bool {
	nd, ok := cm.nodeCache.Get(n)
	if !ok {
		return false
	}
	path := nd.GetFullPath()
	if subNode := cm.getSubNode(path); subNode != nil {
		if subNode.GetWeight() != 0 {
			return true
		}
	}
	return cm.CheckPath(path)
}

func (cm *CacheManager) IsNodeDeleted(n string) bool {
	nd, ok := cm.nodeCache.Get(n)
	if ok {
		return nd.IsDeleted()
	}
	return false
}

func (cm *CacheManager) GetCachedNode(n string) (ifc.BaseNodeInterface, bool) {
	return cm.nodeCache.Get(n)
}

func (cm *CacheManager) WriteToCache(parentId, id string, base ifc.BaseNodeInterface) {
	logging.Debugf("%s:WriteToCache:: id = %s prop = %s", cm.name, id, base.GetImmediateBaseProperties())
	parentFound := false
	pnd, pndok := cm.nodeCache.Get(parentId)
	if parentId != "" && pndok {
		base.SetBaseParent(pnd)
		parentFound = true
	} else {

		// TODO: remove this once we're more stable and the templates have been cleaned up.
		npath := base.GetFullPath()
		if len(npath) > 1 {
			//We return since current node is expected to be deleted via recursive deletes in the parent chain
			logging.Warnf("%s: Trying to store node with id %s without parent", cm.name, base.GetId())
			return
		}
	}
	path := base.GetFullPath()
	spath := utils.JsonMarshal(path)
	logging.Debugf("%s==>Adding / Executing CB for %s with parent %s.%t path=%s",
		cm.name, id, parentId, parentFound, string(spath))
	cm.doNodeAddCallbacks(base, path)
	cm.nodeCache.Set(id, base)
	cm.updateTreeAdd(id, path, base)
	// cm.log.Call("debug", fmt.Sprintf("==>Adding to Cache %s with parent %s.%b path=%s",
	// 	id, parentId, parentFound, js.Global.Get("JSON").Call("stringify", path)))
}

func (cm *CacheManager) updateTreeAdd(nodeId string, path ifc.NodePathList, base ifc.BaseNodeInterface) {
	logging.Debugf("%s:updateTreeAdd:: nodeId %s path %s", cm.name, nodeId, path)
	root := path[0][ifc.NodePathName_nodeType]
	cur, ok := cm.subTree.Get(root)
	if !ok {
		cur = NewSubNodeType()
		cm.subTree.Set(root, cur)
	}
	for idx, val := range path {
		if idx != 0 {
			nType := val[ifc.NodePathName_nodeType]
			nValue := val[ifc.NodePathName_nodeIdentifier]
			logging.Debugf("%s:updateTreeAdd:       move to %s.%s", cm.name, nType, nValue)
			cur.Add(nType, nValue)
			cur = cur.GetChild(nType, nValue)
		}
	}
	logging.Debugf("%s:updateTreeAdd:       Adding to current node. %s", cm.name, nodeId)
	cur.AddCachedNode(nodeId, &CachedNodeData{NodeInCache: true, NodeRequested: false})
}

func (cm *CacheManager) updateTreeDelInternal(tree *SubNodeType, idx int, nodeId string, path ifc.NodePathList) bool {
	if tree == nil {
		return true
	}
	if idx >= len(path) {
		tree.DelCachedNode(nodeId)
	} else {
		nType := path[idx][ifc.NodePathName_nodeType]
		nValue := path[idx][ifc.NodePathName_nodeIdentifier]
		chnd := tree.GetChild(nType, nValue)
		if chnd != nil {
			if cm.updateTreeDelInternal(chnd, idx+1, nodeId, path) {
				tree.DelChild(nType, nValue)
			}
		}
	}
	return tree.IsEmpty()
}

func (cm *CacheManager) updateTreeDel(nodeId string, path ifc.NodePathList) {
	root := path[0][ifc.NodePathName_nodeType]
	rootSN, _ := cm.subTree.Get(root)
	cm.updateTreeDelInternal(rootSN, 1, nodeId, path)
}

func (cm *CacheManager) PurgeCacheTree(nodeId string) {
	node, nodeOk := cm.GetCachedNode(nodeId)
	if !nodeOk {
		return
	}
	if node.IsDeleteCompleted() {
		return
	}
	// remove all the children
	node.GetLinks().ForEach(func(ntype, childKey string, lnk *ifc.GLink) {
		if lnk.Properties[common.LinkFixedProp_HardLink] == "true" {
			// do link delete callback
			cm.PurgeCacheTree(lnk.DestinationNodeId)
		} else {
			// do a link delete callback for all child nodes when soft link
			// hard link will be handled by parent link delete notification
			node.DeleteImmediateLink(ntype, childKey)
		}
	})
	// parent link delete notification cb
	parent := node.GetBaseParent()
	if parent != nil {
		logging.Debugf("Deleting immediate link nodeId = %s parentId = %s\n", node.GetId(), parent.GetId())
		parent.DeleteImmediateLink(node.GetType(), node.GetKeyValue())
	}
	// delete all the soft link pointing to this node
	node.ReverseLinkIterate(func(linkKeyValue, nid, lid string) {
		logging.Debugf("Sushil1 %s Reverse Link Delete rev Nid = %s lid = %s linkKeyVale=%s\n", node.GetId(), nid, lid, linkKeyValue)
		if nsrc, nsrcOk := cm.GetCachedNode(nid); nsrcOk {
			logging.Debugf("Sushil2 %s Reverse Link Delete rev Nid = %s lid = %s \n", node.GetId(), nid, lid)
			if nsrc.GetLinks().Has(node.GetType(), linkKeyValue) {
				logging.Debugf("Sushil3 %s Reverse Link Delete rev Nid = %s lid = %s \n", node.GetId(), nid, lid)
				nsrc.DeleteImmediateLink(node.GetType(), linkKeyValue)
			}
		} else {
			return
		}
	})
	// remove the reverse link

	cm.purgeCache(nodeId)
}

func (cm *CacheManager) purgeCache(id string) {
	if nd, ndOk := cm.GetCachedNode(id); ndOk {
		path := nd.GetFullPath()
		cm.doNodeDeleteCallbacks(nd)
		cm.deleteNode(id)
		cm.updateTreeDel(id, path)
	}
	// cm.log.debug(`Deleting form Cache ${id}`);
}

func (cm *CacheManager) deleteNode(n string) {
	if dt, ok := cm.nodeCache.Get(n); ok {
		dt.SetDeleteCompleted()
		cm.nodeCache.Del(n)
		logging.Debugf("%s==>Deleting node from cache %s", cm.name, n)

	}
}
func indexOf(arr []string, dt string) int {
	for idx, v := range arr {
		if v == dt {
			return idx
		}
	}
	return -1
}

func (cm *CacheManager) AddSubscription(
	path ifc.NodePathList,
	depth uint32) []string {
	return cm.addSubscription(path, int(depth), nil, nil, false, "")
}

// add subscription for a path, with callback function
func (cm *CacheManager) AddNodeSubscription(
	path ifc.NodePathList,
	cbfn *ifc.CallbackFuncNode) []string {
	return cm.addSubscription(path, 0, nil, cbfn, false, "")
}

// add subscription for a path, with callback function
func (cm *CacheManager) AddLinkSubscription(
	path ifc.NodePathList,
	cbfn *ifc.CallbackFuncLink,
	destNodeType string) []string {
	return cm.addSubscription(path, 0, cbfn, nil, true, destNodeType)
}

func (cm *CacheManager) addSubscriptionProcessor(
	path ifc.NodePathList,
	curNode *SubNodeType,
	idx int,
	depth int,
	syncList map[string]bool,
	cbfnLink *ifc.CallbackFuncLink,
	cbfnNode *ifc.CallbackFuncNode,
	destNodeType string,
	matchOnly bool) int {

	nType := path[idx][ifc.NodePathName_nodeType]
	nValue := path[idx][ifc.NodePathName_nodeIdentifier]
	logging.Debugf("%s:: idx=%d/%d node=%s/%s matchOnly=%t",
		cm.name, idx, len(path), nType, nValue, matchOnly)
	// incChildWeight on nValue == '*' should return a sync list which has all
	//  the child nodes
	if !matchOnly {
		cn := curNode.IncChildWeight(nType, nValue)
		for _, cni := range cn {
			syncList[cni] = true
		}
	} else {
		// if there are any nodes present that are matching wildcard they need to be
		// synced to db. to get latest state post sync request. if the weight on this node is 0.
		if curNode.GetWeight() == 0 {
			cn := curNode.GetCachedNodeKeys()
			for _, cni := range cn {
				syncList[cni] = true
			}
		}
	}
	nextNode := curNode.GetChild(nType, nValue)
	if nextNode == nil && !matchOnly {
		logging.Debugf("%s:: no nextNode for %d, %s", cm.name, idx, path)
		return 0
	}
	if idx != len(path)-1 {
		if nValue == "*" {
			// have to checkout all child nodes for the current node
			allChildKey := curNode.GetChildKeyForType(nType)
			for _, ival := range allChildKey {
				if ival != "*" {
					ichild := curNode.GetChild(nType, ival)
					// do investigate on match only
					tSyncList := make(map[string]bool)
					logging.Debugf("%s:: Exploring Child path for node=%s/%s", cm.name, nType, ival)
					cm.addSubscriptionProcessor(path, ichild, idx+1, depth, tSyncList,
						cbfnLink, cbfnNode, destNodeType, true)
					logging.Debugf("%s:: Exploring Child path for node=%s/%s Got %v",
						cm.name, nType, ival, syncList)
					// For now sync every node that matches, this is expensive in a deep tree.
					// but without syncing it's not possible to know the current status on the path.
					for key := range tSyncList {
						syncList[key] = true
					}
				}
			}
		}
		if idx == len(path)-2 && nextNode != nil {
			// the parent node for all sub's need to be resynced ...
			cn := nextNode.GetCachedNodeKeys()
			for _, cni := range cn {
				syncList[cni] = true
			}
		}
		if nextNode == nil {
			return 0
		}
		// if not last iteration
		return cm.addSubscriptionProcessor(path, nextNode, idx+1, depth, syncList,
			cbfnLink, cbfnNode, destNodeType, matchOnly)
	}

	// last iteration
	cbCnt := 0
	if depth != 0 {
		if !matchOnly {
			nextNode.SetDepth(uint32(depth))
		}
		cbCnt++
	}
	if cbfnLink != nil || cbfnNode != nil {
		if !matchOnly {
			if cbfnLink != nil {
				nextNode.AddCBLink(cbfnLink, destNodeType)
			} else {
				nextNode.AddCBNode(cbfnNode)
			}
		}
		// if this  node has a cached object we need to do the callback on it
		cachedNodes := []string{}
		if nextNode != nil {
			cachedNodes = nextNode.GetCachedNodeKeys()
		}
		logging.Debugf("%s ---> Extract node %s/%s len=%d",
			cm.name, nType, nValue, len(cachedNodes))
		if nValue == "*" {
			curNode.ChildForEach(func(nTypeInt, nValueInt string, ndInt *SubNodeType) {
				logging.Debugf("%s ---> Extract * node for %s/%s len = %d",
					cm.name, nTypeInt, nValueInt, len(ndInt.GetCachedNodeKeys()))
				if nTypeInt == nType && nValueInt != "*" {
					cachedNodes = append(cachedNodes, ndInt.GetCachedNodeKeys()...)
				}
			})
		}
		if len(cachedNodes) != 0 {
			logging.Debugf("%s: cachedNodes %d %s", cm.name, len(cachedNodes), cachedNodes)
			if cbfnLink != nil && destNodeType != "" {
				for _, cItm := range cachedNodes {
					if base, baseOk := cm.GetCachedNode(cItm); baseOk {

						// TODO: once goroutine/locking is done at link/node level the goroutines in here can be removed.
						go func(cItm string, base ifc.BaseNodeInterface) {
							base.GetLinks().ForEachType(destNodeType, func(nt, nv string, lnk *ifc.GLink) {
								logging.Debugf("%s:LINK CBFN Called from subscribe ===> %s", cm.name, cItm)
								cm.dmc.ReadyDestNodeBeforeCB(base, nv, destNodeType, lnk.Properties)
								(*cbfnLink)(base, common.UpdateType_LinkAdd, nv, ifc.PropertyType{}, lnk.Properties)
							})
						}(cItm, base)
						cbCnt++
					}
				}
			} else {
				for _, cItm := range cachedNodes {
					if base, baseOk := cm.GetCachedNode(cItm); baseOk {

						// TODO: once goroutine/locking is done at link/node level the goroutines in here can be removed.
						go func(cItm string, base ifc.BaseNodeInterface) {
							logging.Debugf("%s:Node CBFN Called from subscribe ===> %d", cm.name, len(cachedNodes))
							(*cbfnNode)(base, common.UpdateType_NodeAdd, ifc.PropertyType{}, base.GetImmediateBaseProperties())
						}(cItm, base)
						cbCnt++
					}
				}
			}
		}
	}
	return cbCnt
}

func (cm *CacheManager) addSubscription(
	path ifc.NodePathList,
	depth int,
	cbfnLink *ifc.CallbackFuncLink,
	cbfnNode *ifc.CallbackFuncNode,
	isLinkCB bool,
	destNodeType string) []string {
	pathjson := utils.JsonMarshal(path)

	spath := string(pathjson)
	logging.Debugf("%s, Add Sub called with PATH=%s, depth:%d linkcb=%t %s", cm.name, spath, depth, isLinkCB, destNodeType)
	if depth != 0 {
		if _, ok := cm.subList.Get(spath); ok {
			logging.Errorf(cm.name + ":This Path " + spath + "is already already subscribed.")
			return []string{}
		}
		cm.subList.Set(spath, uint32(depth))
	}
	syncList := make(map[string]bool)
	nType0 := path[0][ifc.NodePathName_nodeType]
	curNodeType, curNodeTypeOk := cm.subTree.Get(nType0)
	if !curNodeTypeOk {
		curNodeType = NewSubNodeType()
		cm.subTree.Set(nType0, curNodeType)
		curNodeType.InitWeight()
	} else {
		// syncList = syncList.concat(this.subTree[nType0].incWeight());
		nodeKeys := curNodeType.IncWeight()
		for _, key := range nodeKeys {
			syncList[key] = true
		}
	}
	cm.addSubscriptionProcessor(path, curNodeType, 1, depth, syncList, cbfnLink, cbfnNode, destNodeType, false)
	// cm.printSubTree();
	logging.Debugf("%s:, Sync Init for %v", cm.name, syncList)
	syncListKeys := make([]string, 0, len(syncList))
	for k := range syncList {
		syncListKeys = append(syncListKeys, k)
	}
	return syncListKeys
}

func (cm *CacheManager) getNodeCBforPathInternal(
	path ifc.NodePathList, curNode *SubNodeType, idx int,
	fn []*ifc.CallbackFuncNode) []*ifc.CallbackFuncNode {
	if curNode == nil {
		return fn
	}
	if len(path) == idx {
		curNode.CbfnNodeForEach(func(v *ifc.CallbackFuncNode) {
			fn = append(fn, v)
		})
		return fn
	}
	nType := path[idx][ifc.NodePathName_nodeType]
	nValue := path[idx][ifc.NodePathName_nodeIdentifier]
	if !curNode.IsTypeInChild(nType) {
		return fn
	}
	rdt := curNode.GetChild(nType, nValue)
	if rdt != nil {
		fn = cm.getNodeCBforPathInternal(path, rdt, idx+1, fn)
	}
	rdtw := curNode.GetChild(nType, "*")
	if rdtw != nil {
		fn = cm.getNodeCBforPathInternal(path, rdtw, idx+1, fn)
	}
	return fn
}

func (cm *CacheManager) getNodeCBforPath(path ifc.NodePathList) []*ifc.CallbackFuncNode {
	var ret []*ifc.CallbackFuncNode = []*ifc.CallbackFuncNode{}
	if len(path) != 0 {
		v, ok := cm.subTree.Get(path[0][ifc.NodePathName_nodeType])
		if ok {
			ret = cm.getNodeCBforPathInternal(path, v, 1, ret)
		}
	}
	return ret
}

func (cm *CacheManager) getLinkCBforPathInternal(
	path ifc.NodePathList,
	curNode *SubNodeType,
	idx int,
	nodeType string,
	fn []*ifc.CallbackFuncLink) []*ifc.CallbackFuncLink {

	if curNode == nil {
		return fn
	}
	if len(path) == idx {
		curNode.CbfnLinkForEach(nodeType, func(cb *ifc.CallbackFuncLink) {
			fn = append(fn, cb)
		})
		return fn
	}
	nType := path[idx][ifc.NodePathName_nodeType]
	nValue := path[idx][ifc.NodePathName_nodeIdentifier]
	if !curNode.IsTypeInChild(nType) {
		return fn
	}
	rdt := curNode.GetChild(nType, nValue)
	if rdt != nil {
		fn = cm.getLinkCBforPathInternal(path, rdt, idx+1, nodeType, fn)
	}
	rdtw := curNode.GetChild(nType, "*")
	if rdtw != nil {
		//     // wildcard match
		fn = cm.getLinkCBforPathInternal(path, rdtw, idx+1, nodeType, fn)
	}
	return fn
}

func (cm *CacheManager) getLinkCBforPath(
	path ifc.NodePathList, nType string) []*ifc.CallbackFuncLink {
	var ret []*ifc.CallbackFuncLink = []*ifc.CallbackFuncLink{}
	if len(path) != 0 {
		v, ok := cm.subTree.Get(path[0][ifc.NodePathName_nodeType])
		if ok {
			ret = cm.getLinkCBforPathInternal(path, v, 1, nType, ret)
		}
	}
	if len(ret) != 0 {
		bpath := utils.JsonMarshal(path)

		logging.Debugf("%s, LinkCB for path=%s nType=%s returning cbfn %d",
			cm.name, string(bpath), nType, len(ret))
	}
	return ret
}

func (cm *CacheManager) getSubNode(path ifc.NodePathList) *SubNodeType {
	logging.Debugf("%s :: %s", cm.name, path)
	if len(path) == 0 {
		return nil
	}
	var curNodeType *SubNodeType
	for idx, pathv := range path {
		nType := pathv[ifc.NodePathName_nodeType]
		nValue := pathv[ifc.NodePathName_nodeIdentifier]
		if idx == 0 {
			var ok bool
			curNodeType, ok = cm.subTree.Get(nType)
			if !ok {
				break
			}
			if len(path) == 1 {
				return curNodeType
			}
			continue
		}
		rdt := curNodeType.GetChild(nType, nValue)
		if rdt != nil {
			break
		}
		if idx == len(path)-1 {
			// last iteration
			return rdt
		}
		curNodeType = rdt
	}
	return nil
}

func (cm *CacheManager) DelSubscription(path ifc.NodePathList) {
	spathBytes := utils.JsonMarshal(path)

	spath := string(spathBytes)
	_, spathOk := cm.subList.Get(spath)
	if !spathOk {
		return
	}
	cm.subList.Del(spath)

	var curNodeType *SubNodeType
	for idx, pathv := range path {
		nType := pathv[ifc.NodePathName_nodeType]
		nValue := pathv[ifc.NodePathName_nodeIdentifier]
		if idx == 0 {
			subTreeV, subTreeOk := cm.subTree.Get(nType)
			if !subTreeOk {
				panic(errors.New(cm.name + ":Error while navigating cache datastructure when unsubscribing 1."))
			} else {
				subTreeV.DecWeight()
				curNodeType = subTreeV
			}
			continue
		}
		if !curNodeType.IsTypeInChild(nType) {
			panic(errors.New(cm.name + ":Error while navigating cache datastructure when unsubscribing 2."))
		}
		rdt := curNodeType.GetChild(nType, nValue)
		if rdt == nil {
			panic(errors.New(cm.name + ":Error while navigating cache datastructure when unsubscribing 3."))
		} else {
			rdt.DecWeight()
		}
		if idx == len(path)-1 {
			// last iteration
			rdt.SetDepth(0)
		}
		curNodeType = rdt
	}
}

func (cm *CacheManager) checkPathInternal(path ifc.NodePathList,
	idx int, curNode *SubNodeType) bool {
	if curNode == nil {
		return false
	}
	lleft := len(path) - idx
	if idx != 1 && int(curNode.GetDepth()) > lleft {
		// found a valid notificaiton
		return true
	}
	if curNode.GetWeight() != 0 && idx == len(path) {
		return true
	}
	if len(path) == idx {
		return false
	} else {
		nType := path[idx][ifc.NodePathName_nodeType]
		nValue := path[idx][ifc.NodePathName_nodeIdentifier]
		snt := curNode.GetChild(nType, nValue)
		sntw := curNode.GetChild(nType, "*")
		if snt != nil && cm.checkPathInternal(path, idx+1, snt) {
			return true
		}
		if sntw != nil && cm.checkPathInternal(path, idx+1, sntw) {
			return true
		}
		return false
	}
}

func (cm *CacheManager) CheckPath(path ifc.NodePathList) bool {
	if cm.fetchAll {
		return true
	}
	if len(path) == 0 {
		return false
	}
	nt := path[0][ifc.NodePathName_nodeType]
	found := false
	stv, stok := cm.subTree.Get(nt)
	if stok {
		if len(path) == 1 {
			found = stv.GetWeight() != 0 || stv.GetDepth() >= 1
		} else {
			cmsn, _ := cm.subTree.Get(nt)
			found = cm.checkPathInternal(path, 1, cmsn)
		}
	}
	logging.Debugf("%s: path=%v found=%t", cm.name, path, found)
	return found
}

func (cm *CacheManager) checkPropDiff(p1 ifc.PropertyType,
	p2 ifc.PropertyType,
	exclude []string) bool {
	if len(exclude) == 0 && len(p1) != len(p2) {
		return true
	}

	/* Do not generate an update event for the node if the update is specific to
	   ToBeDeleted property. Ignore changeId and updateTime update changes in this case
	*/
	ignoreChangeId := false
	ignoreUpdateTime := false
	if _, ok := p2[common.NodeFixedProp_ToBeDeleted]; ok {
		if _, ok1 := p1[common.NodeFixedProp_ToBeDeleted]; !ok1 {
			ignoreChangeId = true
			ignoreUpdateTime = true
			logging.Debugf("checkPropDiff: exclusion propset: %v \noldpropset = %v \n newpropset = %v\n", exclude, p1, p2)
		}
	}

	for key, kval := range p1 {
		if (key == common.NodeFixedProp_updateTime && ignoreUpdateTime == true) ||
			(key == common.NodeFixedProp_changeId && ignoreChangeId == true) {
			continue
		}
		if indexOf(exclude, key) == -1 {
			if p2v, p2vOk := p2[key]; p2vOk {

				// The data model sometimes emulates enums via string/iota constants.
				// While this works fine when everything is in memory we lose this typing when
				// those properties go to the database and are later read back. It's possible for
				// a user-defined type to reduce back to an int or even a float64. We can catch and
				// fix this here by comparing the types of the two properties so that we can restore
				// the original (user-defined) type that was used.

				// Note that without this detection any properties that use user-defined types (enums)
				// will incorrectly trigger update events from the data model which increases load and
				// can sometimes ruin application logic.
				t1 := reflect.TypeOf(kval)
				t2 := reflect.TypeOf(p2v)
				if t1 != t2 {
					p2v = reflect.ValueOf(p2v).Convert(t1).Interface()
				}

				if !reflect.DeepEqual(kval, p2v) {
					logging.Debugf("checkPropDiff: Return true for p1 key %s\n", key)
					return true
				}
			}
		}
	}
	for key := range p2 {
		if indexOf(exclude, key) == -1 {
			_, p1ok := p1[key]
			if !p1ok {
				logging.Debugf("checkPropDiff: Return true for p2 key %s\n", key)
				return true
			}
		}
	}
	logging.Debugf("checkPropDiff: Return false: No diff between old and new prop\n")
	return false
}

// TODO: caller needs to synchronize on base node ID (path) before calling this.
func (cm *CacheManager) doNodeAddCallbacks(base ifc.BaseNodeInterface, path ifc.NodePathList) {
	cblist := cm.getNodeCBforPath(path)
	for _, cb := range cblist {
		spath := utils.JsonMarshal(path)

		logging.Debugf("%s CBFN Called from subscribe path=%s", cm.name, string(spath))
		go func(cb *ifc.CallbackFuncNode, base ifc.BaseNodeInterface) {
			(*cb)(base, common.UpdateType_NodeAdd, nil, base.GetImmediateBaseProperties())
		}(cb, base)
	}
}

// TODO: caller needs to synchronize on base node ID (path) before calling this.
func (cm *CacheManager) doNodeDeleteCallbacks(base ifc.BaseNodeInterface) {
	path := base.GetFullPath()
	cblist := cm.getNodeCBforPath(path)
	for _, cb := range cblist {
		spath := utils.JsonMarshal(path)

		logging.Debugf("%s Delete Called from subscribe path=%s", cm.name, string(spath))
		go func(cb *ifc.CallbackFuncNode, base ifc.BaseNodeInterface) {
			(*cb)(base, common.UpdateType_NodeDelete, base.GetImmediateBaseProperties(), nil)
		}(cb, base)
	}
}

// TODO: caller needs to synchronize on base node ID (path) before calling this.
func (cm *CacheManager) DoNodeUpdateCallbacks(
	base ifc.BaseNodeInterface,
	curProp,
	newProp ifc.PropertyType) {
	path := base.GetFullPath()
	base.Lock()
	if !cm.checkPropDiff(curProp, newProp, []string{
		common.NodeFixedProp_creationTime, common.NodeFixedProp_createdBy,
		common.NodeFixedProp_Revision, common.NodeFixedProp_ToBeDeleted}) {
		base.Unlock()
		return
	}

	base.Unlock()
	cblist := cm.getNodeCBforPath(path)
	for _, cb := range cblist {
		spath := utils.JsonMarshal(path)
		logging.Debugf("%s Called from subscribe path=%s nd=%s", cm.name, string(spath), newProp)
		go func(cb *ifc.CallbackFuncNode, base ifc.BaseNodeInterface) {
			(*cb)(base, common.UpdateType_NodeUpdate, curProp, newProp)
		}(cb, base)
	}
}

// TODO: caller needs to synchronize on base node ID (path) before calling this.
func (cm *CacheManager) DoLinkDeleteCallbacks(
	base ifc.BaseNodeInterface,
	destType, dstKey string,
	oldLink ifc.PropertyType) {
	srcPath := base.GetFullPath()
	cblist := cm.getLinkCBforPath(srcPath, destType)
	spath := utils.JsonMarshal(srcPath)
	logging.Debugf("%s:Doing a link delete cb for %s destType=%s.%s cbsize=%d %s",
		cm.name, string(spath), destType, dstKey, len(cblist), base.GetKeyValue())
	for _, cb := range cblist {
		go func(cb *ifc.CallbackFuncLink, base ifc.BaseNodeInterface) {
			(*cb)(base, common.UpdateType_LinkDelete, dstKey, oldLink, nil)
		}(cb, base)
	}
}

// TODO: caller needs to synchronize on base node ID (path) before calling this.
func (cm *CacheManager) DoLinkUpdateCallbacks(
	base ifc.BaseNodeInterface,
	dstKey, destType string,
	curLink, newLink *ifc.GLink) {
	destKeyId := base.GetId() + "/" + destType + "/" + dstKey
	/*
		NOTE (TODO): The callback funcitions need to be executed in order that they got triggered
		The call to "ReadyDestNodebeforeCB" can be blocking and also there is no predictability in the go routines in case they get backlogged the order will be changes
		Adding a lock here will serialize the go routines but this function is changed to blocking from non-blocking. This will affect the callers. All callers need to be checked for this transition.
		Also the node callback funcition needs to be evaluated as well for similar change.
		All use of go routines needs a through check to make sure the ordering does not cause problem.
	*/
	lck := cm.cbsch.Wait(destKeyId)
	defer cm.cbsch.Done(lck)
	srcPath := base.GetFullPath()
	spath := utils.JsonMarshal(srcPath)
	logging.Debugf("%s:doing a link update for %s destType=%s.%s curLink=%p",
		cm.name, string(spath), destType, dstKey, curLink)

	// TODO: remove this once we're more stable and the templates have been cleaned up.
	{
		bpath := base.GetFullPath()
		if len(bpath) > 1 && base.GetBaseParent() == nil {
			logging.Warnf("%s : No Parent set during callback processing\n", base.GetId())
			return
		}
	}
	if curLink != nil {
		diff := cm.checkPropDiff(
			curLink.Properties, newLink.Properties,
			[]string{common.LinkFixedProp_createdBy, common.LinkFixedProp_creationTime, common.LinkFixedProp_Revision})
		if diff {
			cblist := cm.getLinkCBforPath(srcPath, destType)
			if len(cblist) > 0 {
				logging.Debugf("%s:Link update callback[%s] cblist length is %d key = %s ut=%d", cm.name, base.GetId(), len(cblist), dstKey, common.UpdateType_LinkUpdate)
				cm.dmc.ReadyDestNodeBeforeCB(base, dstKey, destType, newLink.Properties)
			}
			for _, cb := range cblist {
				// go func(cb *ifc.CallbackFuncLink, base ifc.BaseNodeInterface) {
				(*cb)(base, common.UpdateType_LinkUpdate, dstKey, curLink.Properties, newLink.Properties)
				// }(cb, base)
			}
		} else {
			logging.Debugf("%s: Link update callback skipped no diff found on link properties", cm.name)
		}
	} else {
		cblist := cm.getLinkCBforPath(srcPath, destType)
		if len(cblist) > 0 {
			logging.Debugf("%s:link update callback2[%s] cblist length is %d key = %s ut=%d", cm.name, base.GetId(), len(cblist), dstKey, common.UpdateType_LinkAdd)
			cm.dmc.ReadyDestNodeBeforeCB(base, dstKey, destType, newLink.Properties)
		}
		for _, cb := range cblist {
			// go func(cb *ifc.CallbackFuncLink, base ifc.BaseNodeInterface) {
			(*cb)(base, common.UpdateType_LinkAdd, dstKey, nil, newLink.Properties)
			//}(cb, base)
		}
	}
}
