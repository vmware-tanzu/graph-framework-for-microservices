package base

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/internal/base/link"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/internal/utils"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/ifc"
	linkItr "gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/linkiterator"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/logging"

	"github.com/golang/protobuf/proto"
	"github.com/mitchellh/mapstructure"
)

type ReverseLink struct {
	NodeKey string
	LinkId  string
}
type BaseNode struct {
	ID          string                           `json:"id"`
	Ccnt        uint64                           `json:"ccnt"`
	NodeType    string                           `json:"nodeType"`
	DM          ifc.DataModelInterface           `json:"-"`
	NodeKeyName string                           `json:"nodeKeyName"`
	TypeRefNode bool                             `json:"typeRefNode"`
	Properties  ifc.PropertyType                 `json:"properties"`
	Links       ifc.BaseNodeLinkManagerInterface `json:"links"`
	RSoftLink   map[string]ReverseLink           `json:"rSoftLink"`
	sync.RWMutex
	Deleted         bool                  `json:"deleted"`
	ToBeDeleted     bool                  `json:"toBeDeleted"`
	BaseParent      ifc.BaseNodeInterface `json:"baseParent"`
	DeleteCompleted bool                  `json:"deleteCompleted"`
}

var baseNodeCreateCounter uint64 = 1000

func NewBaseNode(
	dm ifc.DataModelInterface, nodeType, keyName, id string,
	prop ifc.PropertyType, typeRef bool) *BaseNode {
	e := &BaseNode{
		ID:              id,
		Ccnt:            atomic.AddUint64(&baseNodeCreateCounter, 1),
		NodeType:        nodeType,
		DM:              dm,
		NodeKeyName:     keyName,
		TypeRefNode:     typeRef,
		Properties:      prop,
		Links:           link.NewBaseNodeLinkManager(),
		RSoftLink:       make(map[string]ReverseLink),
		Deleted:         false,
		BaseParent:      nil,
		DeleteCompleted: false,
	}
	return e
}

func (b *BaseNode) GetLinks() ifc.BaseNodeLinkManagerInterface {
	return b.Links
}
func (b *BaseNode) GetCCnt() uint64 {
	return b.Ccnt
}
func (b *BaseNode) Sync(forceSync bool) {
	b.DM.Sync(b, forceSync)
}

func (b *BaseNode) GetCreatorId() string {
	return b.DM.GetId()
}

func (b *BaseNode) GetLatestRevision() int64 {
	return b.DM.GetLatestRevision()
}

func (b *BaseNode) CompactRevision(revision, compactCtxtTimeout int64, compactPhysical bool) (*ifc.EtcdCompactResponse, error) {
	return b.DM.CompactRevision(revision, compactCtxtTimeout, compactPhysical)
}

func (b *BaseNode) checkRLinkFeatureFlagStatus() {
	if !b.IsRLinkFeatureFlagEnabled() {
		logging.Fatalf(fmt.Sprintf("DM: %s: Unsupported operation. Reverse link feature flag support is disabled.",
			b.DM.GetName()))
	}
}

func (b *BaseNode) SetCreatorId(seed string) string {
	b.DM.SetId(seed)
	return b.DM.GetId()
}

func (b *BaseNode) ClaimOwnership(path ifc.NodePathList, timeDelay uint32, prefix string) {
	b.DM.ClaimOwnership(path, timeDelay, prefix)
}

func (b *BaseNode) ClaimLinkOwnership(path ifc.NodePathList, destNodeType string, timeDelay uint32) {
	b.DM.ClaimLinkOwnership(path, destNodeType, timeDelay)
}

func (b *BaseNode) RegisterCB(pattern ifc.NodePathList, cbfn ifc.CallbackFuncNode) {
	b.DM.RegisterCB(pattern, cbfn)
}

func (b *BaseNode) RegisterLinkCB(pattern ifc.NodePathList, destNodeType string, cbfn ifc.CallbackFuncLink) {
	b.DM.RegisterLinkCB(pattern, destNodeType, cbfn)
}

func (b *BaseNode) GetTypeRefChildNode(ntype, nkeyname, nkeyvalue string) ifc.BaseNodeInterface {
	p := make(ifc.PropertyType)
	p[nkeyname] = nkeyvalue
	n := NewBaseNode(b.DM, ntype, nkeyname, "", p,
		true)
	var bifc ifc.BaseNodeInterface = b
	n.SetBaseParent(bifc)
	return n
}

func (b *BaseNode) Subscribe(pattern ifc.NodePathList, depth uint32) {
	b.DM.Subscribe(pattern, depth)
}

func (b *BaseNode) Unsubscribe(pattern ifc.NodePathList) {
	b.DM.Unsubscribe(pattern)
}

