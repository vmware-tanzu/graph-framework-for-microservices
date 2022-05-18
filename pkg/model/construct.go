package model

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sync"
)

var (
	GlobalRestURIChan        = make(chan []nexus.RestURIs, 100)
	GlobalEndpointCache      = make(map[string][]nexus.RestURIs)
	globalEndpointCacheMutex = &sync.Mutex{}

	GlobalURIToCRDTypes     = make(map[string]string)
	globalURIToCRDTypeMutex = &sync.Mutex{}

	GlobalCRDTypeToNodes     = make(map[string]NodeInfo)
	globalCRDTypeToNodeMutex = &sync.Mutex{}
	GlobalCRDChan            = make(chan string, 100)

	GlobalCRDTypeToSpec      = make(map[string]apiextensionsv1.CustomResourceDefinitionSpec)
	globalCRDTypeToSpecMutex = &sync.Mutex{}
)

func ConstructMapURIToCRDType(eventType EventType, crdType string, apiURIs []nexus.RestURIs) {
	globalURIToCRDTypeMutex.Lock()
	defer globalURIToCRDTypeMutex.Unlock()

	if eventType == Delete {
		for uri, cType := range GlobalURIToCRDTypes {
			if cType == crdType {
				delete(GlobalURIToCRDTypes, uri)
			}
		}
	}

	for _, u := range apiURIs {
		GlobalURIToCRDTypes[u.Uri] = crdType
	}
}

func ConstructMapCRDTypeToNode(eventType EventType, crdType, name string, parentHierarchy []string, children map[string]NodeHelperChild) {
	globalCRDTypeToNodeMutex.Lock()
	defer globalCRDTypeToNodeMutex.Unlock()

	if eventType == Delete {
		delete(GlobalCRDTypeToNodes, crdType)
	}

	GlobalCRDTypeToNodes[crdType] = NodeInfo{
		Name:            name,
		ParentHierarchy: parentHierarchy,
		Children:        children,
	}

	GlobalCRDChan <- crdType
}

func ConstructMapCRDTypeToSpec(eventType EventType, crdType string, spec apiextensionsv1.CustomResourceDefinitionSpec) {
	globalCRDTypeToSpecMutex.Lock()
	defer globalCRDTypeToSpecMutex.Unlock()

	if eventType == Delete {
		delete(GlobalCRDTypeToNodes, crdType)
	}

	GlobalCRDTypeToSpec[crdType] = spec
}

// revisit for simplification
func addOrRemoveEndpointCache(endpointCache, removed, added []nexus.RestURIs) []nexus.RestURIs {
	for _, r := range removed {
		idx := -1
		for i, c := range endpointCache {
			if r.Uri == c.Uri {
				idx = i
				for rm := range r.Methods {
					pdx := -1
					for nm := range c.Methods {
						pdx = i
						if rm == nm {
							break
						}
					}
					if pdx > -1 {
						delete(c.Methods, rm)
					}
				}
			}
		}
		if idx > -1 && len(endpointCache[idx].Methods) == 0 {
			endpointCache = append(endpointCache[:idx], endpointCache[idx+1:]...)
		}
	}

	for _, a := range added {
		endpointCache = append(endpointCache, a)
	}

	return endpointCache
}

func findDifference(existingEndpoints, newURIs []nexus.RestURIs) []nexus.RestURIs {
	removedOrAddedURLs := make([]nexus.RestURIs, 0)
	removedOrAddedMethods := make([]nexus.RestURIs, 0)
	for _, existingEp := range existingEndpoints {
		urlNotEncountered := true
		for _, newEndpoint := range newURIs {
			if existingEp.Uri == newEndpoint.Uri {
				urlNotEncountered = false
				for m, r := range existingEp.Methods {
					methodNotEncountered := true
					for n := range newEndpoint.Methods {
						if m == n {
							methodNotEncountered = false
							break
						}
					}
					if methodNotEncountered {
						removedOrAddedMethods = append(removedOrAddedMethods, nexus.RestURIs{
							Uri:     existingEp.Uri,
							Methods: nexus.HTTPMethodsResponses{m: r},
						})
					}
				}
			}
		}
		if urlNotEncountered {
			removedOrAddedURLs = append(removedOrAddedURLs, existingEp)
		}
	}

	if len(removedOrAddedMethods) > 0 {
		removedOrAddedURLs = append(removedOrAddedURLs, removedOrAddedMethods...)
	}

	return removedOrAddedURLs
}

func PopulateEndpointCache(evenType EventType, crdType string, restURIs []nexus.RestURIs) (int, []nexus.RestURIs) {
	globalEndpointCacheMutex.Lock()
	defer globalEndpointCacheMutex.Unlock()

	removed := []nexus.RestURIs{}
	cachedEndpoints, ok := GlobalEndpointCache[crdType]
	if ok {
		if evenType == Delete {
			delete(GlobalEndpointCache, crdType)
			return len(cachedEndpoints), nil
		}
		removed = findDifference(cachedEndpoints, restURIs)
	}
	added := findDifference(restURIs, cachedEndpoints)

	// update the cache with removed and added entries
	finalEPCache := addOrRemoveEndpointCache(cachedEndpoints, removed, added)
	GlobalEndpointCache[crdType] = finalEPCache

	return len(removed), added
}

func GetGlobalEndpointCache() (epCacheCopy []nexus.RestURIs) {
	globalEndpointCacheMutex.Lock()
	defer globalEndpointCacheMutex.Unlock()

	for _, ep := range GlobalEndpointCache {
		epCacheCopy = append(epCacheCopy, ep...)
	}
	return epCacheCopy
}
