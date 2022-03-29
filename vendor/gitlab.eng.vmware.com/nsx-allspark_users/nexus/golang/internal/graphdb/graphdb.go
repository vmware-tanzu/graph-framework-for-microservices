package graphdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"sync/atomic"

	"strings"
	"time"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/internal/utils"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/ifc"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/logging"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/scheduler"

	"github.com/lithammer/shortuuid"
	"go.etcd.io/etcd/clientv3"
)

type GraphDB struct {
	client            *clientv3.Client
	stats             common.GraphDBStats
	scheduler         *scheduler.Scheduler
	fixedProp         map[string]bool
	lockRetryDelay    int
	txnRetryDelay     int
	maxTxnRetryCount  int
	maxLockRetryCount int
	retryDelay        int
	maxRetryCount     int
	name              string
	featureFlags      []string
}

const (
	FixedPath_Created = "_created"
	FixedPath_Links   = "_links"
	FixedPath_RLinks  = "_rlinks"
	FixedPath_Lock    = "_lock"
)

func New(name, etcdLoc string, dmFeatureFlags []string) *GraphDB {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdLoc}, // "http://127.0.0.1:2379"},
		DialTimeout: 10 * time.Second,
	})
	if err == context.DeadlineExceeded {
		log.Fatal(err)
	} else if err != nil {
		log.Fatal(err)
	}
	g := &GraphDB{
		lockRetryDelay:    500,
		txnRetryDelay:     100,
		maxLockRetryCount: 100,
		maxTxnRetryCount:  100,
		retryDelay:        10,
		maxRetryCount:     32,
		name:              name,
		featureFlags:      dmFeatureFlags,
		stats:             common.GraphDBStats{},
		scheduler:         scheduler.NewScheduler(),
		fixedProp:         make(map[string]bool),
		client:            cli}
	g.fixedProp[common.NodeFixedProp_NodeKeyName] = true
	g.fixedProp[common.NodeFixedProp_IsRoot] = true
	g.fixedProp[common.NodeFixedProp_NodeDefaultKeyName] = true
	g.fixedProp[common.NodeFixedProp_NodeSingletonKeyValue] = true
	g.fixedProp[common.NodeFixedProp_Revision] = true
	g.fixedProp[common.NodeFixedProp_createdBy] = true
	g.fixedProp[common.NodeFixedProp_updatedBy] = true
	g.fixedProp[common.NodeFixedProp_creationTime] = true
	g.fixedProp[common.NodeFixedProp_updateTime] = true
	g.fixedProp[common.NodeFixedProp_changeId] = true
	g.fixedProp[common.LinkFixedProp_HardLink] = true
	g.fixedProp[common.LinkFixedProp_NodeKeyName] = true
	g.fixedProp[common.LinkFixedProp_NodeKeyValue] = true
	g.fixedProp[common.LinkFixedProp_NodeType] = true
	g.fixedProp[common.LinkFixedProp_SoftLinkDestinationPath] = true
	g.fixedProp[common.LinkFixedProp_Revision] = true
	g.fixedProp[common.LinkFixedProp_createdBy] = true
	g.fixedProp[common.LinkFixedProp_updatedBy] = true
	g.fixedProp[common.LinkFixedProp_creationTime] = true
	g.fixedProp[common.LinkFixedProp_updateTime] = true
	logging.Debugf("CREATE: Creating new GraphDB[%+v] %s Location = %s\n", *g, name, etcdLoc)

	return g
}

func (g *GraphDB) isRLinkFeatureFlagEnabled() bool {
	for _, v := range g.featureFlags {
		if strings.Compare(v, common.Rlink_EnableFeatureFlag) == 0 {
			return true
		}
	}
	return false
}

func (g *GraphDB) Shutdown() {
	g.client.Close()
	logging.Infof("SHUTDOWN: Stopping GraphDB %s[%p] Done\n", g.name, g)
}
func (g *GraphDB) GetClient() *clientv3.Client {
	return g.client
}
func (g *GraphDB) retryPut(fn func() (*clientv3.PutResponse, error)) *clientv3.PutResponse {
	r := g.retry_(func() (interface{}, error) {
		return fn()
	})
	return r.(*clientv3.PutResponse)
}
func (g *GraphDB) retryGet(fn func() (*clientv3.GetResponse, error)) *clientv3.GetResponse {
	r := g.retry_(func() (interface{}, error) {
		return fn()
	})
	return r.(*clientv3.GetResponse)
}
func (g *GraphDB) retryTxn(fn func() (*clientv3.TxnResponse, error)) *clientv3.TxnResponse {
	r := g.retry_(func() (interface{}, error) {
		return fn()
	})
	return r.(*clientv3.TxnResponse)
}
func (g *GraphDB) GetLatestRevision() int64 {
	nodeId := "/Root/root"
	curPropS := g.retryGet(func() (*clientv3.GetResponse, error) {
		return g.client.KV.Get(context.TODO(), nodeId)
	})
	atomic.AddUint32(&g.stats.DbRead, 1)
	if curPropS == nil {
		logging.Warnf("Unable to get revision for %s node. Get op returned empty object", nodeId)
		fmt.Printf("Unable to get revision for %s node. Get op returned empty object", nodeId)
		return -1
	}
	if len(curPropS.Kvs) == 0 {
		// may be node was deleted
		logging.Warnf("Unable to get revision for %s node. Node not found", nodeId)
		fmt.Printf("Unable to get revision for %s node. Node not found", nodeId)
		return -1
	}
	return curPropS.Header.Revision
}
func (g *GraphDB) CompactRevision(revision, compactCtxtTimeout int64, compactPhysical bool) (*ifc.EtcdCompactResponse, error) {
	var opts []clientv3.CompactOption
	if compactPhysical == true {
		opts = append(opts, clientv3.WithCompactPhysical())
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(compactCtxtTimeout)*time.Second)
	resp, err := g.client.KV.Compact(ctx, revision, opts...)
	cancel()
	if err != nil {
		logging.Errorf("Error in compacting revision %d. Error: %s", revision, err.Error())
		errNew := fmt.Errorf("Error in compacting revision %d. Error: %s", revision, err.Error())
		return nil, errNew
	}
	if resp == nil || resp.Header == nil {
		logging.Errorf("Error in compacting revision %d. Txn response is empty.", revision)
		errNew := fmt.Errorf("Error in compacting revision %d. Txn response is empty.", revision)
		return nil, errNew
	}

	compactResp := &ifc.EtcdCompactResponse{ClusterId: resp.Header.ClusterId,
		MemberId: resp.Header.MemberId,
		RaftTerm: resp.Header.RaftTerm,
		Revision: resp.Header.Revision}
	return compactResp, nil
}
func (g *GraphDB) retry_(fn func() (interface{}, error)) interface{} {
	cnt := 0
	backoff := g.retryDelay
	for {
		ret, err := fn()
		if err == nil {
			return ret
		}
		logging.Warnf("Warning (Will Retry %d) while executing Txn on etcd %s", cnt, err)
		debug.PrintStack()
		cnt++
		time.Sleep(time.Duration(backoff) * time.Millisecond)
		if backoff < 5000 {
			backoff = backoff * 2
		}
		if cnt > g.maxRetryCount {
			logging.Fatalf("ERROR while trying to execute operation (retrycnt = %d) on etcd %s", cnt, err)
		}
	}
}
func (g *GraphDB) checkId(nodeId string) {
	if nodeId == "" {
		logging.Fatalf("cant use empty nodeId")
	}
}
func (g *GraphDB) copyPropWithoutFixed(propIn ifc.PropertyType) ifc.PropertyType {
	propOut := make(map[string]interface{})
	for k, v := range propIn {
		if _, ok := g.fixedProp[k]; !ok {
			propOut[k] = v
		}
	}
	return propOut
}

