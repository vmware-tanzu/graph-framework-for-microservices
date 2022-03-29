package link

import (
	"sort"
	"sync"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/ifc"
)

type BaseNodeLinkManager struct {
	Links  map[string](map[string]*ifc.GLink)
	RLinks map[string](map[string]*ifc.GLink)
	sync.RWMutex
}

func NewBaseNodeLinkManager() *BaseNodeLinkManager {
	e := &BaseNodeLinkManager{
		Links:  make(map[string](map[string]*ifc.GLink)),
		RLinks: make(map[string](map[string]*ifc.GLink))}
	return e
}

func (lm *BaseNodeLinkManager) AddRLink(nodeType, nodeKey string, lnk *ifc.GLink) {
	ch := lm.UpsertRLinkType(nodeType)
	lm.Lock()
	ch[nodeKey] = lnk
	lm.Unlock()
}

func (lm *BaseNodeLinkManager) GetNextKey(nodeType, key string) string {
	lm.RLock()
	defer lm.RUnlock()
	l, lok := lm.Links[nodeType]
	if !lok {
		return ""
	}
	keys := make([]string, len(l))
	cnt := 0
	for k := range l {
		keys[cnt] = k
		cnt++
	}
	sort.Strings(keys)
	pkey := ""
	for _, k := range keys {

		if key == pkey {
			return k
		}
		pkey = k
	}
	return ""
}
func (lm *BaseNodeLinkManager) UpsertType(nodeType string) map[string]*ifc.GLink {
	lm.Lock()
	defer lm.Unlock()
	if ch, ok := lm.Links[nodeType]; !ok {
		ret := make(map[string]*ifc.GLink)
		lm.Links[nodeType] = ret
		return ret
	} else {
		return ch
	}
}

func (lm *BaseNodeLinkManager) UpsertRLinkType(nodeType string) map[string]*ifc.GLink {
	lm.Lock()
	defer lm.Unlock()
	if ch, ok := lm.RLinks[nodeType]; !ok {
		ret := make(map[string]*ifc.GLink)
		lm.RLinks[nodeType] = ret
		return ret
	} else {
		return ch
	}
}
func (lm *BaseNodeLinkManager) HasType(nodeType string) bool {
	lm.RLock()
	defer lm.RUnlock()
	dt, ok := lm.Links[nodeType]
	if !ok {
		return false
	}
	return len(dt) > 0
}

func (lm *BaseNodeLinkManager) HasRLinkType(nodeType string) bool {
	lm.RLock()
	defer lm.RUnlock()
	dt, ok := lm.RLinks[nodeType]
	if !ok {
		return false
	}
	return len(dt) > 0
}

func (lm *BaseNodeLinkManager) Has(nodeType, childKey string) bool {
	lm.RLock()
	defer lm.RUnlock()
	ch, ok := lm.Links[nodeType]
	if !ok {
		return false
	}
	_, ok2 := ch[childKey]
	return ok2
}

func (lm *BaseNodeLinkManager) HasRLink(nodeType, childKey string) bool {
	lm.RLock()
	defer lm.RUnlock()
	ch, ok := lm.RLinks[nodeType]
	if !ok {
		return false
	}
	_, ok2 := ch[childKey]
	return ok2
}

func (lm *BaseNodeLinkManager) Get(nodeType, childKey string) (*ifc.GLink, bool) {
	lm.RLock()
	defer lm.RUnlock()
	ch, ok := lm.Links[nodeType]
	if !ok {
		return nil, false
	}
	gl, ok2 := ch[childKey]
	return gl, ok2
}

func (lm *BaseNodeLinkManager) GetRLink(nodeType, childKey string) (*ifc.GLink, bool) {
	lm.RLock()
	defer lm.RUnlock()
	ch, ok := lm.RLinks[nodeType]
	if !ok {
		return nil, false
	}
	gl, ok2 := ch[childKey]
	return gl, ok2
}

func (lm *BaseNodeLinkManager) GetRLinkId(nodeType, childKey string) (string, bool) {
	lm.RLock()
	defer lm.RUnlock()
	ch, ok := lm.RLinks[nodeType]
	if !ok {
		return "", false
	}
	gl, ok2 := ch[childKey]
	if !ok2 {
		return "", ok2
	}
	return gl.DestinationNodeId, ok2
}