func (b *BaseNode) IsAnyParentDeleted(level ...uint32) bool {
	var l uint32 = 0
	if len(level) > 0 {
		l = level[0]
	}

	b.RLock()
	id := b.ID
	if b.DeleteCompleted {
		b.RUnlock()
		logging.Debugf("IsAnyParentDeleted(): Node %s Returning true for deleted state",
			b.GetId())
		return true
	}
	if b.BaseParent == nil {
		b.RUnlock()
		logging.Debugf("IsAnyParentDeleted(): Node %s Returning false for deleted state, base parent is nil",
			b.GetId())
		return false
	}

	if l > 32 {
		b.RUnlock()
		panic(errors.New("anyParentDeleted Tree Depth Exceeded > 32" + id))
	}

	result := b.BaseParent.IsAnyParentDeleted(l + 1)
	b.RUnlock()
	markedStr := "is marked"
	if !result {
		markedStr = "is not marked"
	}
	logging.Debugf("IsAnyParentDeleted(): Node %s Returning %t for deleted state, one of the parents %s"+
		" for deleteCompleted",
		b.GetId(), result, markedStr)

	return result
}

func (b *BaseNode) GetFullPath() ifc.NodePathList {
	idp := strings.Split(b.ID, "/")
	idpPtr := 0
	if len(idp) > 2 {
		idpPtr++
		var retPath ifc.NodePathList = []ifc.NodePath{}
		for idpPtr < len(idp) {
			nt := idp[idpPtr]
			nk := idp[idpPtr+1]
			idpPtr += 2
			if nt == "Root" {
				nk = "default"
			}
			retPath = append(retPath, ifc.NodePath{nt, nk})
			if idpPtr == len(idp) {
				return retPath
			}
		}
		return retPath
	} else {
		panic(errors.New("GetFullPath called with ID =" + b.ID))
	}
}

func (b *BaseNode) SetImmediateBaseProperties(pnew ifc.PropertyType) {
	/*
		Need a check for tobedeleted condition here
		Skip now to preserve backward compatibility
	*/

	if b.IsDeleted() {
		return
	}
	b.Lock()
	oldProp := b.Properties
	b.Properties = pnew
	b.Unlock()
	b.DM.DoNodeUpdateCallbacks(b, oldProp, pnew)
}

func (b *BaseNode) SetBaseProperties(pnew ifc.PropertyType, merge bool) {
	logging.Debugf("SetBaseProperties for %s prop = %s merge = %t", b.ID, pnew, merge)
	var toDelete []string
	if !merge {
		b.RLock()
		for itmOld := range b.Properties {
			if _, ok := pnew[itmOld]; !ok {
				if itmOld[0] != '_' && itmOld != b.NodeKeyName {
					toDelete = append(toDelete, itmOld)
				}
			}
		}

		b.RUnlock()
		if len(toDelete) != 0 {
			logging.Debugf("SetBaseProperties debug DEL ID=%s pnew=%s", b.ID, toDelete)
			b.DM.UpdateNodeRemoveProperties(b.ID, toDelete)
		}
	}
	b.DM.UpdateNodeAddProperties(b.ID, pnew)

	// UpdateNodeAddProperties does a force sync with the DB so we don't need to set it here.
	// if merge {
	// 	np := ifc.PropertyType{}
	// 	for itm, itmv := range b.Properties {
	// 		np[itm] = itmv
	// 	}
	// 	for itm, itmv := range pnew {
	// 		np[itm] = itmv
	// 	}
	// 	// What happens if the DM updates have already change the local node property
	// 	// should we be overwriting it?
	// 	// b.SetImmediateBaseProperties(np)
	// } else {
	// 	for itm, itmv := range b.Properties {
	// 		if itm[0] == '_' {
	// 			pnew[itm] = itmv
	// 		}
	// 	}
	// 	logging.Debugf("SetBaseProperties debug ADD ID=%s pnew=%s", b.ID, pnew)
	// 	// b.SetImmediateBaseProperties(pnew)
	// }
}

func (b *BaseNode) DelBaseProperties(keys []string) {
	logging.Debugf("DelBaseProperties for %s prop = %s", b.ID, keys)
	b.DM.UpdateNodeRemoveProperties(b.ID, keys)
	// b.SetImmediateBaseProperties(pnew)
}

func (b *BaseNode) DeleteImmediateLink(nodeType, nodeKey string) {
	olink, ok := b.Links.Get(nodeType, nodeKey)
	b.Links.Delete(nodeType, nodeKey)
	logging.Debugf("deleteImmediateLink: %s -> %s/%s %v/%v",
		b.ID, nodeType, nodeKey, ok, b.Links.Has(nodeType, nodeKey))
	if ok {
		b.DM.DoLinkDeleteCallbacks(b, nodeType, nodeKey, olink.Properties)
	}
}