func (g *GraphDB) getLinks(path string) []*ifc.GLink {
	ret := []*ifc.GLink{}
	// get all the hard links
	cpath := path + "/" + FixedPath_Links
	resp := g.retryGet(func() (*clientv3.GetResponse, error) {
		return g.client.KV.Get(context.TODO(), cpath, clientv3.WithPrefix())
	})
	atomic.AddUint32(&g.stats.DbRead, 1)
	for _, items := range resp.Kvs {
		itemkey := string(items.Key)
		subStr := strings.Split(strings.Replace(itemkey, cpath, "", 1), "/")
		leafItem := subStr[len(subStr)-1]
		if leafItem == FixedPath_Created || leafItem == FixedPath_Lock {
			continue
		}
		leafType := subStr[1]
		node := g.retryGet(func() (*clientv3.GetResponse, error) {
			return g.client.KV.Get(context.TODO(), itemkey)
		})
		cnode := g.retryGet(func() (*clientv3.GetResponse, error) {
			return g.client.KV.Get(context.TODO(), itemkey+"/"+FixedPath_Created)
		})
		atomic.AddUint32(&g.stats.DbRead, 2)
		if len(node.Kvs) == 0 || len(cnode.Kvs) == 0 {
			continue
		}
		nodeDt := node.Kvs[0].Value
		cnodeDt := cnode.Kvs[0].Value
		logging.Debugf("==> %s;; leafType = %s   nodeDt = %s cnodeDt=%s", itemkey, leafType, string(nodeDt), string(cnodeDt))
		var nodeProp ifc.PropertyType
		var cnodeProp ifc.PropertyType
		nodePropErr := json.Unmarshal(nodeDt, &nodeProp)
		cnodePropErr := json.Unmarshal(cnodeDt, &cnodeProp)
		if nodePropErr != nil {
			panic(errors.New(fmt.Sprintf("Error when parsing prop from db path=%s data %s err %s", itemkey, string(nodeDt), nodePropErr)))
		}
		if cnodePropErr != nil {
			panic(errors.New(fmt.Sprintf("Error when parsing create prop from db path=%s data %s err %s", itemkey+"/"+FixedPath_Created, string(cnodeDt), cnodePropErr)))
		}
		lnk := &ifc.GLink{
			Id:                itemkey,
			LinkType:          leafType,
			Properties:        nodeProp,
			SourceNodeId:      path,
			DestinationNodeId: nodeProp[common.LinkFixedProp_destNodeId].(string),
		}
		for k, v := range cnodeProp {
			lnk.Properties[k] = v
		}
		lnk.Properties[common.LinkFixedProp_Revision] = node.Kvs[0].ModRevision
		ret = append(ret, lnk)
	}
	return ret
}

func (g *GraphDB) getRLinks(path string) []*ifc.GLink {
	ret := []*ifc.GLink{}
	if !g.isRLinkFeatureFlagEnabled() {
		return ret
	}
	cpath := path + "/" + FixedPath_RLinks
	resp := g.retryGet(func() (*clientv3.GetResponse, error) {
		return g.client.KV.Get(context.TODO(), cpath, clientv3.WithPrefix())
	})
	atomic.AddUint32(&g.stats.DbRead, 1)
	for _, items := range resp.Kvs {
		itemkey := string(items.Key)
		subStr := strings.Split(strings.Replace(itemkey, cpath, "", 1), "/")
		leafItem := subStr[len(subStr)-1]
		if leafItem == FixedPath_Created || leafItem == FixedPath_Lock {
			continue
		}
		leafType := subStr[1]
		node := g.retryGet(func() (*clientv3.GetResponse, error) {
			return g.client.KV.Get(context.TODO(), itemkey)
		})
		cnode := g.retryGet(func() (*clientv3.GetResponse, error) {
			return g.client.KV.Get(context.TODO(), itemkey+"/"+FixedPath_Created)
		})
		atomic.AddUint32(&g.stats.DbRead, 2)
		if len(node.Kvs) == 0 || len(cnode.Kvs) == 0 {
			continue
		}
		nodeDt := node.Kvs[0].Value
		cnodeDt := cnode.Kvs[0].Value
		logging.Debugf("==> %s;; leafType = %s   nodeDt = %s cnodeDt=%s", itemkey, leafType,
			string(nodeDt), string(cnodeDt))
		var nodeProp ifc.PropertyType
		var cnodeProp ifc.PropertyType
		nodePropErr := json.Unmarshal(nodeDt, &nodeProp)
		cnodePropErr := json.Unmarshal(cnodeDt, &cnodeProp)
		if nodePropErr != nil {
			panic(fmt.Errorf(fmt.Sprintf("Error when parsing prop from db path=%s data %s err %s",
				itemkey, string(nodeDt), nodePropErr)))
		}
		if cnodePropErr != nil {
			panic(fmt.Errorf(fmt.Sprintf("Error when parsing create prop from db path=%s data %s err %s",
				itemkey+"/"+FixedPath_Created, string(cnodeDt), cnodePropErr)))
		}
		lnk := &ifc.GLink{
			Id:                itemkey,
			LinkType:          leafType,
			Properties:        nodeProp,
			SourceNodeId:      path,
			DestinationNodeId: nodeProp[common.LinkFixedProp_destNodeId].(string),
		}
		for k, v := range cnodeProp {
			lnk.Properties[k] = v
		}
		lnk.Properties[common.LinkFixedProp_Revision] = node.Kvs[0].ModRevision
		ret = append(ret, lnk)
	}
	return ret
}

func (g *GraphDB) UpsertNode(requestorId, nodeType, nodeName string,
	nodePropIn ifc.PropertyType) *ifc.GNode {
	logging.Debugf("upsertNode: %s.%s prop = %s", nodeType, nodeName, nodePropIn)
	nodeProp := g.copyPropWithoutFixed(nodePropIn)
	nodeProp[common.NodeFixedProp_NodeDefaultKeyName] = nodeName
	now := time.Now().UTC().Format(utils.ISOTimeFormat)
	changeId := shortuuid.New()
	nodeProp[common.NodeFixedProp_IsRoot] = "true"
	nodeProp[common.NodeFixedProp_NodeKeyName] = common.NodeFixedProp_NodeDefaultKeyName
	nodeProp[common.NodeFixedProp_changeId] = changeId

	nodeProp[common.NodeFixedProp_updatedBy] = requestorId
	nodeProp[common.NodeFixedProp_updateTime] = now
	prop := string(utils.JsonMarshal(nodeProp))

	cpropdt := make(map[string]string)
	cpropdt[common.NodeFixedProp_createdBy] = requestorId
	cpropdt[common.NodeFixedProp_creationTime] = now
	createProp := string(utils.JsonMarshal(cpropdt))

	// path is /<nodeType>/<defaultkey>
	path := "/" + nodeType + "/" + nodeName
	cpath := path + "/" + FixedPath_Created
	resp := g.retryTxn(func() (*clientv3.TxnResponse, error) {
		return g.client.KV.Txn(context.TODO()).If(
			clientv3.Compare(clientv3.CreateRevision(cpath), "=", 0),
		).Then(
			clientv3.OpPut(cpath, string(createProp)),
			clientv3.OpPut(path, string(prop)),
		).Else(
			clientv3.OpPut(path, string(prop)),
		).Commit()
	})
	atomic.AddUint32(&g.stats.DbWrite, 1)
	ret := &ifc.GNode{
		Id:         path,
		Type:       nodeType,
		Properties: nodeProp,
		Links:      g.getLinks(path),
		RLinks:     g.getRLinks(path),
	}
	ret.Properties[common.NodeFixedProp_Revision] = resp.Header.Revision
	return ret
}