func (lm *BaseNodeLinkManager) GetRLinksForNodeType(nodeType string) ([]*ifc.GLink, bool) {
	lm.RLock()
	defer lm.RUnlock()
	_, ok := lm.RLinks[nodeType]
	if !ok {
		return nil, false
	}
	res := []*ifc.GLink{}
	for _, lnk := range lm.RLinks[nodeType] {
		res = append(res, lnk)
	}
	if len(res) == 0 {
		return nil, false
	}
	return res, true
}

func (lm *BaseNodeLinkManager) GetRLinks() ([]*ifc.GLink, bool) {
	lm.RLock()
	defer lm.RUnlock()
	res := []*ifc.GLink{}
	for _, nKeyMap := range lm.RLinks {
		for _, rLnk := range nKeyMap {
			res = append(res, rLnk)
		}
	}
	if len(res) == 0 {
		return nil, false
	}
	return res, true
}

func (lm *BaseNodeLinkManager) isChildLink(lnk *ifc.GLink) bool {
	if lnk != nil && len(lnk.Properties) != 0 {
		if lProp, lOk := lnk.Properties[common.LinkFixedProp_HardLink]; lOk {
			if hLink, hOk := lProp.(bool); hOk {
				return hLink
			}
			if hLinkS, sOk := lProp.(string); sOk {
				return hLinkS == "true"
			}
		}
	}
	return false
}

func (lm *BaseNodeLinkManager) GetChildLink(nodeType, childKey string) (*ifc.GLink, bool) {
	lm.RLock()
	defer lm.RUnlock()
	ch, ok := lm.Links[nodeType]
	if !ok {
		return nil, false
	}
	gl, ok2 := ch[childKey]
	if lm.isChildLink(gl) {
		return gl, ok2
	}
	return &ifc.GLink{}, false
}

func (lm *BaseNodeLinkManager) GetChildLinksForNodeType(nodeType string) ([]*ifc.GLink, bool) {
	lm.RLock()
	defer lm.RUnlock()
	_, ok := lm.Links[nodeType]
	if !ok {
		return nil, false
	}
	res := []*ifc.GLink{}
	for _, lnk := range lm.Links[nodeType] {
		if lm.isChildLink(lnk) {
			res = append(res, lnk)
		}
	}
	if len(res) == 0 {
		return nil, false
	}
	return res, true
}

func (lm *BaseNodeLinkManager) GetChildLinks() ([]*ifc.GLink, bool) {
	lm.RLock()
	defer lm.RUnlock()
	res := []*ifc.GLink{}
	for _, nKeyMap := range lm.Links {
		for _, hLnk := range nKeyMap {
			if lm.isChildLink(hLnk) {
				res = append(res, hLnk)
			}
		}
	}
	if len(res) == 0 {
		return nil, false
	}
	return res, true
}

func (lm *BaseNodeLinkManager) GetRLinkIdsForNodeType(nodeType string) ([]string, bool) {
	lm.RLock()
	defer lm.RUnlock()
	_, ok := lm.RLinks[nodeType]
	if !ok {
		return nil, false
	}
	res := []string{}
	for _, lnk := range lm.RLinks[nodeType] {
		if lnk != nil {
			if lnk.DestinationNodeId != "" {
				res = append(res, lnk.DestinationNodeId)
			}
		}
	}
	if len(res) == 0 {
		return nil, false
	}
	return res, true
}

func (lm *BaseNodeLinkManager) GetRLinkIds() ([]string, bool) {
	lm.RLock()
	defer lm.RUnlock()
	res := []string{}
	for _, nKeyMap := range lm.RLinks {
		for _, rLnk := range nKeyMap {
			if rLnk != nil {
				if rLnk.DestinationNodeId != "" {
					res = append(res, rLnk.DestinationNodeId)
				}
			}
		}
	}
	if len(res) == 0 {
		return nil, false
	}
	return res, true
}