func (b *BaseNode) DeleteImmediateRLink(nodeType, nodeKey string) {
	b.checkRLinkFeatureFlagStatus()
	/*
		We may not need link delete db for rlinks for now
	*/
	b.Links.DeleteRLink(nodeType, nodeKey)
	logging.Debugf("deleteImmediateRLink: %s -> %s/%s %v",
		b.ID, nodeType, nodeKey, b.Links.HasRLink(nodeType, nodeKey))
	/*
		We may not need link delete db for rlinks for now
	*/
}

func (b *BaseNode) IsRLinkFeatureFlagEnabled() bool {
	return b.DM.IsRLinkFeatureFlagEnabled()
}

// Creates rlink GLink object using source node GLink object
// Please note that rlink is always created in the db when a soft link is created
// However, the soft link prop could be updated by the creator DM or any other DM.
// Hence, we while generating GLink--> prop map--> created-by information
// we look for the presence of 'created-by' key in the softlink GLink obj
// and use the same if present
// Updated-by key is not generated in such cases
// Updated-by key === Created by key in soft link GLink prop
// Update time === creation time in soft link GLink prop
func createLinkObj(srcNodePath ifc.NodePathList, srcNodeType string, srcNodeKey string,
	srcNodeKeyName string, srcLnk *ifc.GLink) *ifc.GLink {
	rLnk := &ifc.GLink{
		Id:                srcLnk.DestinationNodeId + "/_rlinks/" + srcNodeType + "/" + srcNodeKey,
		LinkType:          srcNodeType,
		Properties:        map[string]interface{}{},
		SourceNodeId:      srcLnk.DestinationNodeId,
		DestinationNodeId: srcLnk.SourceNodeId,
	}

	if _, ok := srcLnk.Properties[common.LinkFixedProp_createdBy]; ok {
		rLnk.Properties[common.LinkFixedProp_createdBy] = srcLnk.Properties[common.LinkFixedProp_createdBy]
		rLnk.Properties[common.LinkFixedProp_updatedBy] = srcLnk.Properties[common.LinkFixedProp_createdBy]
	}

	if _, ok := srcLnk.Properties[common.LinkFixedProp_creationTime]; ok {
		rLnk.Properties[common.LinkFixedProp_creationTime] = srcLnk.Properties[common.LinkFixedProp_creationTime]
		rLnk.Properties[common.LinkFixedProp_updateTime] = srcLnk.Properties[common.LinkFixedProp_updateTime]
	}

	rLnk.Properties[common.LinkFixedProp_HardLink] = "false"
	rLnk.Properties[common.LinkFixedProp_NodeKeyName] = srcNodeKeyName
	rLnk.Properties[common.LinkFixedProp_NodeKeyValue] = srcNodeKey
	rLnk.Properties[common.LinkFixedProp_NodeType] = srcNodeType
	rLnk.Properties[common.LinkFixedProp_RSoftLinkDestinationPath] = srcNodePath
	rLnk.Properties[common.LinkFixedProp_destNodeId] = srcLnk.SourceNodeId
	return rLnk
}

// TODO: serialization is needed here. Goroutine is needed here instead of dataModelCache.go
func (b *BaseNode) UpsertImmediateLink(nodeType, nodeKey string, lnk *ifc.GLink) {
	olink, _ := b.Links.Get(nodeType, nodeKey)
	logging.Debugf("node Id %s adding link %s/%s lnk = %s", b.ID, nodeType, nodeKey, lnk.DestinationNodeId)
	b.Links.Add(nodeType, nodeKey, lnk)
	// Populate rlink for dest node if dest node exists in the cache
	// This is handled during link ntfn from the database
	// We dont process rlink ntfn, runtime processing is expensive, too many additional
	// cycles consumed
	// This is the focal point to handle rlink updates to cache via subscription ntfns
	// arriving for Links
	// Check only for soft Links
	// Wrap this op in a feature flag check for rlink
	if b.IsRLinkFeatureFlagEnabled() {
		if _, hLnkOk := lnk.Properties[common.LinkFixedProp_HardLink]; hLnkOk {
			if lnk.Properties[common.LinkFixedProp_HardLink] == "false" {
				if _, rLnkOk := lnk.Properties[common.LinkFixedProp_RSoftLinkDestinationPath]; !rLnkOk {
					if _, sLnkOk := lnk.Properties[common.LinkFixedProp_SoftLinkDestinationPath]; sLnkOk {
						dNode, dOk := b.DM.GetCachedNode(lnk.DestinationNodeId)
						if dNode != nil && dOk {
							rLnk := createLinkObj(b.GetFullPath(), b.GetType(), b.GetKeyValue(),
								b.GetKeyName(), lnk)
							logging.Debugf("Populating rlink here from %s--->%s.%s",
								lnk.DestinationNodeId, b.GetType(), b.GetKeyValue())
							dNode.UpsertImmediateRLink(b.GetType(), b.GetKeyValue(), rLnk)
						}
					}
				}
			}
		}
	}
	nkv := lnk.Properties[common.LinkFixedProp_NodeKeyValue]
	nt := lnk.Properties[common.LinkFixedProp_NodeType]
	b.DM.DoLinkUpdateCallbacks(
		b,
		nkv.(string),
		nt.(string),
		olink, lnk)
}