func (g *GraphDB) UpsertChildNode(
	requestorId, parentNodeId, linkType, linkKey string,
	linkPropIn ifc.PropertyType,
	nodeType string,
	nodePropIn ifc.PropertyType) (*ifc.GNode, *ifc.GLink) {
	g.checkId(parentNodeId)
	lck := g.scheduler.Wait(parentNodeId)
	defer g.scheduler.Done(lck)
	key := common.NodeFixedProp_NodeDefaultKeyName
	if linkKey != "" {
		key = linkKey
	}
	nodeProp := g.copyPropWithoutFixed(nodePropIn)
	linkProp := g.copyPropWithoutFixed(linkPropIn)

	now := time.Now().UTC().Format(utils.ISOTimeFormat)
	changeId := shortuuid.New()
	nodeName := nodePropIn[key].(string)
	if strings.Contains(nodeName, "/") {
		panic(fmt.Sprintf("Attempting database write with node name %s, slashes not allowed!", nodeName))
	}

	logging.Debugf("upserChildNode %s :: %s.%s P=%s", parentNodeId, nodeType, nodeName, nodeProp)
	ndpath := parentNodeId + "/" + nodeType + "/" + nodeName
	ndpathC := ndpath + "/" + FixedPath_Created
	ndprop := nodeProp
	ndprop[common.NodeFixedProp_IsRoot] = "false"
	ndprop[common.NodeFixedProp_NodeKeyName] = key
	ndprop[key] = nodeName
	ndprop[common.NodeFixedProp_updatedBy] = requestorId
	ndprop[common.NodeFixedProp_updateTime] = now
	ndprop[common.NodeFixedProp_changeId] = changeId

	ndpropC := make(map[string]string)
	ndpropC[common.NodeFixedProp_createdBy] = requestorId
	ndpropC[common.NodeFixedProp_creationTime] = now

	lnprop := linkProp
	lnpath := parentNodeId + "/" + FixedPath_Links + "/" + nodeType + "/" + nodeName
	lnpathC := lnpath + "/" + FixedPath_Created
	lnprop[common.LinkFixedProp_NodeKeyName] = key
	lnprop[common.LinkFixedProp_NodeKeyValue] = nodeName
	lnprop[common.LinkFixedProp_NodeType] = nodeType
	lnprop[common.LinkFixedProp_HardLink] = "true"
	lnprop[common.LinkFixedProp_updatedBy] = requestorId
	lnprop[common.LinkFixedProp_updateTime] = now
	lnprop[common.LinkFixedProp_destNodeId] = ndpath
	lnpropC := make(map[string]string)
	lnpropC[common.LinkFixedProp_createdBy] = requestorId
	lnpropC[common.LinkFixedProp_creationTime] = now
	ndpropStr := string(utils.JsonMarshal(ndprop))
	lnpropStr := string(utils.JsonMarshal(lnprop))
	ndpropCStr := string(utils.JsonMarshal(ndpropC))
	lnpropCStr := string(utils.JsonMarshal(lnpropC))

	resp := g.retryTxn(func() (*clientv3.TxnResponse, error) {
		return g.client.KV.Txn(context.TODO()).If(
			clientv3.Compare(clientv3.CreateRevision(ndpathC), "=", 0),
		).Then(
			clientv3.OpPut(ndpathC, ndpropCStr),
			clientv3.OpPut(lnpathC, lnpropCStr),
			clientv3.OpPut(ndpath, ndpropStr),
			clientv3.OpPut(lnpath, lnpropStr),
		).Else(
			clientv3.OpPut(ndpath, ndpropStr),
			clientv3.OpPut(lnpath, lnpropStr),
		).Commit()
	})
	atomic.AddUint32(&g.stats.DbWrite, uint32(len(resp.Responses)))
	nret := &ifc.GNode{
		Id:         ndpath,
		Type:       nodeType,
		Properties: ndprop,
		Links:      g.getLinks(ndpath),
		RLinks:     g.getRLinks((ndpath)),
	}
	for k, v := range ndpropC {
		nret.Properties[k] = v
	}
	lret := &ifc.GLink{
		Id:                lnpath,
		LinkType:          linkType,
		Properties:        lnprop,
		DestinationNodeId: ndpath,
		SourceNodeId:      parentNodeId,
	}
	for k, v := range lnpropC {
		lret.Properties[k] = v
	}
	nret.Properties[common.NodeFixedProp_Revision] = resp.Responses[0].GetResponsePut().Header.Revision
	lret.Properties[common.LinkFixedProp_Revision] = resp.Responses[1].GetResponsePut().Header.Revision
	return nret, lret

}

func (g *GraphDB) DescribeNode(nodeId, nodeType string) *ifc.GNode {
	lck := g.scheduler.Wait(nodeId)
	defer g.scheduler.Done(lck)
	g.checkId(nodeId)
	npropSKVP := g.retryGet(func() (*clientv3.GetResponse, error) {
		return g.client.KV.Get(context.TODO(), nodeId)
	})
	npropSCP := g.retryGet(func() (*clientv3.GetResponse, error) {
		return g.client.KV.Get(context.TODO(), nodeId+"/"+FixedPath_Created)
	})
	if len(npropSKVP.Kvs) == 0 || len(npropSCP.Kvs) == 0 {
		logging.Debugf("describeNode: nodeId=%s %s", nodeId, npropSKVP)
		return nil
	}

	atomic.AddUint32(&g.stats.DbRead, 2)
	atomic.AddUint32(&g.stats.NodeReadCnt, 1)
	var nprop ifc.PropertyType
	if err := json.Unmarshal(npropSKVP.Kvs[0].Value, &nprop); err != nil {
		logging.Fatalf("Unable to json parse the property on node %s = %s err=%s", nodeId, nprop, err)
	}
	var ncprop ifc.PropertyType
	if err := json.Unmarshal(npropSCP.Kvs[0].Value, &ncprop); err != nil {
		logging.Fatalf("Unable to json parse the property on nodec %s = %s err=%s", nodeId, ncprop, err)
	}
	links := g.getLinks(nodeId)
	rlinks := g.getRLinks(nodeId)
	pa := strings.Split(nodeId, "/")
	nodeTypeL := pa[len(pa)-2]
	ret := &ifc.GNode{
		Id:         nodeId,
		Links:      links,
		RLinks:     rlinks,
		Properties: nprop, // { ...nprop, ...ncprop },
		Type:       nodeTypeL,
	}
	for k, v := range ncprop {
		ret.Properties[k] = v
	}
	ret.Properties[common.NodeFixedProp_Revision] = npropSKVP.Kvs[0].ModRevision
	// fmt.Printf("DescribeNode nprop = %s node ptr = %p ", nprop, ret)
	// logger.debug(`describeNode: ${nodeId} RETURNING: ${JSON.stringify(ret)}`);
	// logger.debug(`         LINKS = ${JSON.stringify(links)}`);
	return ret

}

func (g *GraphDB) DeleteNode(nodeId string) {
	logging.Debugf("deleteNode nodeId=%s", nodeId)
	g.checkId(nodeId)
	// Also need to delete the link on parent node.
	np := strings.Split(nodeId, "/")

	if len(np) > 3 {
		// np.splice(np.length - 2, 0, FixedPath_Links);
		np = append(np[:len(np)-2], append([]string{FixedPath_Links}, np[len(np)-2:]...)...)
		pnl := strings.Join(np, "/")
		logging.Debugf("deleteNode Link %s\n\n", pnl)

		// g.client.KV.Delete(context.TODO(), pnl)
		// g.client.KV.Delete(context.TODO(), pnl + "/" + FixedPath_Created)
		g.retryTxn(func() (*clientv3.TxnResponse, error) {
			return g.client.KV.Txn(context.TODO()).
				If().Then(
				clientv3.OpDelete(nodeId+"/", clientv3.WithPrefix()),
				clientv3.OpDelete(nodeId),
				clientv3.OpDelete(pnl),
				clientv3.OpDelete(pnl+"/"+FixedPath_Created),
			).Commit()
		})
	} else {
		// TODO: Add support for deleting soft links as well.
		g.retryTxn(func() (*clientv3.TxnResponse, error) {
			return g.client.KV.Txn(context.TODO()).
				If().Then(
				clientv3.OpDelete(nodeId, clientv3.WithPrefix()),
			).Commit()
		})
	}

	if len(np) > 3 {
		atomic.AddUint32(&g.stats.DbWrite, 4)
	} else {
		atomic.AddUint32(&g.stats.DbWrite, 2)
	}
	// done with delete
}

