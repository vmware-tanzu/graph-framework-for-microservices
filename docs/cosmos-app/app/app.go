package app

import (
	"fmt"

	v2 "golang-appnet.eng.vmware.com/cosmos-datamodel/apis/global.cosmos.tanzu.vmware.org/v1"
	nexus_client "golang-appnet.eng.vmware.com/cosmos-datamodel/nexus-client"
)

func App(nexusClient *nexus_client.Clientset) {
	fmt.Println("My APP")
	// Subscribe Org
	nexusClient.GlobalOrg("*").Subscribe()
	// Org: -> RegisterAddCallback
	nexusClient.GlobalOrg("*").RegisterAddCallback(onAddOrg)
	// Org: -> RegisterUpdateCallback
	nexusClient.GlobalOrg("*").RegisterUpdateCallback(onUpdateOrg)
	// Org: -> RegisterDeleteCallback
	nexusClient.GlobalOrg("*").RegisterDeleteCallback(onDeleteOrg)
	// Project: -> RegisterAddCallback
	nexusClient.GlobalOrg("*").Project("*").RegisterAddCallback(onAddProject)
	// GlobalNamespace: -> RegisterAddCallback
	nexusClient.GlobalOrg("*").Project("*").Config().GlobalNamespace("*").RegisterAddCallback(onAddGNS)

}

// onAdd is the function executed when the kubernetes informer notified the
// presence of a new kubernetes node in the cluster
func onAddOrg(obj *v2.Org) {
	// Cast the obj as node
	fmt.Printf("ORG =====>>>>>>>>>>>ON ADD EVENT!: %s\n", obj.Name)
}

func onUpdateOrg(oldObj, newObj *v2.Org) {
	// Cast the obj as node
	fmt.Printf("ORG: ===> UPDATE EVENT! OLD:>>> %s \n SPEC: NAME >>> %s\n", oldObj.Name, oldObj.Spec.Name)
	fmt.Printf("ORG: ===> UPDATE EVENT! New:>>> %s \n SPEC: NAME >>> %s\n", newObj.Name, newObj.Spec.Name)
}

func onDeleteOrg(obj *v2.Org) {
	// Cast the obj as node
	fmt.Printf("ORG: ===> DELETE EVENT! :%s\n", obj.Name)
}

func onAddProject(obj *v2.Project) {
	// Cast the obj as node
	fmt.Printf("PROJECT =====>>>>>>>>>>>ADD EVENT!: %s\n\n", obj.Name)
}

func onAddGNS(obj *v2.GlobalNamespace) {
	// Cast the obj as node
	fmt.Printf("GNS =====>>>>>>>>>>>ADD EVENT!: %s\n\n", obj.Name)
}
