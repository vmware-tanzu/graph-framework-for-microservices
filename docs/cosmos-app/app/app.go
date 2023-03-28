package app

import (
	"fmt"

	orgv2 "golang-appnet.eng.vmware.com/cosmos-datamodel/apis/global.cosmos.tanzu.vmware.org/v1"
	nexus_client "golang-appnet.eng.vmware.com/cosmos-datamodel/nexus-client"
)

func App(nexusClient *nexus_client.Clientset) {
	fmt.Println("My APP")
	nexusClient.GlobalOrg("*").Subscribe()
	_, err := nexusClient.GlobalOrg("*").RegisterAddCallback(onAdd)
	if err != nil {
		fmt.Printf("ERROR: %v+", err)
	}
	nexusClient.GlobalOrg("*").Project("*").Config().GlobalNamespace("*").Gns("*")
	nexusClient.GlobalOrg("*").RegisterUpdateCallback(onUpdate)
	nexusClient.GlobalOrg("*").RegisterDeleteCallback(onDelete)
	regId, err := nexusClient.GlobalOrg("*").Project("*").RegisterAddCallback(onAddProject)
	if err != nil {
		fmt.Printf("ERROR: %v+", err)
	}
	nexusClient.GlobalOrg("*").Project("*").UnRegisterAddCallback(regId)
	// fmt.Printf("rID: %v+", rID)

}

// onAdd is the function executed when the kubernetes informer notified the
// presence of a new kubernetes node in the cluster
func onAdd(obj *orgv2.Org) {
	// Cast the obj as node
	fmt.Printf("NEW =====>>>>>>>>>>>ADD EVENT!: %s\n\n", obj.Name)
}

func onUpdate(oldObj, newObj *orgv2.Org) {
	// Cast the obj as node
	fmt.Printf("UPDATE EVENT! OLD:>>> %s \n SPEC: NAME >>> %s\n", oldObj.Name, oldObj.Spec.Name)
	fmt.Printf("UPDATE EVENT! New:>>> %s \n SPEC: NAME >>> %s\n", newObj.Name, newObj.Spec.Name)
}

func onDelete(obj *orgv2.Org) {
	// Cast the obj as node
	fmt.Printf("DELETE EVENT! :%s\n", obj.Name)
}

func onAddProject(obj *orgv2.Project) {
	// Cast the obj as node
	fmt.Printf("PROJECT =====>>>>>>>>>>>ADD EVENT!: %s\n\n", obj.Name)
}

/*<-
func (c *orgGlobalCosmosV1Chainer) RegisterAddCallback(cbfn func(obj *baseglobalcosmostanzuvmwareorgv1.Org)) (cache.ResourceEventHandlerRegistration, error) {
    fmt.Println("GlobalOrg -->  RegisterAddCallback!")
    var (
        registrationId cache.ResourceEventHandlerRegistration
        err            error
    )
    key := "orgs.global.cosmos.tanzu.vmware.org"
    stopper := make(chan struct{})
    if s, ok := subscriptionMap.Load(key); ok {
        fmt.Println("[GlobalOrg] ---SUBSCRIBE-INFORMER---->")
        sub := s.(subscription)
        registrationId, err = sub.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
            AddFunc: func(obj interface{}) {
                cbfn(obj.(*baseglobalcosmostanzuvmwareorgv1.Org))
            },
        })
    } else {
        fmt.Println("[GlobalOrg] ---NEW-INFORMER---->")
        informer := informerglobalcosmostanzuvmwareorgv1.NewOrgInformer(c.client.baseClient, 0, cache.Indexers{})
        registrationId, err = informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
            AddFunc: func(obj interface{}) {
                cbfn(obj.(*baseglobalcosmostanzuvmwareorgv1.Org))
            },
        })
        go informer.Run(stopper)
    }
    return registrationId, err
}
func (c *orgGlobalCosmosV1Chainer) RegisterUpdateCallback(cbfn func(oldObj, newObj *baseglobalcosmostanzuvmwareorgv1.Org)) (cache.ResourceEventHandlerRegistration, error) {
    fmt.Println("GlobalOrg -->  RegisterUpdateCallback!")
    var (
        registrationId cache.ResourceEventHandlerRegistration
        err            error
    )
    key := "orgs.global.cosmos.tanzu.vmware.org"
    stopper := make(chan struct{})
    if s, ok := subscriptionMap.Load(key); ok {
        fmt.Println("[GlobalOrg] ---SUBSCRIBE-INFORMER---->")
        sub := s.(subscription)
        registrationId, err = sub.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
            UpdateFunc: func(oldObj, newObj interface{}) {
                cbfn(oldObj.(*baseglobalcosmostanzuvmwareorgv1.Org), newObj.(*baseglobalcosmostanzuvmwareorgv1.Org))
            },
        })
    } else {
        fmt.Println("[GlobalOrg] ---NEW-INFORMER---->")
        informer := informerglobalcosmostanzuvmwareorgv1.NewOrgInformer(c.client.baseClient, 0, cache.Indexers{})
        registrationId, err = informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
            UpdateFunc: func(oldObj, newObj interface{}) {
                cbfn(oldObj.(*baseglobalcosmostanzuvmwareorgv1.Org), newObj.(*baseglobalcosmostanzuvmwareorgv1.Org))
            },
        })
        go informer.Run(stopper)
    }
    return registrationId, err
}
func (c *orgGlobalCosmosV1Chainer) RegisterDeleteCallback(cbfn func(obj *baseglobalcosmostanzuvmwareorgv1.Org)) (cache.ResourceEventHandlerRegistration, error) {
    fmt.Println("GlobalOrg -->  RegisterDeleteCallback!")
    var (
        registrationId cache.ResourceEventHandlerRegistration
        err            error
    )
    key := "orgs.global.cosmos.tanzu.vmware.org"
    stopper := make(chan struct{})
    if s, ok := subscriptionMap.Load(key); ok {
        fmt.Println("[GlobalOrg] ---SUBSCRIBE-INFORMER---->")
        sub := s.(subscription)
        registrationId, err = sub.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
            DeleteFunc: func(obj interface{}) {
                cbfn(obj.(*baseglobalcosmostanzuvmwareorgv1.Org))
            },
        })
    } else {
        fmt.Println("[GlobalOrg] ---NEW-INFORMER---->")
        informer := informerglobalcosmostanzuvmwareorgv1.NewOrgInformer(c.client.baseClient, 0, cache.Indexers{})
        registrationId, err = informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
            DeleteFunc: func(obj interface{}) {
                cbfn(obj.(*baseglobalcosmostanzuvmwareorgv1.Org))
            },
        })
        go informer.Run(stopper)
    }
    return registrationId, err
}
*/