func (g *GraphDB) UpdateNodeAddProperties(requestorId, nodeId string, propertiesIn ifc.PropertyType) int64 {
	logging.Debugf("updateNodeAddProperties nodeId=%s, requestor=%s prop = %s", nodeId, requestorId, propertiesIn)
	g.checkId(nodeId)
	properties := g.copyPropWithoutFixed(propertiesIn)
	g.checkId(nodeId)
	if len(properties) == 0 {
		return 0
	}
	now := time.Now().UTC().Format(utils.ISOTimeFormat)
	changeId := shortuuid.New()
	properties[common.NodeFixedProp_updatedBy] = requestorId
	properties[common.NodeFixedProp_updateTime] = now
	properties[common.NodeFixedProp_changeId] = changeId
	logging.Debugf("Lock and add properties nodeId=%s %s", nodeId, properties)
	ctxn := 0
	rev := int64(0)
	for {
		logging.Debugf("Update node add prop for node %s", nodeId)
		curPropS := g.retryGet(func() (*clientv3.GetResponse, error) {
			return g.client.KV.Get(context.TODO(), nodeId)
		})
		atomic.AddUint32(&g.stats.DbRead, 1)
		if len(curPropS.Kvs) == 0 {
			// may be node was deleted
			logging.Warnf("Unable to update property for %s node not found (Node add)", nodeId)
			return 0
		}
		var curProp ifc.PropertyType
		if err := json.Unmarshal(curPropS.Kvs[0].Value, &curProp); err != nil {
			panic(errors.New(fmt.Sprintf("Unable to json parse the property on node %s (Node add)", nodeId)))
		}
		logging.Debugf(" Before %s", curProp)
		for k, v := range properties {
			curProp[k] = v
		}
		logging.Debugf(" After %s", curProp)
		logging.Debugf(
			"Update node add prop for node %s, prior to txn commit, node KV mod revision is "+
				"%d,  node KV version is %d "+
				"node header rev is %d",
			nodeId, curPropS.Kvs[0].ModRevision,
			curPropS.Kvs[0].Version,
			curPropS.Header.Revision)
		curPropOut := string(utils.JsonMarshal(curProp))
		/* PLEASE READ: We can use the KV mod revision or the KV version to check if the KV has been
		   modified in transit. In case this comparison fails, the revision number observed in
		   the txn response should not be incremented as compared to the one read prior. This
		   behavior should be accounted for and tested in the DM/sys test scripts.
		*/

		presp := g.retryTxn(func() (*clientv3.TxnResponse, error) {
			return g.client.KV.Txn(context.TODO()).If(
				clientv3.Compare(clientv3.ModRevision(nodeId), "=", curPropS.Kvs[0].ModRevision),
			).Then(clientv3.OpPut(nodeId, curPropOut)).Else().Commit()
		})

		/*
			nPropS := g.retryGet(func() (*clientv3.GetResponse, error) {
				return g.client.KV.Get(context.TODO(), nodeId)
			})
			if len(nPropS.Kvs) != 0 {
				logging.Debugf(
					"Update node add prop for node %s, after txn commit, "+
					"node KV mod revision is %d "+
					"node KV version is %d",
					nodeId, nPropS.Kvs[0].ModRevision,
					nPropS.Kvs[0].Version)
			} else {
				logging.Debugf("Update node add prop for node %s: KVS len is 0 in Txn response", nodeId)
			}
		*/

		rev = presp.Header.Revision
		prevRev := curPropS.Header.Revision

		/* PLEASE READ: We need to use the revision number in the range response to identify
			   a successful prop update. Txn commit call does not include KV version info.
			   A revision upgrade should also correspond to a mod_version upgrade on subsequent reads
			   for the KV pair, in case txn commit is successful.
		           Should the incoming prop key value map not include any update for non fixed prop,
		           updateTime and changeId are updated for every successful commit.
		           Hence, we expect the header revision and KV[0] mod_version to be bumped up.
		           Please note that increments in version and revision numbers are not linear.
		           No strict assumption should be made in the DM /sys test scripts.
		           Please note that in case the same prop is written to the KV pair, mod revisions and versions
		           are still updated, even though they are no differences between the prop map and the existing
		           KV entry in the DB. No DB txn optimization occurs here. Optimize for throughput.
		*/

		logging.Debugf("UpdateNodeAddProperties: cur rev is %d, prev rev is %d", rev, prevRev)
		if prevRev >= rev {
			ctxn++
			if ctxn > g.maxTxnRetryCount {
				logging.Fatalf("ERROR: MaxTxn Retry Count Exceeded for node add prop for node ID %s", nodeId)
			}
			time.Sleep(time.Duration(g.txnRetryDelay) * time.Millisecond)
			continue
		}
		atomic.AddUint32(&g.stats.DbWrite, 1)
		break
	}
	return rev
}
func (g *GraphDB) UpdateNodeRemoveProperties(requestorId, nodeId string, propertiesIn []string) int64 {
	g.checkId(nodeId)
	propertieKeys := []string{}
	for _, k := range propertiesIn {
		if _, ok := g.fixedProp[k]; !ok {
			propertieKeys = append(propertieKeys, k)
		}
	}
	if len(propertieKeys) == 0 {
		return 0
	}
	logging.Debugf("updateNodeRemoveProperties nodeId=%s, prop = %s", nodeId, propertieKeys)
	// do not remove fixed properties.
	properties := make(ifc.PropertyType)
	now := time.Now().UTC().Format(utils.ISOTimeFormat)
	changeId := shortuuid.New()
	properties[common.NodeFixedProp_updatedBy] = requestorId
	properties[common.NodeFixedProp_updateTime] = now
	properties[common.NodeFixedProp_changeId] = changeId
	logging.Debugf("updateNodeRemoveProperties: Lock and Del properties nodeId=%s %s", nodeId, propertieKeys)
	ctxn := 0
	rev := int64(0)
	for {
		logging.Debugf("Update node remove prop for node %s", nodeId)
		curPropS := g.retryGet(func() (*clientv3.GetResponse, error) {
			return g.client.KV.Get(context.TODO(), nodeId)
		})
		atomic.AddUint32(&g.stats.DbRead, 1)
		if len(curPropS.Kvs) == 0 {
			// may be node was deleted
			logging.Warnf("Unable to update property for %s node not found (Node del)", nodeId)
			return 0
		}
		var curProp ifc.PropertyType
		if err := json.Unmarshal(curPropS.Kvs[0].Value, &curProp); err != nil {
			panic(errors.New(fmt.Sprintf("Unable to json parse the property on node %s (Node del)", nodeId)))
		}
		logging.Debugf(" Before %s", curProp)
		for _, k := range propertieKeys {
			delete(curProp, k)
		}
		logging.Debugf(" After %s", curProp)
		logging.Debugf(
			"Update node remove prop for node %s, prior to txn commit, node KV mod revision is "+
				"%d,  node KV version is %d "+
				"node header rev is %d",
			nodeId, curPropS.Kvs[0].ModRevision,
			curPropS.Kvs[0].Version,
			curPropS.Header.Revision)
		curPropOut := string(utils.JsonMarshal(curProp))
		/* PLEASE READ: We can use the KV mod revision or the KV version to check if the KV has been
		   modified in transit. In case this comparison fails, the revision number observed in
		   the txn response should not be incremented as compared to the one read prior. This
		   behavior should be accounted for and tested in the DM/sys test scripts.
		*/

		presp := g.retryTxn(func() (*clientv3.TxnResponse, error) {
			return g.client.KV.Txn(context.TODO()).If(
				clientv3.Compare(clientv3.ModRevision(nodeId), "=", curPropS.Kvs[0].ModRevision),
			).Then(clientv3.OpPut(nodeId, curPropOut)).Else().Commit()
		})
		/*
			nPropS := g.retryGet(func() (*clientv3.GetResponse, error) {
				return g.client.KV.Get(context.TODO(), nodeId)
			})
			if len(nPropS.Kvs) != 0 {
				logging.Debugf(
					"Update node remove prop for node %s, after txn commit, "+
						"node KV mod revision is %d "+
						"node KV version is %d",
					nodeId, nPropS.Kvs[0].ModRevision,
					nPropS.Kvs[0].Version)
			} else {
				logging.Debugf("Update node remove prop for node %s: KVS len is 0 in Txn response", nodeId)
			}
		*/
		rev = presp.Header.Revision
		prevRev := curPropS.Header.Revision

		/* PLEASE READ: We need to use the revision number in the range response to identify
			   a successful prop update. Txn commit call does not include KV version info.
			   A revision upgrade should also correspond to a mod_version upgrade on subsequent reads
			   for the KV pair, in case txn commit is successful.
		           Should the incoming prop key value map not include any update for non fixed prop,
		           updateTime and changeId are updated for every successful commit.
		           Hence, we expect the header revision and KV[0] mod_version to be bumped up.
		           Please note that increments in version and revision numbers are not linear.
		           No strict assumption should be made in the DM /sys test scripts.
		           Please note that in case the same prop is written to the KV pair, mod revisions and versions
		           are still updated, even though they are no differences between the prop map and the existing
		           KV entry in the DB. No DB txn optimization occurs here. Optimize for throughput.
		*/

		logging.Debugf("UpdateNodeRemoveProperties: cur rev is %d, prev rev is %d", rev, prevRev)
		if prevRev >= rev {
			ctxn++
			if ctxn > g.maxTxnRetryCount {
				logging.Fatalf("ERROR: MaxTxn Retry Count Exceeded for node rem prop for node ID %s", nodeId)
			}
			time.Sleep(time.Duration(g.txnRetryDelay) * time.Millisecond)
			continue
		}
		atomic.AddUint32(&g.stats.DbWrite, 1)
		break
	}
	atomic.AddUint32(&g.stats.NodeUpdateCnt, 1)
	return rev
}