func (b *BaseNode) UpsertImmediateRLink(nodeType, nodeKey string, lnk *ifc.GLink) {
	b.checkRLinkFeatureFlagStatus()
	/*
		We may not need link delete db for rlinks for now
		olink, _ := b.Links.GetRLink(NodeType, nodeKey)
	*/
	logging.Debugf("Adding rlink from  %s --->  %s/%s rlnk dest node ID = %s",
		b.ID, nodeType, nodeKey, lnk.DestinationNodeId)
	b.Links.AddRLink(nodeType, nodeKey, lnk)
	/*
		We may not need link delete db for rlinks for now
		nkv := lnk.Properties[common.LinkFixedProp_NodeKeyValue]
		nt := lnk.Properties[common.LinkFixedProp_NodeType]
		b.DM.DoLinkUpdateCallbacks(
			b,
			nkv.(string),
			nt.(string),
			olink, lnk)
	*/
}

func (b *BaseNode) GetImmediateBaseProperties() ifc.PropertyType {
	b.RLock()
	defer b.RUnlock()
	return b.Properties
}

func (b *BaseNode) GetBaseProperties() (ifc.PropertyType, bool) {
	b.RLock()
	logging.Debugf("GetBaseProperties for %s. del=%t tobedel=%t", b.ID, b.Deleted, b.ToBeDeleted)
	/*
		Need a check for tobedeleted condition here
		Skip now to preserve backward compatibility
	*/
	if b.Deleted {
		b.RUnlock()
		return nil, false
	}
	b.RUnlock()
	b.DM.SyncNodeProperty(b, false)
	b.RLock()
	defer b.RUnlock()
	logging.Debugf("GetBaseProperties for %s returning prop = %s", b.ID, b.Properties)
	return b.Properties, true
}

func (b *BaseNode) LinkAddProperty(childKey, nodeType string, linkProp ifc.PropertyType) bool {
	/*
		Need a check for tobedeleted condition here
		Skip now to preserve backward compatibility
	*/

	if b.IsDeleted() {
		return false
	}
	return b.DM.UpdateLinkAddProperty(b.ID, nodeType, childKey, linkProp)
}

func (b *BaseNode) LinkRemoveProperty(childKey, nodeType string, linkPropKeys []string) bool {
	/*
		Need a check for tobedeleted condition here
		Skip now to preserve backward compatibility
	*/

	if b.IsDeleted() {
		return false
	}
	return b.DM.UpdateLinkRemoveProperty(b.ID, nodeType, childKey, linkPropKeys)
}

func (b *BaseNode) GetChild(childKey, nodeType string, forcedSync bool) (ifc.BaseNodeInterface, ifc.PropertyType, bool) {
	if childKey == "" || nodeType == "" {
		logging.Fatalf("invalid arguments key=%s nodeType=%s", childKey, nodeType)
	}
	logging.Debugf("GetChild called for node=%s and child=%s.%s force = %v", b.ID, nodeType, childKey, forcedSync)
	b.DM.SyncChild(b, nodeType, childKey, forcedSync)
	lnk, lnkOk := b.Links.Get(nodeType, childKey)
	if !lnkOk {
		logging.Debugf("GetChild returning nil")
		return nil, nil, false
	}
	//logging.Debugf("GetChild2 called for node=%s and child=%s.%s destination=%s", b.ID, nodeType, childKey, lnk.DestinationNodeId)

	r, rOk := b.DM.GetNode(b.ID, nodeType, lnk.DestinationNodeId, false, false)
	if !rOk {
		return nil, nil, false
		// panic(errors.New("Unexpected error in GetChild"))
	}
	logging.Debugf("GetChild called for %s.%s Returning %v %v lnkProp=%s", nodeType, childKey, forcedSync, rOk, lnk.Properties)
	return r, lnk.Properties, true
}

func (b *BaseNode) GetCachedChild(childKey, nodeType string) (ifc.BaseNodeInterface, ifc.PropertyType, bool) {
	if childKey == "" || nodeType == "" {
		logging.Fatalf("invalid arguments key=%s NodeType=%s", childKey, nodeType)
	}
	lnk, lnkOk := b.Links.Get(nodeType, childKey)
	if !lnkOk {
		return nil, nil, false
	}
	r, rOk := b.DM.GetCachedNode(lnk.DestinationNodeId)
	if !rOk {
		return nil, nil, false
		// panic(errors.New("Unexpected error in GetCachedChild"))
	}
	return r, lnk.Properties, true
}

