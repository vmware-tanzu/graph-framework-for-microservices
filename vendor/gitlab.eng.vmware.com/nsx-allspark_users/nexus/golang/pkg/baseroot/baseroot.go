package baseroot

import (
	"sync"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/internal/datamodel"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/ifc"
)

var (
	globalDM       ifc.DataModelInterface = nil
	globalRootBase ifc.BaseNodeInterface  = nil
	globalRootLock sync.Mutex             // used to create root singleton
)

func NewBaseRootNode(dmName, rootName, etcdLocation string, featureFlags string, singleton bool) ifc.BaseNodeInterface {
	globalRootLock.Lock()
	defer globalRootLock.Unlock()
	rootId := "/Root/" + rootName

	if singleton {
		if globalRootBase != nil {
			return globalRootBase
		}
		globalDM := datamodel.NewDataModel(dmName, etcdLocation, featureFlags)
		globalDM.Init()
		globalRootBase, nodeOk := globalDM.GetRootNode(rootId, false)
		if !nodeOk {
			globalRootBase = globalDM.UpsertRootNode("Root", rootName, nil)
		}
		return globalRootBase
	} else {
		dm := datamodel.NewDataModel(dmName, etcdLocation, featureFlags)
		dm.Init()
		node, nodeOk := dm.GetRootNode(rootId, false)
		if !nodeOk {
			node = dm.UpsertRootNode("Root", rootName, nil)
		}
		return node
	}
}