/* Caller checks for argument validity*/
func (g *GraphDB) getSoftLinkProps(
	requestorId, linkType, srcId, destId, destNodeType string,
	destKeyValue, destNodePath string,
	destNodeProperties, linkPropertiesIn ifc.PropertyType,
	isSingletonLink bool) (ifc.PropertyType, ifc.PropertyType, string, string) {

	key := destNodeProperties[common.NodeFixedProp_NodeKeyName]
	nodeName := destKeyValue
	if isSingletonLink {
		key = common.NodeFixedProp_NodeDefaultKeyName
		nodeName = common.NodeFixedProp_NodeSingletonKeyValue
	}
	lnprop := g.copyPropWithoutFixed(linkPropertiesIn)
	now := time.Now().UTC().Format(utils.ISOTimeFormat)
	// changeId := shortuuid.New()
	lnpropC := make(map[string]interface{})
	lnpropC[common.LinkFixedProp_createdBy] = requestorId
	lnpropC[common.LinkFixedProp_creationTime] = now
	lnprop[common.LinkFixedProp_updatedBy] = requestorId
	lnprop[common.LinkFixedProp_updateTime] = now

	lnprop[common.LinkFixedProp_HardLink] = "false"
	lnprop[common.LinkFixedProp_NodeKeyName] = key
	lnprop[common.LinkFixedProp_NodeKeyValue] = nodeName
	lnprop[common.LinkFixedProp_NodeType] = destNodeType
	lnprop[common.LinkFixedProp_destNodeId] = destId

	lnpath := srcId + "/" + FixedPath_Links + "/" + destNodeType + "/" + nodeName
	lncpath := lnpath + "/" + FixedPath_Created

	lnprop[common.LinkFixedProp_SoftLinkDestinationPath] = destNodePath
	logging.Debugf("upsertLink path=%s :: %s", lnpath, lncpath)
	return lnprop, lnpropC, lnpath, lncpath

}

/* Caller checks for argument validity*/
func (g *GraphDB) getRSoftLinkProps(
	requestorId, linkType, srcId, destId, destNodeType string,
	destKeyValue, destNodePath string,
	destNodeProperties, linkPropertiesIn ifc.PropertyType,
	isSingletonLink bool) (ifc.PropertyType, ifc.PropertyType, string, string) {

	key := destNodeProperties[common.NodeFixedProp_NodeKeyName]
	nodeName := destKeyValue
	if isSingletonLink {
		key = common.NodeFixedProp_NodeDefaultKeyName
		nodeName = common.NodeFixedProp_NodeSingletonKeyValue
	}
	lnprop := g.copyPropWithoutFixed(linkPropertiesIn)
	now := time.Now().UTC().Format(utils.ISOTimeFormat)
	// changeId := shortuuid.New()
	lnpropC := make(map[string]interface{})
	lnpropC[common.LinkFixedProp_createdBy] = requestorId
	lnpropC[common.LinkFixedProp_creationTime] = now
	lnprop[common.LinkFixedProp_updatedBy] = requestorId
	lnprop[common.LinkFixedProp_updateTime] = now

	lnprop[common.LinkFixedProp_HardLink] = "false"
	lnprop[common.LinkFixedProp_NodeKeyName] = key
	lnprop[common.LinkFixedProp_NodeKeyValue] = nodeName
	lnprop[common.LinkFixedProp_NodeType] = destNodeType
	lnprop[common.LinkFixedProp_destNodeId] = destId

	lnpath := srcId + "/" + FixedPath_RLinks + "/" + destNodeType + "/" + nodeName
	lncpath := lnpath + "/" + FixedPath_Created
	logging.Debugf("upsertRLink path=%s :: %s", lnpath, lncpath)

	lnprop[common.LinkFixedProp_RSoftLinkDestinationPath] = destNodePath
	return lnprop, lnpropC, lnpath, lncpath

}

func (g *GraphDB) UpsertLink(
	requestorId string, lnk common.UpsertLinkOpObj) (*ifc.GLink, *ifc.GLink) {
	srcId := lnk.SrcNodeObj.NodeId
	destId := lnk.DestNodeObj.NodeId
	var rlnProp, rlnPropC ifc.PropertyType
	var rlnPath, rlncPath string
	var resp *clientv3.TxnResponse
	var rlret *ifc.GLink

	logging.Debugf("UpsertLink: srcId=%s destId=%s prop = %s", srcId, destId, lnk.LinkPropIn)
	g.checkId(destId)
	g.checkId(srcId)

	rLinkEnabled := g.isRLinkFeatureFlagEnabled()

	lnProp, lnPropC, lnPath, lncPath := g.getSoftLinkProps(requestorId, lnk.LinkType, lnk.SrcNodeObj.NodeId,
		lnk.DestNodeObj.NodeId, lnk.DestNodeObj.NodeType, lnk.DestNodeObj.NodeKeyValue, lnk.DestNodeObj.NodePath,
		lnk.DestNodeObj.NodeProp, lnk.LinkPropIn, lnk.IsSnglton)

	if rLinkEnabled {
		rlnProp, rlnPropC, rlnPath, rlncPath = g.getRSoftLinkProps(requestorId, common.LinkType_ROwner, lnk.DestNodeObj.NodeId,
			lnk.SrcNodeObj.NodeId, lnk.SrcNodeObj.NodeType, lnk.SrcNodeObj.NodeKeyValue, lnk.SrcNodeObj.NodePath,
			lnk.SrcNodeObj.NodeProp, ifc.PropertyType{}, lnk.IsSnglton)

		resp = g.retryTxn(func() (*clientv3.TxnResponse, error) {
			return g.client.KV.Txn(context.TODO()).If(
				clientv3.Compare(clientv3.CreateRevision(lncPath), "=", 0),
			).Then(
				clientv3.OpPut(lncPath, string(string(utils.JsonMarshal(lnPropC)))),
				clientv3.OpPut(lnPath, string(string(utils.JsonMarshal(lnProp)))),
				clientv3.OpPut(rlncPath, string(string(utils.JsonMarshal(rlnPropC)))),
				clientv3.OpPut(rlnPath, string(string(utils.JsonMarshal(rlnProp)))),
			).Else(
				clientv3.OpPut(lnPath, string(string(utils.JsonMarshal(lnProp)))),
				clientv3.OpPut(rlnPath, string(string(utils.JsonMarshal(rlnProp)))),
			).Commit()
		})
		logging.Debugf("upsertRLink path=%s DONE.", rlnPath)
		atomic.AddUint32(&g.stats.RLinkAddCnt, 1)
	} else {
		resp = g.retryTxn(func() (*clientv3.TxnResponse, error) {
			return g.client.KV.Txn(context.TODO()).If(
				clientv3.Compare(clientv3.CreateRevision(lncPath), "=", 0),
			).Then(
				clientv3.OpPut(lncPath, string(string(utils.JsonMarshal(lnPropC)))),
				clientv3.OpPut(lnPath, string(string(utils.JsonMarshal(lnProp)))),
			).Else(
				clientv3.OpPut(lnPath, string(string(utils.JsonMarshal(lnProp)))),
			).Commit()
		})
	}

	logging.Debugf("upsertLink path=%s DONE.", lnPath)
	atomic.AddUint32(&g.stats.DbWrite, uint32(len(resp.Responses)))
	atomic.AddUint32(&g.stats.LinkAddCnt, 1)

	lret := &ifc.GLink{
		Id:                lnPath,
		LinkType:          lnk.LinkType,
		Properties:        lnProp,
		DestinationNodeId: destId,
		SourceNodeId:      srcId,
	}
	for k, v := range lnPropC {
		lret.Properties[k] = v
	}
	lret.Properties[common.LinkFixedProp_Revision] = resp.Responses[0].GetResponsePut().Header.Revision
	if rLinkEnabled {
		rlret = &ifc.GLink{
			Id:                rlnPath,
			LinkType:          common.LinkType_ROwner,
			Properties:        rlnProp,
			DestinationNodeId: srcId,
			SourceNodeId:      destId,
		}
		for k, v := range rlnPropC {
			rlret.Properties[k] = v
		}
		rlret.Properties[common.LinkFixedProp_Revision] = resp.Responses[0].GetResponsePut().Header.Revision
	}
	return lret, rlret
}