func (b *BaseNode) GetSoftLinkedChildLeg(childKey, nodeType string) []ifc.BaseNodeInterface {
	if childKey == "" || nodeType == "" {
		logging.Fatalf("invalid arguments key=%s NodeType=%s", childKey, nodeType)
	}
	b.DM.SyncChild(b, nodeType, childKey, false)
	lnk, lnkOk := b.Links.Get(nodeType, childKey)
	var r []ifc.BaseNodeInterface
	if !lnkOk {
		return r
	}
	lt := lnk.Properties[common.LinkFixedProp_HardLink]
	if lt.(string) == "true" {
		panic(errors.New("soft link method called for a hard link object"))
	}
	var destNodePath ifc.NodePathList
	byt := []byte((lnk.Properties[common.LinkFixedProp_SoftLinkDestinationPath]).(string))
	if err := json.Unmarshal(byt, &destNodePath); err != nil {
		panic(err)
	}
	return b.DM.PopulatePathAndFetchNodes(destNodePath, false)
}

func (b *BaseNode) GetChildLink(childKey, nodeType string) (ifc.PropertyType, bool) {
	if childKey == "" || nodeType == "" {
		logging.Fatalf("invalid arguments key=%s NodeType=%s", childKey, nodeType)
	}
	b.DM.SyncChild(b, nodeType, childKey, false)
	if ln, lnOk := b.Links.Get(nodeType, childKey); !lnOk {
		return nil, false
	} else {
		return ln.Properties, true
	}
}

func (b *BaseNode) getCachedRLinkChildKeyList(nodeType string) []ifc.BaseNodeLinkManagerIteratorInterface {
	b.checkRLinkFeatureFlagStatus()
	var rLinkItrList []ifc.BaseNodeLinkManagerIteratorInterface
	var ln []*ifc.GLink
	var lnOk bool
	if nodeType != "" {
		ln, lnOk = b.Links.GetRLinksForNodeType(nodeType)
	} else {
		ln, lnOk = b.Links.GetRLinks()
	}
	if !lnOk {
		return rLinkItrList
	}
	for _, lnkObj := range ln {
		if dNode, ok := b.DM.GetCachedNode(lnkObj.DestinationNodeId); ok {
			lnkItrIntf := linkItr.NewBaseNodeLinkManagerIterator(dNode, lnkObj)
			if lnkItrIntf != nil {
				rLinkItrList = append(rLinkItrList, lnkItrIntf)
				continue
			}
			logging.Debugf("RLink iterator handler itr is nil for link ID: %v", lnkObj.DestinationNodeId)
		}
	}
	return rLinkItrList
}

func (b *BaseNode) GetCachedRLinkHdlr(nodeType, childKey string) ifc.BaseNodeLinkManagerIteratorInterface {
	b.checkRLinkFeatureFlagStatus()
	if childKey == "" || nodeType == "" {
		logging.Fatalf("invalid arguments for rlnk cache read key=%s nodeType=%s", childKey, nodeType)
	}
	if ln, lnOk := b.Links.GetRLink(nodeType, childKey); !lnOk {
		return nil
	} else {
		if dNode, ok := b.DM.GetCachedNode(ln.DestinationNodeId); ok {
			return linkItr.NewBaseNodeLinkManagerIterator(dNode, ln)
		}
	}
	return nil
}

func (b *BaseNode) ForEachCachedRLinkWithNodeType(fn func(rLinkHdlr ifc.BaseNodeLinkManagerIteratorInterface) interface{},
	nodeType string) {
	b.checkRLinkFeatureFlagStatus()
	rLinkHdlrList := b.getCachedRLinkChildKeyList(nodeType)
	if len(rLinkHdlrList) == 0 {
		return
	}
	for _, hdlr := range rLinkHdlrList {
		fn(hdlr)
	}
}

func (b *BaseNode) ForEachCachedRLink(fn func(rLinkHdlr ifc.BaseNodeLinkManagerIteratorInterface) interface{}) {
	b.checkRLinkFeatureFlagStatus()
	rLinkHdlrList := b.getCachedRLinkChildKeyList("")
	if len(rLinkHdlrList) == 0 {
		return
	}
	for _, hdlr := range rLinkHdlrList {
		fn(hdlr)
	}
}

func (b *BaseNode) getCachedChildLinkChildKeyList(nodeType string) []ifc.BaseNodeLinkManagerIteratorInterface {
	var hLinkItrList []ifc.BaseNodeLinkManagerIteratorInterface
	var ln []*ifc.GLink
	var lnOk bool
	if nodeType != "" {
		ln, lnOk = b.Links.GetChildLinksForNodeType(nodeType)
	} else {
		ln, lnOk = b.Links.GetChildLinks()
	}
	if !lnOk {
		return hLinkItrList
	}
	for _, lnkObj := range ln {
		if dNode, ok := b.DM.GetCachedNode(lnkObj.DestinationNodeId); ok {
			lnkItrIntf := linkItr.NewBaseNodeLinkManagerIterator(dNode, lnkObj)
			if lnkItrIntf != nil {
				hLinkItrList = append(hLinkItrList, lnkItrIntf)
				continue
			}
			logging.Debugf("RLink iterator handler itr is nil for link ID: %v", lnkObj.DestinationNodeId)
		}
	}
	return hLinkItrList
}

