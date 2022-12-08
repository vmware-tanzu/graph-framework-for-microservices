# Subscribe API

Subscribe API feature is to provide a way to create cache for the nexus nodes. If GET and LIST calls will be made fequently for the nexus node then creating cache for the nexus node will be the best thing to avoid making a lot of calls to the api-server and hence to reduce the load to the api-server.

Subscribe API feature in nexus-client library is implemeted using informers from k8s generated code.
Link to see the example of k8s informers generated code https://github.com/vmware-tanzu/graph-framework-for-microservices/tree/main/compiler/example/output/generated/client/informers/externalversions

### Below methods are present in nexus-client library to use Subscribe API feature:

1. Subscribe() method to create cache for the node.

2. Unsubscribe() method to remove cache for the node(if node has been subscribed for cache).

3. SubscribeAll() to create cache for all the nexus nodes.

4. UnsubscribeAll() to remove cache for all the nexus nodes.

5. IsSubscribed() to check if node is subscribed or not.


### Demo code For Subscribe API Feature:

```
func main.go(){

    // get the nexus client 
    config := &rest.Config{
        Host: "nexus-apiserver:8080",
    }
    nexusClient, _ := nexus_client.NewForConfig(config)

    UseSubscribeAPIFeatue(nexusClient)
}

// UseSubscribeAPIFeature method is show the subsribe API Feature demo
func UseSubscribeAPIFeature(nexusClient *nexus_client.Clientset) {

    // create root
    rootDef := &orgchartv1.Root{
        ObjectMeta: metav1.ObjectMeta{
            Name: "default",
        },
    }
    nexusClient.AddOrgchartRoot(context.TODO(), rootDef)

    // create leader
    leaderDef := &managementvmwareorgv1.Leader{
        ObjectMeta: metav1.ObjectMeta{
            Name: "default",
        },
        Spec: managementvmwareorgv1.LeaderSpec{
            Designation: "Chief",
            Name:        "ABC",
            EmployeeID:  1,
        },
    }
    root, _ := nexusClient.GetOrgchartRoot(context.TODO())
    root.AddCEO(context.TODO(), leaderDef)

    // subscribe all nodes
    nexusClient.SubscribeAll()
    fmt.Println("---------subscribed all nodes-------------")
    time.Sleep(5 * time.Second) // because cache is eventually consistent, may need some time to update cache
    GetList(nexusClient)

    // unsubscribe for particular node(Root node)
    nexusClient.OrgchartRoot().Unsubscribe()
    fmt.Println("---------unsubscribed root-------------")
    GetList(nexusClient)

    // check if a node is subscirbed
    if !nexusClient.OrgchartRoot().IsSubscribed(){
        fmt.Println("Root is not subscribed.")
    }else{
        fmt.Println("Root is subscribed.")
    }

    // unsubscribe all nodes
    nexusClient.UnsubscribeAll()
    fmt.Println("---------unsubscribed all nodes-------------")
    GetList(nexusClient)

    // subscribe for particular node(Leader node)
    nexusClient.OrgchartRoot().CEO().Subscribe()
    fmt.Println("---------subscribed Leader-------------")
    time.Sleep(5 * time.Second) // because cache is eventually consistent, may need some time to update cache
    GetList(nexusClient)
}

func GetList(nexusClient *nexus_client.Clientset) {
    fmt.Println("############ GET and LIST calls ###########")
    root, _ := nexusClient.GetOrgchartRoot(context.TODO())
    fmt.Printf("root: %+v\n", root.Root.Spec)
    roots, _ := nexusClient.Orgchart().ListRoots(context.TODO(), metav1.ListOptions{})
    for _, r := range roots {
        fmt.Printf("Roots: %+v\n", r.Root.Spec)
    }
    ceo, _ := root.GetCEO(context.TODO())
    fmt.Printf("leader: %+v\n", ceo.Leader.Spec)
    leaders, _ := nexusClient.Management().ListLeaders(context.TODO(), metav1.ListOptions{})
    for _, l := range leaders {
        fmt.Printf("Leaders: %+v\n", l.Leader.Spec)
    }
}
```