func (g *GraphDB) DescribeLinkById(linkId string) (*ifc.GLink, bool) {
	logging.Debugf("describeLinkById %s", linkId)
	lnke := g.retryGet(func() (*clientv3.GetResponse, error) {
		return g.client.KV.Get(context.TODO(), linkId)
	})
	if len(lnke.Kvs) == 0 {
		// logging.Fatalf("link id = %s not found", linkId)
		return nil, false
	}
	lnk := string(lnke.Kvs[0].Value)
	atomic.AddUint32(&g.stats.DbRead, 1)
	atomic.AddUint32(&g.stats.LinkReadCnt, 1)
	lnkc := g.retryGet(func() (*clientv3.GetResponse, error) {
		return g.client.KV.Get(context.TODO(), linkId+"/"+FixedPath_Created)
	})
	if len(lnkc.Kvs) == 0 {
		return nil, false
	}
	lnkcs := lnkc.Kvs[0].Value
	var lprop ifc.PropertyType
	var lcprop ifc.PropertyType
	if err := json.Unmarshal(lnke.Kvs[0].Value, &lprop); err != nil {
		logging.Fatalf("DescribeLinkById error on unmarshal %s", lnk)
	}
	if err := json.Unmarshal(lnkcs, &lcprop); err != nil {
		logging.Fatalf("DescribeLinkByIfd error on unmarshal %s", lnkc)
	}
	srcSplit := strings.Split(linkId, "/")
	srcId := strings.Join(srcSplit[:len(srcSplit)-3], "/")

	r := &ifc.GLink{
		Id:                linkId,
		LinkType:          lprop[common.LinkFixedProp_NodeType].(string), // ["type"].(string),
		Properties:        lprop,
		DestinationNodeId: lprop[common.LinkFixedProp_destNodeId].(string),
		SourceNodeId:      srcId,
	}
	for k, v := range lcprop {
		r.Properties[k] = v
	}
	r.Properties[common.LinkFixedProp_Revision] = lnke.Kvs[0].ModRevision
	return r, true
}

func (g *GraphDB) DescribeLink(linkType, sourceId, destinaitonId string) (*ifc.GLink, bool) {
	logging.Debugf("describeLink sourceId=%s linkType=%s destinaitonId=%s", sourceId, linkType, destinaitonId)
	destKeyS := strings.Split(destinaitonId, "/")
	destKey := destKeyS[len(destKeyS)-1]
	lid := sourceId + "/" + FixedPath_Links + "/" + linkType + "/" + destKey
	return g.DescribeLinkById(lid)
}

/*This API is used only for slink deletion
Deletion of parent link from the db is done in DeleteNode()
*/
func (g *GraphDB) DeleteLink(linkId string) {
	logging.Debugf("DeleteLink: linkId = %s", linkId)
	g.retryTxn(func() (*clientv3.TxnResponse, error) {
		return g.client.KV.Txn(context.TODO()).If().Then(
			clientv3.OpDelete(linkId),
			clientv3.OpDelete(linkId+"/"+FixedPath_Created),
		).Commit()
	})
	atomic.AddUint32(&g.stats.DbWrite, 1)
	if strings.Contains(linkId, FixedPath_RLinks) {
		atomic.AddUint32(&g.stats.RLinkDelCnt, 1)
	}
	if strings.Contains(linkId, FixedPath_Links) {
		atomic.AddUint32(&g.stats.LinkDelCnt, 1)
	}
}

/*This API is used only for slink iand rlink pair deletion
Deletion of parent link from the db is done in DeleteNode()
*/
func (g *GraphDB) DeleteSoftLinkWithRSoftLink(linkId string, rLinkId string) {
	logging.Debugf("DeleteSoftLinkWithRSoftLink: linkId = %s , rLinkId = %s",
		linkId, rLinkId)
	g.retryTxn(func() (*clientv3.TxnResponse, error) {
		return g.client.KV.Txn(context.TODO()).If().Then(
			clientv3.OpDelete(linkId),
			clientv3.OpDelete(linkId+"/"+FixedPath_Created),
			clientv3.OpDelete(rLinkId),
			clientv3.OpDelete(rLinkId+"/"+FixedPath_Created),
		).Commit()
	})
	atomic.AddUint32(&g.stats.DbWrite, 2)
	atomic.AddUint32(&g.stats.RLinkDelCnt, 1)
	atomic.AddUint32(&g.stats.LinkDelCnt, 1)
}

func (g *GraphDB) deleteLinks(sLinkId, rLinkId string) {
	logging.Debugf("DeleteLink: Forward softlinkId = %s", sLinkId)
	logging.Debugf("DeleteLink: Reverse softlinkId = %s", rLinkId)

	g.retryTxn(func() (*clientv3.TxnResponse, error) {
		return g.client.KV.Txn(context.TODO()).If().Then(
			clientv3.OpDelete(sLinkId),
			clientv3.OpDelete(sLinkId+"/"+FixedPath_Created),
			clientv3.OpDelete(rLinkId),
			clientv3.OpDelete(rLinkId+"/"+FixedPath_Created),
		).Commit()
	})
	atomic.AddUint32(&g.stats.DbWrite, 1)
	atomic.AddUint32(&g.stats.RLinkDelCnt, 1)
	atomic.AddUint32(&g.stats.LinkDelCnt, 1)
}

func (g *GraphDB) DeleteLinks(nodeId string) []*ifc.GLink {
	logging.Debugf("Looking up reverse links for node id %s in db", nodeId)
	rlinks := g.getRLinks(nodeId)
	if len(rlinks) == 0 {
		logging.Debugf("Reverse links not found for node %s in db\n", nodeId)
		return rlinks
	}
	gLinkIdSet := strings.Split(nodeId, "/")
	l := len(gLinkIdSet)
	nodeType := gLinkIdSet[l-2]
	nodeKeyValue := gLinkIdSet[l-1]

	/*
			Predicate delete of forward softlinks based on reverse soft links,
			to preserve backward compatibility.
		 	We need to exclude links that pertain to child relationships.
	*/
	for _, gLinkObj := range rlinks {
		logging.Debugf("Reverse link id is %s \n Link object: %v", gLinkObj.Id, gLinkObj)
		fPath := gLinkObj.DestinationNodeId + "/" + FixedPath_Links + "/" + nodeType + "/" + nodeKeyValue
		/*
			No need to check here if the node in dest Path exists
			Graph db interface does this check
		*/
		g.deleteLinks(fPath, gLinkObj.Id)
		logging.Debugf("Deleted softlink %s in db\nDeleted reverse softlink %s in db",
			fPath, gLinkObj.Id)
	}

	return rlinks
}

