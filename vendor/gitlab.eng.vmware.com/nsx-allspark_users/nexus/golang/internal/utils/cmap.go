package utils

import (
	"C"
	"sync"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/ifc"
)
import (
	"errors"
	"fmt"
	"strconv"
)

var SHARDS = uint64(64)
var CMCnt uint64 = 1

type ConcurrentMap []*ConcurrentMapShared
type ConcurrentMapShared struct {
	items map[uint64]interface{}
	sync.RWMutex
}

func NewConcurrentMap() ConcurrentMap {
	m := make(ConcurrentMap, SHARDS)
	for i := uint64(0); i < SHARDS; i++ {
		m[i] = &ConcurrentMapShared{items: make(map[uint64]interface{})}
	}
	return m
}
func (m ConcurrentMap) GetShard(key uint64) *ConcurrentMapShared {
	return m[key%SHARDS]
}
func (m ConcurrentMap) Store(value interface{}) uint64 {
	key := CMCnt
	CMCnt++
	shard := m.GetShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
	return key
}
func (m ConcurrentMap) Set(key uint64, value interface{}) {
	shard := m.GetShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}
func (m ConcurrentMap) Load(key uint64) (interface{}, bool) {
	shard := m.GetShard(key)
	shard.RLock()
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}
func (m ConcurrentMap) Delete(key uint64) {
	shard := m.GetShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

func (m ConcurrentMap) GetDM(key uint64) (ifc.DataModelInterface, bool) {
	d, dok := m.Load(key)
	if dok {
		return d.(ifc.DataModelInterface), true
	}
	return nil, false
}

func (m ConcurrentMap) GetLinkManager(key uint64) (ifc.BaseNodeLinkManagerInterface, bool) {
	d, dok := m.Load(key)
	if dok {
		return d.(ifc.BaseNodeLinkManagerInterface), true
	}
	return nil, false
}
func (m ConcurrentMap) GetLinkManagerP(key uint64) ifc.BaseNodeLinkManagerInterface {
	d, dok := m.GetLinkManager(key)
	if dok {
		return d
	}
	panic(errors.New(fmt.Sprintf("Not found LinkManager for for handle %d", key)))
}

func (m ConcurrentMap) GetDMP(key uint64) ifc.DataModelInterface {
	d, dok := m.GetDM(key)
	if dok {
		return d
	}
	panic(errors.New(fmt.Sprintf("Not found dm for handle %d", key)))
}

var CM ConcurrentMap = NewConcurrentMap()
var B ConcurrentMap = NewConcurrentMap()

func SetBase(bn ifc.BaseNodeInterface) uint64 {
	if bn == nil {
		return 0
	}
	var k uint64 = bn.GetCCnt()
	B.Set(k, bn)
	return k
}
func GetBase(key uint64) (ifc.BaseNodeInterface, bool) {
	d, dok := B.Load(key)
	if dok {
		return d.(ifc.BaseNodeInterface), true
	}
	return nil, false
}
func GetBaseP(keyStr string) ifc.BaseNodeInterface {
	key, err := strconv.ParseUint(keyStr, 10, 64)
	if err != nil {
		panic(errors.New(fmt.Sprintf("GetBaseP:Error when converting key %s to uint64", keyStr)))
	}
	d, dok := GetBase(key)
	if dok {
		return d
	}
	panic(errors.New(fmt.Sprintf("Not found Base for handle %d", key)))
}

var DM ConcurrentMap = NewConcurrentMap()

func SetDM(bn ifc.DataModelInterface) uint64 {
	if bn == nil {
		return 0
	}
	var k uint64 = bn.GetCCnt()
	DM.Set(k, bn)
	return k
}
func GetDM(key uint64) (ifc.DataModelInterface, bool) {
	d, dok := DM.Load(key)
	if dok {
		return d.(ifc.DataModelInterface), true
	}
	return nil, false
}
func GetDMP(keyStr string) ifc.DataModelInterface {
	key, err := strconv.ParseUint(keyStr, 10, 64)
	if err != nil {
		panic(errors.New(fmt.Sprintf("GetDMP:Error when converting key %s to uint64", keyStr)))
	}
	d, dok := GetDM(key)
	if dok {
		return d
	}
	panic(errors.New(fmt.Sprintf("GetDMP:Not found Datamodel for handle %d", key)))
}