func (b *BaseNode) GetCachedChildLinkHdlr(nodeType, childKey string) ifc.BaseNodeLinkManagerIteratorInterface {
	if childKey == "" || nodeType == "" {
		logging.Fatalf("invalid arguments for rlnk cache read key=%s nodeType=%s", childKey, nodeType)
	}
	if ln, lnOk := b.Links.GetChildLink(nodeType, childKey); !lnOk {
		return nil
	} else {
		if dNode, ok := b.DM.GetCachedNode(ln.DestinationNodeId); ok {
			return linkItr.NewBaseNodeLinkManagerIterator(dNode, ln)
		}
	}
	return nil
}

func (b *BaseNode) ForEachCachedChildLinkWithNodeType(fn func(hLinkHdlr ifc.BaseNodeLinkManagerIteratorInterface) interface{},
	nodeType string) {
	hLinkHdlrList := b.getCachedChildLinkChildKeyList(nodeType)
	if len(hLinkHdlrList) == 0 {
		return
	}
	for _, hdlr := range hLinkHdlrList {
		fn(hdlr)
	}
}

func (b *BaseNode) ForEachCachedChildLink(fn func(hLinkHdlr ifc.BaseNodeLinkManagerIteratorInterface) interface{}) {
	hLinkHdlrList := b.getCachedChildLinkChildKeyList("")
	if len(hLinkHdlrList) == 0 {
		return
	}
	for _, hdlr := range hLinkHdlrList {
		fn(hdlr)
	}
}

func (b *BaseNode) SetBaseParent(n ifc.BaseNodeInterface) {
	if n == b {
		logging.Fatalf("Recursion is not allowed when setting base node parent")
	}
	b.Lock()
	b.BaseParent = n
	b.Unlock()
}
func (b *BaseNode) GetNextChildKeyList(nodeName string, nodeType string, cnt uint32) []string {
	return b.DM.GetNextChildKeyList(b.ID, nodeType, nodeName, cnt)
}
func (b *BaseNode) GetNextRLinkChildKeyList(nodeName string, nodeType string, cnt uint32) []string {
	b.checkRLinkFeatureFlagStatus()
	return b.DM.GetNextRLinkChildKeyList(b.ID, nodeType, nodeName, cnt)
}

// DEPRECATED. Use GetNextChildKeyList instead.
func (b *BaseNode) GetNextChildKey(nodeName, nodeType string) string {
	if nodeType == "" {
		logging.Fatalf("invalid arguments NodeType=%s", nodeType)
	}
	if nodeName == "" {
		go func() {
			b.Sync(false)
		}()
	}
	return b.Links.GetNextKey(nodeType, nodeName)
	/*
	       r := b.GetNextChildKeyList(nodeName, NodeType, 1)
	   	if len(r) > 0 {
	   		return r[0]
	   	} else {
	   		return ""
	   	}*/
}
func (b *BaseNode) UpsertChild(nodeType, nodeKeyName string, nodeProp ifc.PropertyType, linkType string,
	linkProp ifc.PropertyType) ifc.BaseNodeInterface {
	if nodeKeyName == "" || nodeType == "" {
		logging.Fatalf("invalid arguments key=%s nodeType=%s", nodeKeyName, nodeType)
	}
	logging.Debugf("UpsertChild for node=%s child=%s.%s", b.ID, nodeType, nodeKeyName)
	newNode, _ := b.DM.UpsertNode(b.ID, linkType, nodeKeyName,
		linkProp, nodeType, nodeProp)
	return newNode
}