func (g *GraphDB) UpdateLinkAddProperties(
	requestorId, linkId string, propertiesIn ifc.PropertyType) int64 {
	logging.Debugf("UpdateLinkAddProperties: linkId = %s prop = %s", linkId, propertiesIn)

	properties := g.copyPropWithoutFixed(propertiesIn)
	if len(properties) == 0 {
		return 0
	}
	now := time.Now().UTC().Format(utils.ISOTimeFormat)
	// changeId := shortuuid.New()
	properties[common.LinkFixedProp_updatedBy] = requestorId
	properties[common.LinkFixedProp_updateTime] = now
	logging.Debugf("Lock link and add properties %s", linkId)
	ctxn := 0
	rev := int64(0)
	for {
		logging.Debugf("Update link add prop for node %s", linkId)
		curPropS := g.retryGet(func() (*clientv3.GetResponse, error) {
			return g.client.KV.Get(context.TODO(), linkId)
		})
		atomic.AddUint32(&g.stats.DbRead, 1)
		if len(curPropS.Kvs) == 0 {
			// may be node was deleted
			logging.Warnf("Unable to update property for link %s not found (Link add)", linkId)
			return 0
		}
		var curProp ifc.PropertyType
		if err := json.Unmarshal(curPropS.Kvs[0].Value, &curProp); err != nil {
			panic(errors.New(fmt.Sprintf("Unable to json parse the property on link %s (Link add)", linkId)))
		}
		logging.Debugf(" Before %s", curProp)
		for k, v := range properties {
			curProp[k] = v
		}
		logging.Debugf(" After %s", curProp)
		logging.Debugf(
			"Update link add prop for link %s, prior to txn commit, link KV mod revision is "+
				"%d,  link KV version is %d "+
				"link header rev is %d",
			linkId, curPropS.Kvs[0].ModRevision,
			curPropS.Kvs[0].Version,
			curPropS.Header.Revision)
		curPropOut := string(utils.JsonMarshal(curProp))

		/* PLEASE READ: We can use the KV mod revision or the KV version to check if the KV has been
		   modified in transit. In case this comparison fails, the revision number observed in
		   the txn response should not be incremented as compared to the one read prior. This
		   behavior should be accounted for and tested in the DM/sys test scripts.
		*/

		presp := g.retryTxn(func() (*clientv3.TxnResponse, error) {
			return g.client.KV.Txn(context.TODO()).If(
				clientv3.Compare(clientv3.ModRevision(linkId), "=", curPropS.Kvs[0].ModRevision),
			).Then(clientv3.OpPut(linkId, curPropOut)).Else().Commit()
		})
		/*
			nPropS := g.retryGet(func() (*clientv3.GetResponse, error) {
				return g.client.KV.Get(context.TODO(), linkId)
			})
			if len(nPropS.Kvs) != 0 {
				logging.Debugf(
					"Update link add prop for node %s, after txn commit, "+
						"link KV mod revision is %d "+
						"link KV version is %d",
					linkId, nPropS.Kvs[0].ModRevision,
					nPropS.Kvs[0].Version)
			} else {
				logging.Debugf("Update link add prop for node %s: KVS len is 0 in Txn response", linkId)
			}
		*/
		rev = presp.Header.Revision
		prevRev := curPropS.Header.Revision

		/* PLEASE READ: We need to use the revision number in the range response to identify
			   a successful prop update. Txn commit call does not include KV version info.
			   A revision upgrade should also correspond to a mod_version upgrade on subsequent reads
			   for the KV pair, in case txn commit is successful.
		           Should the incoming prop key value map not include any update for non fixed prop,
		           updateTime and changeId are updated for every successful commit.
		           Hence, we expect the header revision and KV[0] mod_version to be bumped up.
		           Please note that increments in version and revision numbers are not linear.
		           No strict assumption should be made in the DM /sys test scripts.
		           Please note that in case the same prop is written to the KV pair, mod revisions and versions
		           are still updated, even though they are no differences between the prop map and the existing
		           KV entry in the DB. No DB txn optimization occurs here. Optimize for throughput.
		*/

		logging.Debugf("UpdateLinkAddProperties: cur rev is %d, prev rev is %d", rev, prevRev)
		if prevRev >= rev {
			ctxn++
			if ctxn > g.maxTxnRetryCount {
				logging.Fatalf("ERROR: MaxTxn Retry Count Exceeded for link add prop for node ID %s", linkId)
			}
			time.Sleep(time.Duration(g.txnRetryDelay) * time.Millisecond)
			continue
		}
		atomic.AddUint32(&g.stats.DbWrite, 1)
		break
	}
	atomic.AddUint32(&g.stats.LinkUpdateCnt, 1)
	return rev
}
func (g *GraphDB) UpdateLinkRemoveProperties(requestorId, linkId string, propertiesIn []string) int64 {
	// do not remove fixed properties.
	propertieKeys := []string{}
	for _, k := range propertiesIn {
		if _, ok := g.fixedProp[k]; !ok {
			propertieKeys = append(propertieKeys, k)
		}
	}
	if len(propertieKeys) == 0 {
		return 0
	}
	logging.Debugf("UpdateLinkRemoveProperties linkId=%s, prop = %s", linkId, propertieKeys)
	properties := make(ifc.PropertyType)
	now := time.Now().UTC().Format(utils.ISOTimeFormat)
	changeId := shortuuid.New()
	properties[common.NodeFixedProp_updatedBy] = requestorId
	properties[common.NodeFixedProp_updateTime] = now
	properties[common.NodeFixedProp_changeId] = changeId
	ctxn := 0
	rev := int64(0)
	for {
		logging.Debugf("Update link remove prop for node %s", linkId)
		curPropS := g.retryGet(func() (*clientv3.GetResponse, error) {
			return g.client.KV.Get(context.TODO(), linkId)
		})
		atomic.AddUint32(&g.stats.DbRead, 1)
		if len(curPropS.Kvs) == 0 {
			// may be node was deleted
			logging.Warnf("Unable to update property for link %s not found (Link del)", linkId)
			return 0
		}
		var curProp ifc.PropertyType
		if err := json.Unmarshal(curPropS.Kvs[0].Value, &curProp); err != nil {
			panic(errors.New(fmt.Sprintf("Unable to json parse the property on link %s (Link del)", linkId)))
		}
		logging.Debugf(" Before %s", curProp)
		for _, k := range propertieKeys {
			delete(curProp, k)
		}
		logging.Debugf(" After %s", curProp)
		logging.Debugf(
			"Update link remove prop for link %s, prior to txn commit, link KV mod revision is "+
				"%d,  link KV version is %d "+
				"link header rev is %d",
			linkId, curPropS.Kvs[0].ModRevision,
			curPropS.Kvs[0].Version,
			curPropS.Header.Revision)
		curPropOut := string(utils.JsonMarshal(curProp))

		/* PLEASE READ: We can use the KV mod revision or the KV version to check if the KV has been
		   modified in transit. In case this comparison fails, the revision number observed in
		   the txn response should not be incremented as compared to the one read prior. This
		   behavior should be accounted for and tested in the DM/sys test scripts.
		*/

		presp := g.retryTxn(func() (*clientv3.TxnResponse, error) {
			return g.client.KV.Txn(context.TODO()).If(
				clientv3.Compare(clientv3.ModRevision(linkId), "=", curPropS.Kvs[0].ModRevision),
			).Then(clientv3.OpPut(linkId, curPropOut)).Else().Commit()
		})
		/*
			nPropS := g.retryGet(func() (*clientv3.GetResponse, error) {
				return g.client.KV.Get(context.TODO(), linkId)
			})
			if len(nPropS.Kvs) != 0 {
				logging.Debugf(
					"Update link remove prop for node %s, after txn commit, "+
					"link KV mod revision is %d "+
					"link KV version is %d",
					linkId, nPropS.Kvs[0].ModRevision,
					nPropS.Kvs[0].Version)
			} else {
				logging.Debugf("Update link remove prop for node %s: KVS len is 0 in Txn response", linkId)
			}
		*/
		rev = presp.Header.Revision
		prevRev := curPropS.Header.Revision

		/* PLEASE READ: We need to use the revision number in the range response to identify
			   a successful prop update. Txn commit call does not include KV version info.
			   A revision upgrade should also correspond to a mod_version upgrade on subsequent reads
			   for the KV pair, in case txn commit is successful.
		           Should the incoming prop key value map not include any update for non fixed prop,
		           updateTime and changeId are updated for every successful commit.
		           Hence, we expect the header revision and KV[0] mod_version to be bumped up.
		           Please note that increments in version and revision numbers are not linear.
		           No strict assumption should be made in the DM /sys test scripts.
		           Please note that in case the same prop is written to the KV pair, mod revisions and versions
		           are still updated, even though they are no differences between the prop map and the existing
		           KV entry in the DB. No DB txn optimization occurs here. Optimize for throughput.
		*/

		logging.Debugf("UpdateLinkRemoveProperties: cur rev is %d, prev rev is %d", rev, prevRev)
		if prevRev >= rev {
			ctxn++
			if ctxn > g.maxTxnRetryCount {
				logging.Fatalf("ERROR: MaxTxn Retry Count Exceeded for link remove prop for node ID %s", linkId)
			}
			time.Sleep(time.Duration(g.txnRetryDelay) * time.Millisecond)
			continue
		}
		atomic.AddUint32(&g.stats.DbWrite, 1)
		break
	}
	atomic.AddUint32(&g.stats.LinkUpdateCnt, 1)
	return rev
}