func (lm *BaseNodeLinkManager) SetLinkProperty(nodeType, childKey string, prop ifc.PropertyType) bool {
	lm.RLock()
	defer lm.RUnlock()
	gl, ok := lm.Get(nodeType, childKey)
	if ok {
		gl.Properties = prop
		return true
	}
	return false
}

func (lm *BaseNodeLinkManager) Find(nodeType string, fn func(lnk *ifc.GLink) bool) (string, bool) {
	lm.RLock()
	defer lm.RUnlock()
	ch, ok := lm.Links[nodeType]
	if ok {
		for key, gl := range ch {
			if fn(gl) {
				return key, true
			}
		}
	}
	return "", false
}

func (lm *BaseNodeLinkManager) FindRLink(nodeType string, fn func(lnk *ifc.GLink) bool) (string, bool) {
	lm.RLock()
	defer lm.RUnlock()
	ch, ok := lm.RLinks[nodeType]
	if ok {
		for key, gl := range ch {
			if fn(gl) {
				return key, true
			}
		}
	}
	return "", false
}

func (lm *BaseNodeLinkManager) Add(nodeType, nodeKey string, lnk *ifc.GLink) {
	ch := lm.UpsertType(nodeType)
	lm.Lock()
	defer lm.Unlock()
	ch[nodeKey] = lnk
}

func (lm *BaseNodeLinkManager) Delete(nodeType, nodeKey string) {
	lm.Lock()
	defer lm.Unlock()
	ch, ok := lm.Links[nodeType]
	if !ok {
		return
	}
	delete(ch, nodeKey)
	if len(ch) == 0 {
		delete(lm.Links, nodeType)
	}
}

func (lm *BaseNodeLinkManager) DeleteRLink(nodeType, nodeKey string) {
	lm.Lock()
	defer lm.Unlock()
	ch, ok := lm.RLinks[nodeType]
	if !ok {
		return
	}
	delete(ch, nodeKey)
	if len(ch) == 0 {
		delete(lm.RLinks, nodeType)
	}
}

func (lm *BaseNodeLinkManager) ForEach(fn func(ntype, nkey string, lnk *ifc.GLink)) {
	// make a copy and then iterate?
	lm.RLock()
	typeList := make([]string, len(lm.Links))
	idx := 0
	for t, _ := range lm.Links {
		typeList[idx] = t
		idx++
	}
	lm.RUnlock()
	for _, typeName := range typeList {
		lm.ForEachType(typeName, fn)
	}
}

func (lm *BaseNodeLinkManager) ForEachRLink(fn func(ntype, nkey string, lnk *ifc.GLink)) {
	// make a copy and then iterate?
	lm.RLock()
	typeList := make([]string, len(lm.RLinks))
	idx := 0
	for t, _ := range lm.RLinks {
		typeList[idx] = t
		idx++
	}
	lm.RUnlock()
	for _, typeName := range typeList {
		lm.ForEachRLinkType(typeName, fn)
	}
}

func (lm *BaseNodeLinkManager) ForEachType(ntype string, fn func(ntype, nkey string, lnk *ifc.GLink)) {
	lm.RLock()
	child, ok := lm.Links[ntype]
	if !ok {
		lm.RUnlock()
		return
	}
	keyList := make([]string, len(child))
	idx := 0
	for k, _ := range child {
		keyList[idx] = k
		idx++
	}
	lm.RUnlock()

	for _, nkey := range keyList {
		lm.RLock()
		lnk, ok := child[nkey]
		lm.RUnlock()
		if ok {
			fn(ntype, nkey, lnk)
		}
	}
}

func (lm *BaseNodeLinkManager) ForEachRLinkType(ntype string, fn func(ntype, nkey string, lnk *ifc.GLink)) {
	lm.RLock()
	child, ok := lm.RLinks[ntype]
	if !ok {
		lm.RUnlock()
		return
	}
	keyList := make([]string, len(child))
	idx := 0
	for k, _ := range child {
		keyList[idx] = k
		idx++
	}
	lm.RUnlock()

	for _, nkey := range keyList {
		lm.RLock()
		lnk, ok := child[nkey]
		lm.RUnlock()
		if ok {
			fn(ntype, nkey, lnk)
		}
	}
}
