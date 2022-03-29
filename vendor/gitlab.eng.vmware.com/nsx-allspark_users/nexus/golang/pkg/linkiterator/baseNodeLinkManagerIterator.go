package linkiterator

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/ifc"
)

type BaseNodeLinkManagerIterator struct {
	basenodeIntf ifc.BaseNodeInterface
	lnk          *ifc.GLink
}

func NewBaseNodeLinkManagerIterator(bNodeIntf ifc.BaseNodeInterface, linkIn *ifc.GLink) *BaseNodeLinkManagerIterator {
	e := &BaseNodeLinkManagerIterator{basenodeIntf: bNodeIntf, lnk: linkIn}
	return e
}

func (lmi *BaseNodeLinkManagerIterator) GetNodeId() string {
	return lmi.basenodeIntf.GetId()
}

func (lmi *BaseNodeLinkManagerIterator) GetNodeType() string {
	return lmi.basenodeIntf.GetType()
}

func (lmi *BaseNodeLinkManagerIterator) GetNodeKeyValue() string {
	return lmi.basenodeIntf.GetKeyValue()
}

func (lmi *BaseNodeLinkManagerIterator) GetParentNodeIterator() ifc.BaseNodeLinkManagerIteratorInterface {
	parNodeIntf := lmi.basenodeIntf.GetBaseParent()
	if parNodeIntf != nil {
		return NewBaseNodeLinkManagerIterator(parNodeIntf, nil)
	}
	return nil
}

func (lmi *BaseNodeLinkManagerIterator) GetNodeParentNodeKey() string {
	parNodeIntf := lmi.basenodeIntf.GetBaseParent()
	if parNodeIntf != nil {
		return parNodeIntf.GetKeyValue()
	}
	return ""
}

func (lmi *BaseNodeLinkManagerIterator) GetNodeParentNodeType() string {
	parNodeIntf := lmi.basenodeIntf.GetBaseParent()
	if parNodeIntf != nil {
		return parNodeIntf.GetType()
	}
	return ""
}

func (lmi *BaseNodeLinkManagerIterator) ForEachCachedRLinkWithNodeType(fn func(rLinkHdlr ifc.BaseNodeLinkManagerIteratorInterface) interface{},
	nodeType string) {
	lmi.basenodeIntf.ForEachCachedRLinkWithNodeType(fn, nodeType)
}

func (lmi *BaseNodeLinkManagerIterator) ForEachCachedRLink(fn func(rLinkHdlr ifc.BaseNodeLinkManagerIteratorInterface) interface{}) {
	lmi.basenodeIntf.ForEachCachedRLink(fn)
}

func (lmi *BaseNodeLinkManagerIterator) GetNodeProperties() ifc.PropertyType {
	return lmi.basenodeIntf.GetImmediateBaseProperties()
}

func (lmi *BaseNodeLinkManagerIterator) GetCachedRLinkHdlr(nodeType string,
	nodeKey string) ifc.BaseNodeLinkManagerIteratorInterface {
	return lmi.basenodeIntf.GetCachedRLinkHdlr(nodeType, nodeKey)
}

func (lmi *BaseNodeLinkManagerIterator) IsRLinkFeatureFlagEnabled() bool {
	return lmi.basenodeIntf.IsRLinkFeatureFlagEnabled()
}

func (lmi *BaseNodeLinkManagerIterator) GetLinkProperties() (ifc.PropertyType, bool) {
	if lmi.lnk != nil {
		return lmi.lnk.Properties, true
	}
	return ifc.PropertyType{}, false
}

func (lmi *BaseNodeLinkManagerIterator) ForEachCachedChildLinkWithNodeType(fn func(hLinkHdlr ifc.BaseNodeLinkManagerIteratorInterface) interface{},
	nodeType string) {
	lmi.basenodeIntf.ForEachCachedChildLinkWithNodeType(fn, nodeType)
}

func (lmi *BaseNodeLinkManagerIterator) ForEachCachedChildLink(fn func(hLinkHdlr ifc.BaseNodeLinkManagerIteratorInterface) interface{}) {
	lmi.basenodeIntf.ForEachCachedChildLink(fn)
}

func (lmi *BaseNodeLinkManagerIterator) GetCachedChildLinkHdlr(nodeType string,
	nodeKey string) ifc.BaseNodeLinkManagerIteratorInterface {
	return lmi.basenodeIntf.GetCachedChildLinkHdlr(nodeType, nodeKey)
}