func (b *BaseNode) UpsertLink(destNode ifc.BaseNodeInterface, linkType string, linkProp ifc.PropertyType,
	isSingleton bool, rootName string, useUUIDasKey bool) *ifc.GLink {
	dkv := destNode.GetKeyValue()
	if useUUIDasKey {

		// We don't want to use UUID here since the calls won't be idempotent.
		// Rely on the full path of the destination node instead which is unique for us.
		var path []string
		for _, outerPart := range destNode.GetFullPath() {
			for _, innerPart := range outerPart {
				path = append(path, innerPart)
			}
		}

		dkv = strings.Join(path, ":")
	}
	if isSingleton {
		dkv = common.NodeFixedProp_NodeSingletonKeyValue
	}
	srcPathA := b.GetFullPath()
	srcPathA[0][ifc.NodePathName_nodeIdentifier] = rootName
	srcPath := utils.JsonMarshal(srcPathA)

	destPathA := destNode.GetFullPath()
	destPathA[0][ifc.NodePathName_nodeIdentifier] = rootName
	destPath := utils.JsonMarshal(destPathA)

	lnkOpObj := common.UpsertLinkOpObj{LinkType: common.LinkType_Owner, SrcNodeObj: common.UpsertLinkOpNodeObj{NodeId: b.ID, NodePath: string(srcPath),
		NodeType: b.GetType(), NodeKeyValue: b.GetKeyValue(), NodeProp: b.GetImmediateBaseProperties()},
		DestNodeObj: common.UpsertLinkOpNodeObj{NodeId: destNode.GetId(), NodePath: string(destPath),
			NodeType: destNode.GetType(), NodeKeyValue: dkv, NodeProp: destNode.GetImmediateBaseProperties()},
		IsSnglton: isSingleton, LinkPropIn: linkProp}

	lnk := b.DM.UpsertLink(lnkOpObj)
	destNode.AddReverseLink(b.ID, dkv, lnk.Id)
	return lnk
}

func (b *BaseNode) DeleteLink(childKey, nodeType string) bool {
	cnt := 0
	linkedNodes := b.GetSoftLinkedChildLeg(childKey, nodeType)
	for _, linkedNode := range linkedNodes {
		cnt++
		linkedNode.RemoveReverseLink(b.ID)
	}
	b.DM.DeleteLink(b.ID, nodeType, childKey)
	return cnt != 0
}

func (b *BaseNode) DeleteLinkToNode(destNode ifc.BaseNodeInterface, nodeType string) bool {
	if destKey, ok := b.Links.Find(nodeType, func(links *ifc.GLink) bool {
		return links.DestinationNodeId == destNode.GetId()
	}); ok {
		return b.DeleteLink(destKey, nodeType)
	}
	fmt.Printf("DeleteLinkToNode unable to find destination node with ID =%s\n", destNode.GetId())
	return false
}

func (b *BaseNode) Delete() {
	b.RLock()
	pn := b.BaseParent
	nType := b.NodeType
	nValue := b.Properties[b.NodeKeyName].(string)
	b.RUnlock()
	b.SetToBeDeleted()
	b.DM.Delete(b)
	b.SetDeleted()
	if pn != nil && pn.GetLinks().Has(nType, nValue) {
		b.DM.DeleteLink(pn.GetId(), nType, nValue)
	}
}

func (b *BaseNode) IsTypeRefNode() bool {
	b.RLock()
	defer b.RUnlock()
	return b.TypeRefNode
}
func (b *BaseNode) IsDeleted() bool {
	b.RLock()
	defer b.RUnlock()
	return b.Deleted
}
func (b *BaseNode) SetDeleted() {
	b.Lock()
	b.Deleted = true
	b.Unlock()
}
func (b *BaseNode) IsToBeDeleted() bool {
	b.RLock()
	defer b.RUnlock()
	return b.ToBeDeleted
}
func (b *BaseNode) SetToBeDeleted() {
	b.Lock()
	b.ToBeDeleted = true
	b.Unlock()
}

func (b *BaseNode) SetDeleteCompleted() {
	b.Lock()
	b.DeleteCompleted = true
	b.Unlock()
}
func (b *BaseNode) IsDeleteCompleted() bool {
	b.RLock()
	defer b.RUnlock()
	return b.DeleteCompleted
}
func (b *BaseNode) GetKeyName() string {
	b.RLock()
	defer b.RUnlock()
	return b.NodeKeyName
}
func (b *BaseNode) SetKeyName(n string) {
	b.Lock()
	b.NodeKeyName = n
	b.Unlock()
}
func (b *BaseNode) SetKeyValue(v string) {
	b.Lock()
	b.Properties[b.NodeKeyName] = v
	b.Unlock()
}
func (b *BaseNode) GetId() string {
	b.RLock()
	defer b.RUnlock()
	return b.ID
}
func (b *BaseNode) GetType() string {
	b.RLock()
	defer b.RUnlock()
	return b.NodeType
}
func (b *BaseNode) SetType(t string) {
	b.Lock()
	b.NodeType = t
	b.Unlock()
}
func (b *BaseNode) AddReverseLink(nodeId, linkKeyValue, linkId string) {
	b.Lock()
	defer b.Unlock()
	logging.Debugf("%s:reverse link ADD from %s -> %s \n", b.DM.GetName(), b.ID, nodeId)
	b.RSoftLink[nodeId] = ReverseLink{NodeKey: linkKeyValue, LinkId: linkId}
}
func (b *BaseNode) RemoveReverseLink(nodeId string) {
	b.Lock()
	defer b.Unlock()
	logging.Debugf("%s:reverse link DELETE from %s -> %s \n", b.DM.GetName(), b.ID, nodeId)
	delete(b.RSoftLink, nodeId)
}
func (b *BaseNode) ReverseLinkIterate(fn func(linkKeyValue, nodeId, linkId string)) {
	b.RLock()
	defer b.RUnlock()
	for k, v := range b.RSoftLink {
		fn(v.NodeKey, k, v.LinkId)
	}
}
func (b *BaseNode) GetKeyValue() string {
	b.RLock()
	defer b.RUnlock()
	v, ok := b.Properties[b.NodeKeyName]
	if !ok {
		logging.Infof("%s:basenode GetKeyValue ID=%s[%t %t %t] keyname = %s Property %s",
			b.DM.GetName(), b.ID, b.ToBeDeleted, b.Deleted, b.DeleteCompleted, b.NodeKeyName, b.Properties)
		logging.Debugf("--------------FULL Stack DUMP---------------------\n")
		buf := make([]byte, 1<<16)
		runtime.Stack(buf, true)
		logging.Debugf("%s", buf)
		logging.Debugf("-----------------------------------------\n")
		panic(errors.New("internal error, nodekeyname is not set in Properties of the node"))
	}

	return v.(string)
}

