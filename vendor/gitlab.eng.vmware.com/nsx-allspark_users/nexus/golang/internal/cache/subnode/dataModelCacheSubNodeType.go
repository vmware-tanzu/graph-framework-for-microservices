package subnode

import (
	"sync"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/internal/base/link"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/ifc"
)

// CachedNodeData Status of the node
type CachedNodeData struct {
	NodeRequested bool // node is requested from database
	NodeInCache   bool // node is in cache
}

// type CallbackFunc func()

// SubNodeType Data structure to hold
// all the subscription information and the callback functions
type SubNodeType struct {
	sync.RWMutex

	Weight        uint32                  `json:"weight"`
	Depth         uint32                  `json:"depth"`
	CbfnNode      []*ifc.CallbackFuncNode `json:"-"`
	CbfnNodeMutex sync.RWMutex            `json:"-"`

	CbfnLink      map[string][]*ifc.CallbackFuncLink `json:"-"`
	CbfnLinkMutex sync.RWMutex                       `json:"-"`

	CachedNodes      map[string]*CachedNodeData `json:"cached_nodes"`
	CachedNodesMutex sync.RWMutex               `json:"-"`

	Child      map[string](map[string]*SubNodeType) `json:"child"`
	ChildMutex sync.RWMutex                         `json:"-"`
}

func NewSubNodeType() *SubNodeType {
	e := &SubNodeType{
		Weight:      0,
		Depth:       0,
		CbfnNode:    []*ifc.CallbackFuncNode{},
		CbfnLink:    make(map[string][]*ifc.CallbackFuncLink),
		CachedNodes: make(map[string]*CachedNodeData),
		Child:       make(map[string]map[string]*SubNodeType)}
	return e
}

// IncChildWeight increase the Weight for the node
func (snt *SubNodeType) IncChildWeight(nType, nValue string) []string {
	snt.ChildMutex.Lock()
	defer snt.ChildMutex.Unlock()
	ndt, ntdok := snt.Child[nType]
	if !ntdok {
		ndt = make(map[string]*SubNodeType)
		snt.Child[nType] = ndt
	}
	vdt, vdtok := ndt[nValue]
	if !vdtok {
		st := NewSubNodeType()
		ndt[nValue] = st
		st.InitWeight()
		return []string{}
	}
	return vdt.IncWeight()
}

// GetChildKeyForType get all the value keys for a give type
func (snt *SubNodeType) GetChildKeyForType(nType string) []string {
	snt.ChildMutex.RLock()
	defer snt.ChildMutex.RUnlock()
	dt, dtok := snt.Child[nType]
	link.NewBaseNodeLinkManager()
	keys := []string{}
	if dtok {
		for key := range dt {
			keys = append(keys, key)
		}
	}
	return keys
}

// ChildForEach iterate through all the children
func (snt *SubNodeType) ChildForEach(fn func(nType, nValue string, nd *SubNodeType)) {
	snt.ChildMutex.RLock()
	defer snt.ChildMutex.RUnlock()
	for key1, val2 := range snt.Child {
		for key2, val2 := range val2 {
			fn(key1, key2, val2)
		}
	}
}
func (snt *SubNodeType) DelChild(nType, nValue string) {
	snt.ChildMutex.Lock()
	defer snt.ChildMutex.Unlock()
	nt, ntok := snt.Child[nType]
	if !ntok {
		return
	}
	delete(nt, nValue)
	if len(nt) == 0 {
		delete(snt.Child, nType)
	}
}

// GetChild get a Child node if it exists
func (snt *SubNodeType) GetChild(nType, nValue string) *SubNodeType {
	snt.ChildMutex.RLock()
	defer snt.ChildMutex.RUnlock()
	if ndt, ok := snt.Child[nType]; ok {
		if ndv, ok := ndt[nValue]; ok {
			return ndv
		}
	}
	return nil
}

// IsTypeInChild check if the nType is valid
func (snt *SubNodeType) IsTypeInChild(nType string) bool {
	snt.ChildMutex.RLock()
	defer snt.ChildMutex.RUnlock()
	_, ndtok := snt.Child[nType]
	return ndtok
}

// CbfnLinkHasNodeType check if nodeType is present in link callback
func (snt *SubNodeType) CbfnLinkHasNodeType(nodeType string) bool {
	snt.CbfnLinkMutex.RLock()
	defer snt.CbfnLinkMutex.RUnlock()
	_, ndtok := snt.CbfnLink[nodeType]
	return ndtok
}

// CbfnLinkLength size of the CbfnLink Map
func (snt *SubNodeType) CbfnLinkLength() int {
	snt.CbfnLinkMutex.RLock()
	defer snt.CbfnLinkMutex.RUnlock()
	return len(snt.CbfnLink)
}

// CbfnLinkForEach callback fn for each link for a nodetype
func (snt *SubNodeType) CbfnLinkForEach(nodeType string, fn func(arg *ifc.CallbackFuncLink)) {
	snt.CbfnLinkMutex.RLock()
	defer snt.CbfnLinkMutex.RUnlock()
	nd, ndok := snt.CbfnLink[nodeType]
	if ndok {
		for _, val := range nd {
			fn(val)
		}
	}
}

