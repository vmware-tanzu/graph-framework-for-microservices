package controllers

import (
	"fmt"
	"sync"

	"k8s.io/apimachinery/pkg/util/json"

	"api-gw/pkg/model"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
)

var (
	GlobalRestURIChan        = make(chan []nexus.RestURIs, 100)
	GlobalEndpointCache      = make(map[string][]nexus.RestURIs)
	globalEndpointCacheMutex = &sync.Mutex{}

	GlobalURIToCRDTypes     = make(map[string]string)
	globalURIToCRDTypeMutex = &sync.Mutex{}

	GlobalCRDTypeToNodes     = make(map[string]model.NodeInfo)
	globalCRDTypeToNodeMutex = &sync.Mutex{}
)

func constructMapURIToCRDType(eventType model.EventType, crdType string, apiURIs []nexus.RestURIs) {
	globalURIToCRDTypeMutex.Lock()
	defer globalURIToCRDTypeMutex.Unlock()

	if eventType == model.Delete {
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

func constructMapCRDTypeToNode(eventType model.EventType, crdType, name string, parentHierarchy []string) {
	globalCRDTypeToNodeMutex.Lock()
	defer globalCRDTypeToNodeMutex.Unlock()

	if eventType == model.Delete {
		delete(GlobalCRDTypeToNodes, crdType)
	}

	GlobalCRDTypeToNodes[crdType] = model.NodeInfo{
		Name:            name,
		ParentHierarchy: parentHierarchy,
	}
}

func (r *CustomResourceDefinitionReconciler) ProcessAnnotation(crdType string,
	annotations map[string]string, eventType model.EventType) error {
	n := model.NexusAnnotation{}

	if eventType != model.Delete {
		apiInfo, ok := annotations["nexus"]
		if !ok {
			return nil
		}

		// unmarshall to nexus annotation struct
		err := json.Unmarshal([]byte(apiInfo), &n)
		if err != nil {
			fmt.Printf("Error unmarshaling Nexus annotation %v\n", err)
			return err
		}
	}

	// It has stored the URI with the CRD type and CRD type with the Node Info.
	constructMapURIToCRDType(eventType, crdType, n.NexusRestAPIGen.Uris)
	constructMapCRDTypeToNode(eventType, crdType, n.Name, n.Hierarchy)

	/*
	 populateEndpointCache populates the cache with CRD Type to restURIs attribute ( URL and method [GET, DELETE...]).
	 if any of the attribute removed in the new event notification, that should be removed from the cache and
	 triggers the server restart to remove the routes.
	 If any of the attribute added newly, notify that to `GlobalRestURIChan`.
	*/
	removed, added := populateEndpointCache(eventType, crdType, n.NexusRestAPIGen.Uris)
	if removed > 0 {
		r.StopCh <- struct{}{}
		GlobalRestURIChan <- getGlobalEndpointCache()
		return nil
	}

	if len(added) > 0 {
		GlobalRestURIChan <- added
	}
	return nil
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

func populateEndpointCache(evenType model.EventType, crdType string, restURIs []nexus.RestURIs) (int, []nexus.RestURIs) {
	globalEndpointCacheMutex.Lock()
	defer globalEndpointCacheMutex.Unlock()

	removed := []nexus.RestURIs{}
	cachedEndpoints, ok := GlobalEndpointCache[crdType]
	if ok {
		if evenType == model.Delete {
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

func getGlobalEndpointCache() (epCacheCopy []nexus.RestURIs) {
	globalEndpointCacheMutex.Lock()
	defer globalEndpointCacheMutex.Unlock()

	for _, ep := range GlobalEndpointCache {
		epCacheCopy = append(epCacheCopy, ep...)
	}
	return epCacheCopy
}