func (b *BaseNode) GetBaseParent() ifc.BaseNodeInterface {
	b.RLock()
	defer b.RUnlock()
	return b.BaseParent
}

func (b *BaseNode) GetDBStatsString() string {
	//	dbStats := b.DM.GetDBStats()
	//	var dbStr []string
	//	for k, v := range dbStats {
	//		dbStr = append(dbStr, k+":"+v)
	//	}
	//	return strings.Join(dbStr, ", ")
	return ""
}
func (b *BaseNode) GetMSGStatsString() string {
	//	msgStats := b.DM.GetMsgStats()
	//	var msgStr []string
	//	for k, v := range msgStats {
	//		msgStr = append(msgStr, k+":"+v)
	//	}
	//	return strings.Join(msgStr, ", ")
	return ""
}

func (b *BaseNode) Checksub() {

}

var decodeHooks = []mapstructure.DecodeHookFunc{}

func SetGlobalDecoders(decoders []mapstructure.DecodeHookFunc) {
	decodeHooks = decoders
}

func MapstructureDecode(input interface{}, result interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:   nil,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(decodeHooks...),
		Result:     result,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}

func MarshalProperties(propMap map[string]interface{}, jsonKeys []string, protoKeys map[string]reflect.Type) error {
	for _, jsonKey := range jsonKeys {
		if _, ok := propMap[jsonKey]; ok {
			jsonVal, err := json.Marshal(propMap[jsonKey])
			if err != nil {
				return err
			}
			propMap[jsonKey] = url.QueryEscape(string(jsonVal))
		}
	}
	for protoKey := range protoKeys {
		if _, ok := propMap[protoKey]; ok {
			if propMap[protoKey] == nil || reflect.ValueOf(propMap[protoKey]).IsNil() {
				propMap[protoKey] = ""
				continue
			}
			protoVal, err := proto.Marshal(propMap[protoKey].(proto.Message))
			if err != nil {
				return err
			}
			propMap[protoKey] = base64.StdEncoding.EncodeToString([]byte(protoVal))
		}
	}
	return nil
}

func UnmarshalProperties(propMap map[string]interface{}, jsonKeys []string, protoKeys map[string]reflect.Type) error {
	for _, jsonKey := range jsonKeys {
		if jsonVal, ok := propMap[jsonKey]; ok {
			dec, err := url.QueryUnescape(jsonVal.(string))
			if err != nil {
				return err
			}
			var m interface{}
			err = json.Unmarshal([]byte(dec), &m)
			if err != nil {
				if jsonVal != "undefined" && jsonVal != "" {
					return err
				}
			}
			propMap[jsonKey] = m
		}
	}
	for protoKey, protoType := range protoKeys {
		if protoVal, ok := propMap[protoKey]; ok {
			dec, err := base64.StdEncoding.DecodeString(protoVal.(string))
			if err != nil {
				return err
			}
			m := reflect.New(protoType).Interface().(proto.Message)
			err = proto.Unmarshal([]byte(dec), m)
			if err != nil {
				if protoVal != "undefined" && protoVal != "" {
					return err
				}
			}
			propMap[protoKey] = m
		}
	}
	return nil
}
func Clone(src map[string]interface{}) map[string]interface{} {
	if src == nil {
		return nil
	}
	dst := make(map[string]interface{})
	jsonStr, err := json.Marshal(src)
	if err != nil {
		logging.Fatalf("Error when Cloning from src: %s", err)
	}
	err = json.Unmarshal(jsonStr, &dst)
	if err != nil {
		logging.Fatalf("Error when Cloning to dest: %s", err)
	}
	return dst
}