// AddCachedNode add a cached node
func (snt *SubNodeType) AddCachedNode(id string, dt *CachedNodeData) {
	snt.CachedNodesMutex.Lock()
	defer snt.CachedNodesMutex.Unlock()
	snt.CachedNodes[id] = dt
}

// DelCachedNode ached node
func (snt *SubNodeType) DelCachedNode(id string) {
	snt.CachedNodesMutex.Lock()
	defer snt.CachedNodesMutex.Unlock()
	delete(snt.CachedNodes, id)
}

// CachedNodeForEach iterate over all cached nodes
func (snt *SubNodeType) CachedNodeForEach(fn func(v *CachedNodeData, k string)) {
	snt.CachedNodesMutex.RLock()
	defer snt.CachedNodesMutex.RUnlock()
	for key, val := range snt.CachedNodes {
		fn(val, key)
	}
}

// GetCachedNodeKeys get all the nodeid's for cached nodes
func (snt *SubNodeType) GetCachedNodeKeys() []string {
	snt.CachedNodesMutex.RLock()
	defer snt.CachedNodesMutex.RUnlock()
	keys := []string{}
	for key := range snt.CachedNodes {
		keys = append(keys, key)
	}
	return keys
}

// IsEmpty is the ds empty
func (snt *SubNodeType) IsEmpty() bool {
	snt.CachedNodesMutex.RLock()
	lenCachedNodes := len(snt.CachedNodes)
	snt.CachedNodesMutex.RUnlock()
	snt.ChildMutex.RLock()
	lenChild := len(snt.Child)
	snt.ChildMutex.RUnlock()
	snt.RLock()
	weight := snt.Weight
	depth := snt.Depth
	snt.RUnlock()
	return (weight == 0 &&
		snt.GetCbfnNodeLength() == 0 &&
		snt.CbfnLinkLength() == 0 &&
		depth == 0 &&
		lenCachedNodes == 0 &&
		lenChild == 0)
}

// GetWeight returns the current Weight
func (snt *SubNodeType) GetWeight() uint32 {
	snt.RLock()
	defer snt.RUnlock()
	return snt.Weight
}

// InitWeight initialize the Weight
func (snt *SubNodeType) InitWeight() {
	snt.Lock()
	snt.Weight = 1
	snt.Unlock()
}

// GetDepth get the deopt
func (snt *SubNodeType) GetDepth() uint32 {
	snt.RLock()
	defer snt.RUnlock()
	return snt.Depth
}

//SetDepth set the Depth
func (snt *SubNodeType) SetDepth(n uint32) {
	snt.Lock()
	snt.Depth = n
	snt.Unlock()
}

// GetCbfnNodeLength return the CbfnNode length
func (snt *SubNodeType) GetCbfnNodeLength() int {
	snt.CbfnNodeMutex.RLock()
	defer snt.CbfnNodeMutex.RUnlock()
	return len(snt.CbfnNode)
}

// CbfnNodeForEach foreach of the call back function for this node
func (snt *SubNodeType) CbfnNodeForEach(fn func(f *ifc.CallbackFuncNode)) {
	snt.CbfnNodeMutex.RLock()
	defer snt.CbfnNodeMutex.RUnlock()
	for _, v := range snt.CbfnNode {
		fn(v)
	}
}

// Add add a new node if it's not already there
func (snt *SubNodeType) Add(nType, nValue string) {
	snt.ChildMutex.Lock()
	defer snt.ChildMutex.Unlock()
	tk, tkok := snt.Child[nType]
	if !tkok {
		tk = make(map[string]*SubNodeType)
		snt.Child[nType] = tk
	}
	_, vkok := tk[nValue]
	if !vkok {
		tk[nValue] = NewSubNodeType()
	}
}

// IncWeight for this node
func (snt *SubNodeType) IncWeight() []string {
	snt.Lock()
	if snt.Weight == 0 {
		snt.Weight = 1
		snt.Unlock()

		keys := []string{}
		snt.CachedNodesMutex.Lock()
		for key := range snt.CachedNodes { // /Root/root
			keys = append(keys, key)
		}

		snt.CachedNodesMutex.Unlock()
		return keys
	}

	snt.Weight = snt.Weight + 1
	snt.Unlock()
	return []string{}

}

// DecWeight decrease the Weight
func (snt *SubNodeType) DecWeight() {
	if snt.Weight > 0 {
		snt.Weight = snt.Weight - 1
	}
}

// AddCB add a link/node call back function to tracked nod
// func (snt *SubNodeType) AddCB(cbfn CallbackFunc, isLinkCB bool, destNodeType string) {
func (snt *SubNodeType) AddCBLink(cbfn *ifc.CallbackFuncLink, destNodeType string) {
	snt.CbfnLinkMutex.Lock()
	defer snt.CbfnLinkMutex.Unlock()
	lnd, lndok := snt.CbfnLink[destNodeType]
	if !lndok {
		lnd = []*ifc.CallbackFuncLink{}
		snt.CbfnLink[destNodeType] = lnd
	}
	snt.CbfnLink[destNodeType] = append(lnd, cbfn)
}
func (snt *SubNodeType) AddCBNode(cbfn *ifc.CallbackFuncNode) {
	snt.CbfnNodeMutex.Lock()
	defer snt.CbfnNodeMutex.Unlock()
	snt.CbfnNode = append(snt.CbfnNode, cbfn)
}
