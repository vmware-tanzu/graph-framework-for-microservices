package dmcache

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/ifc"
	log "gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/logging"
)

const defaultDebugServerPort = 6000

// Used to dump DM cache data
var (
	GlobalDMCacheMap   = make(map[string]ifc.DataModelInterface)
	GlobalDMCacheMutex sync.Mutex
	DebugServerStarted bool // used to start debug server
)

func StartDebugServer(port string) {

	http.HandleFunc("/debug/nexus/dm", CollectDM)

	debugServerPort := defaultDebugServerPort
	// Expect app to set the Debug Server Port. If not, set the default one.
	if len(port) != 0 {
		p, err := strconv.Atoi(port)
		if err == nil {
			debugServerPort = p
		} else {
			log.Errorf("Error parsing debug server port %s", err)
		}
	}

	addr := fmt.Sprintf(":%d", debugServerPort)
	log.Infof("Starting DM Debug Server at %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("DM Debug Server Error: %v", err)
	}
}

func CollectDM(w http.ResponseWriter, r *http.Request) {
	// check for query params
	if len(r.URL.Query()) == 0 {
		listDM(w, r)
		return
	}

	//extract dmID segments from request
	dmID := r.URL.Query().Get("dmID")
	if len(dmID) == 0 {
		renderJSON(w, http.StatusBadRequest, "missing DM ID")
		return
	}

	respondWithDMCachedData(w, dmID)
}

func respondWithDMCachedData(w http.ResponseWriter, dmID string) {
	dmCacheMap := GetDMGlobalCacheMap()
	dm, ok := dmCacheMap[dmID]
	if !ok {
		err := fmt.Sprintf("dm (%s) does not exist", dmID)
		log.Debugf(err)
		renderJSON(w, http.StatusBadRequest, err)
		return
	}

	result := dm.CollectDMCachedData()
	renderJSON(w, http.StatusOK, result)
}

type CacheMetatdata struct {
	DMID   string
	DMName string
}

func listDM(w http.ResponseWriter, _ *http.Request) {
	log.Debugf("List DM IDs and names")

	result := []CacheMetatdata{}
	dmCacheMap := GetDMGlobalCacheMap()
	for _, dm := range dmCacheMap {
		result = append(result, CacheMetatdata{
			DMID:   dm.GetId(),
			DMName: dm.GetName(),
		})
	}

	renderJSON(w, http.StatusOK, result)
}

func GetDMGlobalCacheMap() map[string]ifc.DataModelInterface {
	GlobalDMCacheMutex.Lock()
	defer GlobalDMCacheMutex.Unlock()

	newDMCacheMap := map[string]ifc.DataModelInterface{}
	for k, v := range GlobalDMCacheMap {
		newDMCacheMap[k] = v
	}

	return newDMCacheMap
}

func ClearGlobalDMCacheMap() {
	GlobalDMCacheMutex.Lock()
	defer GlobalDMCacheMutex.Unlock()

	for id := range GlobalDMCacheMap {
		delete(GlobalDMCacheMap, id)
	}
}
