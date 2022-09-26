package controllers

import (
	"sync"

	log "github.com/sirupsen/logrus"
	handler "gitlab.eng.vmware.com/nsx-allspark_users/m7/handler.git"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

/* GvrCache stores controller info for a resource type.
The resource type could be CRD Types, Deployments, Services, Pods or any K8s resources.

Eg: GvrCache map entries will look like:
{Group: vmware.org, Version: v1, Resource: roots} => controller{}
{Group: vmware.org, Version: v1beta1, Resource: projects} => controller{}
{Group: apps, Version: v1, Resource: deployments} => controller{}
*/
type GvrCache struct {
	mu     sync.RWMutex
	GvrMap map[schema.GroupVersionResource]*GvrInfo
}

type GvrInfo struct {
	Controller *handler.Controller
}

func NewGvrCache() *GvrCache {
	return &GvrCache{
		mu:     sync.RWMutex{},
		GvrMap: make(map[schema.GroupVersionResource]*GvrInfo),
	}
}

func (c *GvrCache) Get(gvr schema.GroupVersionResource) *GvrInfo {
	log.Debugf("[Get] GVR %s from cache", gvr)
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.GvrMap[gvr]
}

func (c *GvrCache) Delete(gvr schema.GroupVersionResource) {
	log.Debugf("[Delete] GVR %s from cache", gvr)
	c.mu.Lock()
	delete(c.GvrMap, gvr)
	c.mu.Unlock()
}

func (c *GvrCache) UpsertController(gvr schema.GroupVersionResource, controller *handler.Controller) {
	log.Debugf("[Add] %s controller to GVR Cache", gvr)
	c.mu.Lock()
	c.GvrMap[gvr] = &GvrInfo{Controller: controller}
	c.mu.Unlock()
}
