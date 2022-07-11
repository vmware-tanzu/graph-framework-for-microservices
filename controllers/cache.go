package controllers

import (
	"sync"

	log "github.com/sirupsen/logrus"
	handler "gitlab.eng.vmware.com/nsx-allspark_users/m7/handler.git"
)

type CrdCache struct {
	mu     sync.RWMutex
	CrdMap map[string]*CrdInfo
}

type CrdInfo struct {
	Spec       interface{}
	Controller *handler.Controller
}

func NewCrdCache() *CrdCache {
	return &CrdCache{
		mu:     sync.RWMutex{},
		CrdMap: make(map[string]*CrdInfo),
	}
}

func (c *CrdCache) Upsert(crdName string, crdInfo *CrdInfo) {
	log.Debugf("[Add/Update] new CRD %s to cache", crdName)
	c.mu.Lock()
	c.CrdMap[crdName] = crdInfo
	c.mu.Unlock()
}

func (c *CrdCache) Get(crdName string) *CrdInfo {
	log.Debugf("[Get] CRD %s from cache", crdName)
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.CrdMap[crdName]
}

func (c *CrdCache) Delete(crdName string) {
	log.Debugf("[Delete] CRD %s from cache", crdName)
	c.mu.Lock()
	delete(c.CrdMap, crdName)
	c.mu.Unlock()
}

func (c *CrdCache) UpsertController(crdName string, controller *handler.Controller) {
	log.Debugf("[Add] %s controller to CRD Cache", crdName)
	c.mu.Lock()
	crdInfo := c.CrdMap[crdName]
	crdInfo.Controller = controller
	c.CrdMap[crdName] = crdInfo
	c.mu.Unlock()
}