func (g *GraphDB) getNextKeyList(prefix string, cnt int, k string) []string {
	if cnt == 0 {
		return []string{}
	}
	var limit int64 = int64(cnt) * 2
	st := prefix
	var startKey string = ""
	if k != "" {
		st = st + "/" + k
		limit = limit + 2
		startKey = st
	} else {
		st = st + "\x00"
	}
	ed := prefix + "\xff"

	curPropS := g.retryGet(func() (*clientv3.GetResponse, error) {
		return g.client.KV.Get(context.TODO(), st,
			clientv3.WithRange(ed), clientv3.WithLimit(limit))
	})
	logging.Debugf("Iterator: DB Fetched rlink KVS orig len:%d", len(curPropS.Kvs))
	ret := []string{}
	if k == "" && len(curPropS.Kvs) != 0 {
		for _, dt := range curPropS.Kvs {
			ret = append(ret, string(dt.Key))
		}
	} else if len(curPropS.Kvs) >= 2 {
		for idx, dt := range curPropS.Kvs {
			if idx != 0 {
				ret = append(ret, string(dt.Key))
			} else if string(dt.Key) != startKey {
				// case to handle when passed key k does not exist in the DB
				ret = append(ret, string(dt.Key))
			}
		}
	}
	logging.Debugf("Iterator: DB Fetched rlink keys orig set : %v\n", ret)
	finalRet := []string{}
outer:
	for _, v := range ret {
		tv := v[len(prefix)+1:]
		parts := strings.Split(tv, "/")
		for idx, part := range parts {
			if idx%2 != 0 && !strings.Contains(part, FixedPath_Created) && !strings.Contains(part, FixedPath_Links) &&
				!strings.Contains(part, FixedPath_Lock) {
				continue outer
			}
		}

		if !strings.Contains(tv, "/") && len(finalRet) != cnt {
			finalRet = append(finalRet, tv)
		}
	}
	logging.Debugf("Iterator: DB  Return rlink keys final set : %v\n", finalRet)
	return finalRet
}

// iterator
func (g *GraphDB) GetNextLinkKey(sourceId, linkType, currentKey string, cnt int) []string {
	lid := sourceId + "/" + FixedPath_Links + "/" + linkType
	return g.getNextKeyList(lid, cnt, currentKey)
}

func (g *GraphDB) getNextRLinkKeyList(prefix string, cnt int, k string) []string {
	if cnt == 0 {
		return []string{}
	}
	var limit int64 = int64(cnt) * 2
	st := prefix
	var startKey string = ""
	if k != "" {
		st = st + "/" + k
		limit = limit + 2
		startKey = st
	} else {
		st = st + "\x00"
	}
	ed := prefix + "\xff"

	curPropS := g.retryGet(func() (*clientv3.GetResponse, error) {
		return g.client.KV.Get(context.TODO(), st,
			clientv3.WithRange(ed), clientv3.WithLimit(limit))
	})
	m := make(map[string]string)
	ret := []string{}
	var dtPropMap ifc.PropertyType
	if k == "" && len(curPropS.Kvs) != 0 {
		for _, dt := range curPropS.Kvs {
			ret = append(ret, string(dt.Key))
			err := json.Unmarshal(dt.Value, &dtPropMap)
			if err == nil {
				m[string(dt.Key)] = dtPropMap[common.LinkFixedProp_destNodeId].(string)
			}
		}
	} else if len(curPropS.Kvs) >= 2 {
		for idx, dt := range curPropS.Kvs {
			appendOp := false
			if idx != 0 {
				ret = append(ret, string(dt.Key))
				appendOp = true
			} else if string(dt.Key) != startKey {
				// case to handle when passed key k does not exist in the DB
				ret = append(ret, string(dt.Key))
				appendOp = true
			}
			if appendOp {
				err := json.Unmarshal(dt.Value, &dtPropMap)
				// We can encounter a _created entry here, prop map for _Created
				// entries does not include _DID prop key. Include a check
				if err == nil {
					if _, dOk := dtPropMap[common.LinkFixedProp_destNodeId]; dOk {
						m[string(dt.Key)] = dtPropMap[common.LinkFixedProp_destNodeId].(string)
					}
				}
			}
		}
	}

	finalRet := []string{}
outer:
	for _, v := range ret {
		tv := v[len(prefix)+1:]
		parts := strings.Split(tv, "/")
		for idx, part := range parts {
			if idx%2 != 0 && !strings.Contains(part, FixedPath_Created) && !strings.Contains(part, FixedPath_Links) &&
				!strings.Contains(part, FixedPath_Lock) {
				continue outer
			}
		}

		if !strings.Contains(tv, "/") && len(finalRet) != cnt {
			if _, kOk := m[v]; kOk {
				finalRet = append(finalRet, m[v])
			}
		}
	}
	return finalRet
}

// rlink  iterator
func (g *GraphDB) GetNextRLinkKey(sourceId, linkType, currentKey string, cnt int) []string {
	lid := sourceId + "/" + FixedPath_RLinks + "/" + linkType
	return g.getNextRLinkKeyList(lid, cnt, currentKey)
}

// just update node property
func (g *GraphDB) GetNodeProperty(nodeId string) ifc.PropertyType {
	lck := g.scheduler.Wait(nodeId)
	defer g.scheduler.Done(lck)
	logging.Debugf("GETNODEPROPERTY key=%s", nodeId)
	npropSe := g.retryGet(func() (*clientv3.GetResponse, error) {
		return g.client.KV.Get(context.TODO(), nodeId)
	})
	npropSC := g.retryGet(func() (*clientv3.GetResponse, error) {
		return g.client.KV.Get(context.TODO(), nodeId+"/"+FixedPath_Created)
	})
	if len(npropSe.Kvs) == 0 || len(npropSC.Kvs) == 0 {
		return nil
	}
	atomic.AddUint32(&g.stats.DbRead, 2)
	atomic.AddUint32(&g.stats.NodeReadCnt, 1)
	var nprop ifc.PropertyType
	var ncprop ifc.PropertyType
	if err := json.Unmarshal(npropSe.Kvs[0].Value, &nprop); err != nil {
		logging.Fatalf("GetNodeProperty error on unmarshal nprop %s", nprop)
	}
	if err := json.Unmarshal(npropSC.Kvs[0].Value, &ncprop); err != nil {
		logging.Fatalf("GetNodeProperty error on unmarshal ncprop %s", ncprop)
	}
	ret := nprop
	for k, v := range ncprop {
		ret[k] = v
	}
	ret[common.NodeFixedProp_Revision] = npropSe.Kvs[0].ModRevision
	logging.Debugf("getNodeProperty: nodeId=%s Ret:%s", nodeId, ret)
	return ret
}

func (g *GraphDB) GetStats() common.GraphDBStats {
	return g.stats
}
