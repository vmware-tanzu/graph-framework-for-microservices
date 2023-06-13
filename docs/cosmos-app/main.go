package main

import (
	"context"
	"cosmos-app/app"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	v2 "golang-appnet.eng.vmware.com/cosmos-datamodel/apis/global.cosmos.tanzu.vmware.org/v1"
	nexus_client "golang-appnet.eng.vmware.com/cosmos-datamodel/nexus-client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	config := getK8sAPIEndpointConfig()
	// Create a datamodel client handle.
	nexusClient, err := nexus_client.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	app.App(nexusClient)
	crudWithNexusClient(nexusClient)
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()

}

func crudWithNexusClient(nexusClient *nexus_client.Clientset) {
	fmt.Printf("During run of this application we will create following graph:\n\n" +
		"            Org\n" +
		"             | \n" +
		"             |  \n" +
		"          Project \n" +
		"       /     |     \\\n" +
		"      /      |      \\\n" +
		"  Config Inventory  Runtime\n" +
		"     | \n" +
		"     | \n" +
		"  Globalnamespace \n" +
		"     | \n" +
		"     | \n" +
		"   GNS \n" +
		"\n\n")

	// Org Node Definition
	orgDef := &v2.Org{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
	}
	fmt.Println("Creating Org object...")
	// CREATE: --> Org
	org, err := nexusClient.AddGlobalOrg(context.TODO(), orgDef)
	if err != nil {
		panic(err)
	}
	fmt.Println("... checking if Org is created properly, name should be hashed, original name is preserved in "+
		"nexus/display_name label", org.DisplayName())
	// GET: -->  Org
	getOrg, err := nexusClient.GetGlobalOrg(context.TODO(), "default")
	if err != nil {
		panic(err)
	}
	printdata("Org hash name:", getOrg.Name)
	// Project Node Definition
	fmt.Println("Creating Project object...")
	projDef := &v2.Project{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
	}
	// CREATE: --> Project
	proj, err := getOrg.AddProject(context.TODO(), projDef)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Project: %s ----> Created \n", proj.DisplayName())
	// GET: --> Project
	getProj, err := getOrg.GetProject(context.TODO(), "default")
	if err != nil {
		panic(err)
	}
	printdata("Get Project hash name:", getProj.Name)
	// Config Node Definition
	fmt.Println("Creating Config object...")
	configDef := &v2.Config{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
	}
	// CREATE: -->  Config
	conf, err := getProj.AddConfig(context.TODO(), configDef)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Config: %s ----> Created \n", conf.DisplayName())
	// GET: -->  Config [Singleton]
	getConfig, err := getProj.GetConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	printdata("Get Config hash name:", getConfig.Name)
	// GlobalNamespace Node Definition
	fmt.Println("Creating Globalnamespace object...")
	gDef := &v2.GlobalNamespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
	}
	// CREATE: --> GlobalNamespace
	globalN, err := getConfig.AddGlobalNamespace(context.TODO(), gDef)
	if err != nil {
		panic(err)
	}
	fmt.Printf("GlobalNamespace: %s ----> Created \n", globalN.DisplayName())
	// GET: --> GlobalNamespace
	getGlobalNamespace, err := getConfig.GetGlobalNamespace(context.TODO())
	if err != nil {
		panic(err)
	}
	printdata("Get GlobalNamespace hash name:", getGlobalNamespace.Name)
}

func getK8sAPIEndpointConfig() *rest.Config {
	var (
		host, kubeconfig *string
		kubeconfigHome   string
		config           *rest.Config
		err              error
	)

	host = flag.String("host", "", "portfowarded host to reach the app")
	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	if home := homedir.HomeDir(); home != "" {
		kubeconfigHome = filepath.Join(home, ".kube", "config")
	}

	flag.Parse()

	if len(*host) > 0 {
		fmt.Println("Connecting to k8s API at host: ", *host)
		config = &rest.Config{
			Host: *host,
		}
	} else if len(*kubeconfig) > 0 {
		fmt.Println("Connecting to k8s API in kubeconfig file: ", kubeconfigHome)
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	} else if len(kubeconfigHome) > 0 {
		fmt.Println("Connecting to k8s API in kubeconfig in home dir: ", kubeconfigHome)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigHome)
		if err != nil {
			panic(err.Error())
		}
	} else {
		fmt.Println("Unable to determing k8s API server endpoint. Exiting application.")
		os.Exit(1)
	}

	return config
}

func printdata(title string, data interface{}) {
	f, err := json.MarshalIndent(data, "", "        ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s: %v\n", title, string(f))
}
