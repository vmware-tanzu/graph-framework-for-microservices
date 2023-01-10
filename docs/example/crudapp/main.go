package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	hrv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/datamodel-examples.git/org-chart/build/apis/hr.vmware.org/v1"
	managementv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/datamodel-examples.git/org-chart/build/apis/management.vmware.org/v1"
	orgchartv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/datamodel-examples.git/org-chart/build/apis/orgchart.vmware.org/v1"
	rolev1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/datamodel-examples.git/org-chart/build/apis/role.vmware.org/v1"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/datamodel-examples.git/org-chart/build/nexus-client"
)

func main() {

	config := getK8sAPIEndpointConfig()
	// Create a datamodel client handle.
	nexusClient, err := nexus_client.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	crudWithNexusClient(nexusClient)
}

func crudWithNexusClient(nexusClient *nexus_client.Clientset) {
	cleanup(nexusClient)
	fmt.Printf("\nGoal of this application is to show how to build an object graph using Nexus Client and " +
		"do Get, Add, Update, Remove, Link, Unlink operations on the graph.\n")
	prompt()
	fmt.Printf("During run of this application we will create following graph:\n\n" +
		"           Root\n" +
		"             |   \\_________________________________\n" +
		"             |                         |            \\\n" +
		"             |                         |             |\n" +
		"            CEO~~~~~~~~~softlink~~~~ExecutiveRole  EmployeeRole\n" +
		"       (type Leader)              (type Executive) (type Empoloyee)\n" +
		"      /        |     \\___                            |\n" +
		"      |        |          \\                          |\n" +
		"   Manager1   Manager2    HR1~~~~~~softlink~~~~~~~~~~~\n" +
		"  (type Mgr) (type Mgr) (type HR)\n\n" +
		"Root and CEO objects are singletons\n\n")
	prompt()

	fmt.Println("-----------CREATE AND GET OPERATIONS-------------------------")
	fmt.Println("At first we will create base root node using default as name...")
	prompt()
	rootDef := &orgchartv1.Root{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
	}
	fmt.Println("Creating root object...")
	root, err := nexusClient.AddOrgchartRoot(context.TODO(), rootDef)
	if err != nil {
		panic(err)
	}
	fmt.Println("... checking if Root is created properly, name should be hashed, original name is preserved in " +
		"nexus/display_name label")
	printdata("Root object Meta", root.ObjectMeta)
	prompt()
	fmt.Println("Let's also check how to Get object, data should be the same")
	getRoot, err := nexusClient.GetOrgchartRoot(context.TODO())
	if err != nil {
		panic(err)
	}
	printdata("Root object Meta", getRoot.ObjectMeta)
	prompt()
	fmt.Println("And one more way to get object, using it's hashed name, result should also be the same")
	getRoot, err = nexusClient.Orgchart().GetRootByName(context.TODO(), root.GetName())
	if err != nil {
		panic(err)
	}
	printdata("Root object Meta", getRoot.ObjectMeta)

	prompt()
	fmt.Println("Now we create Leader object which is child of Root")
	prompt()
	leaderDef := &managementv1.Leader{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
		Spec: managementv1.LeaderSpec{
			Designation: "Chief",
			EmployeeID:  1,
		},
	}
	fmt.Println("Creating leaders object...")
	ceo, err := root.AddCEO(context.TODO(), leaderDef)
	if err != nil {
		panic(err)
	}

	fmt.Println("Checking if CEO leader is created properly...")
	printdata(fmt.Sprintf("CEO object is %s", ceo.DisplayName()), ceo.Spec)
	fmt.Println("and checking if Root is updated properly with Leader child:")
	ceoChild, err := root.GetCEO(context.TODO())
	if err != nil {
		panic(err)
	}
	printdata(fmt.Sprintf("Child CEO of root Object is %s", ceoChild.DisplayName()), ceoChild.Spec)
	fmt.Println("we can also check CEO object like this")
	ceoChild, err = nexusClient.OrgchartRoot().GetCEO(context.TODO())
	if err != nil {
		panic(err)
	}
	printdata(fmt.Sprintf("Child CEO of root Object is %s", ceoChild.DisplayName()), ceoChild.Spec)

	prompt()
	fmt.Printf("After creating Root and CEO graph state is like this:\n\n" +
		"           Root\n" +
		"             |\n" +
		"             |\n" +
		"             |\n" +
		"            CEO\n")
	prompt()
	fmt.Println("Now Let's add some children objects to Leader, we'll create 2 EngManagers " +
		"(there can be multiple children of this type) and one HR (there can be only one child of this type).")
	prompt()
	mgrDef1 := &managementv1.Mgr{
		ObjectMeta: metav1.ObjectMeta{
			Name: "Manager1",
		},
		Spec: managementv1.MgrSpec{
			EmployeeID: 2,
		},
	}
	fmt.Println("Creating Manager1 object...")
	_, err = ceo.AddEngManagers(context.TODO(), mgrDef1)
	if err != nil {
		panic(err)
	}
	mgrDef2 := &managementv1.Mgr{
		ObjectMeta: metav1.ObjectMeta{
			Name: "Manager2",
		},
		Spec: managementv1.MgrSpec{
			EmployeeID: 3,
		},
	}
	fmt.Println("Creating Manager2 object, using different command right now to show how it works")
	manager2, err := nexusClient.OrgchartRoot().CEO().
		AddEngManagers(context.TODO(), mgrDef2)
	if err != nil {
		panic(err)
	}
	hrDef := &hrv1.HumanResources{
		ObjectMeta: metav1.ObjectMeta{
			Name: "HR1",
		},
		Spec: hrv1.HumanResourcesSpec{
			EmployeeID: 3,
		},
	}
	fmt.Println("Creating HR1 object...")
	_, err = ceo.AddHR(context.TODO(), hrDef)
	if err != nil {
		panic(err)
	}

	fmt.Println("Let's check if CEO's EngManagers Children objects are updated properly...")
	prompt()
	engManagers, err := ceo.GetAllEngManagers(context.TODO())
	if err != nil {
		panic(err)
	}
	for _, manager := range engManagers {
		printdata(fmt.Sprintf("Child %s of CEO Object is", manager.DisplayName()), manager.Spec)
	}
	hr, err := ceo.GetHR(context.TODO())
	if err != nil {
		panic(err)
	}
	printdata("Child HR of CEO Object is", hr.Spec)

	fmt.Println("Named Children can be also obtained by name")
	mgr, err := ceo.GetEngManagers(context.TODO(), "Manager1")
	if err != nil {
		panic(err)
	}
	printdata("Child Manager1 of CEO Object is", mgr.Spec)

	prompt()
	fmt.Printf("Now graph state is like this:\n\n" +
		"           Root\n" +
		"             |\n" +
		"             |\n" +
		"             |\n" +
		"            CEO\n" +
		"       (type Leader) \n" +
		"      /        |     \\___\n" +
		"      |        |          \\\n" +
		"   Manager1   Manager2    HR1\n" +
		"  (type Mgr) (type Mgr) (type HR)\n")
	prompt()
	fmt.Println("-----------UPDATE SPEC OPERATIONS-------------------------")
	fmt.Println("Let's update spec of HR object")
	prompt()
	hr.Spec.EmployeeID = 1000
	fmt.Println("Updating HR1 spec...")
	err = hr.Update(context.TODO())
	if err != nil {
		panic(err)
	}
	printdata("Object HR after an update has spec:", hr.Spec)
	prompt()

	fmt.Println("-----------UPDATE STATUS OPERATIONS-------------------------")
	fmt.Println("Let's update status of Leader object")
	prompt()
	fmt.Println("Updating Leader status...")
	s := managementv1.LeaderState{
		IsOnVacations:            true,
		DaysLeftToEndOfVacations: 2,
	}
	err = ceo.SetStatus(context.TODO(), &s)
	if err != nil {
		panic(err)
	}
	updatedStatus, err := ceo.GetStatus(context.TODO())
	if err != nil {
		panic(err)
	}
	printdata("Status of leader is now", updatedStatus)
	fmt.Println("Let's update status of Leader object again, using chaining this time")
	prompt()
	fmt.Println("Updating Leader status...")
	s = managementv1.LeaderState{
		IsOnVacations:            true,
		DaysLeftToEndOfVacations: 1,
	}
	err = nexusClient.OrgchartRoot().CEO().SetStatus(context.TODO(), &s)
	if err != nil {
		panic(err)
	}
	updatedStatus, err = nexusClient.OrgchartRoot().CEO().GetStatus(context.TODO())
	if err != nil {
		panic(err)
	}
	printdata("Status of leader is now", updatedStatus)
	fmt.Println("We can also clear Leader status...")
	prompt()
	err = ceo.ClearStatus(context.TODO())
	if err != nil {
		panic(err)
	}
	updatedStatus, err = ceo.GetStatus(context.TODO())
	if err != nil {
		panic(err)
	}
	printdata("Status of leader is now", updatedStatus)
	prompt()
	fmt.Println("-----------ADD SOFTLINKS OPERATION-------------------------")
	fmt.Println("Now we will add some soflinks to the objects")
	fmt.Println("Let's create role Executive object first...")
	roleDef := &rolev1.Executive{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ExecutiveRole",
		},
	}

	fmt.Println("Creating Executives Role object...")
	execRole, err := root.AddExecutiveRole(context.TODO(), roleDef)
	if err != nil {
		panic(err)
	}

	prompt()
	fmt.Printf("Now graph looks like this:\n\n" +
		"           Root\n" +
		"             |   \\_____________________\n" +
		"             |                         |            \n" +
		"             |                         |             \n" +
		"            CEO                     ExecutiveRole\n" +
		"       (type Leader)              (type Executive)\n" +
		"      /        |     \\___                            \n" +
		"      |        |          \\                          \n" +
		"   Manager1   Manager2    HR1\n" +
		"  (type Mgr) (type Mgr) (type HR)\n")
	prompt()

	fmt.Println("Let's add role softlink to CEO Leader object...")
	err = ceo.LinkRole(context.TODO(), execRole)
	if err != nil {
		panic(err)
	}
	getExecRole, err := ceo.GetRole(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Role of %s is: %s \n", hr.DisplayName(), getExecRole.DisplayName())
	prompt()

	fmt.Printf("And now graph looks like this:\n\n" +
		"           Root\n" +
		"             |   \\______________________\n" +
		"             |                         |            \n" +
		"             |                         |             \n" +
		"            CEO~~~~~~~~~softlink~~~~ExecutiveRole\n" +
		"       (type Leader)              (type Executive)\n" +
		"      /        |     \\___                            \n" +
		"      |        |          \\                          \n" +
		"   Manager1   Manager2    HR1\n" +
		"  (type Mgr) (type Mgr) (type HR)\n")
	prompt()

	fmt.Println()
	fmt.Println("Let's create employee role and add it to Hr object...")
	emDef := &rolev1.Employee{
		ObjectMeta: metav1.ObjectMeta{
			Name: "EmployeeRole",
		},
	}
	emRole, err := root.AddEmployeeRole(context.TODO(), emDef)
	if err != nil {
		panic(err)
	}

	err = hr.LinkRole(context.TODO(), emRole)
	if err != nil {
		panic(err)
	}
	hrRole, err := hr.GetRole(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Role of %s is: %s \n", hr.DisplayName(), hrRole.DisplayName())
	prompt()
	fmt.Printf("And now we have full graph:\n\n" +
		"           Root\n" +
		"             |   \\_________________________________\n" +
		"             |                         |            \\\n" +
		"             |                         |             |\n" +
		"            CEO~~~~~~~~~softlink~~~~ExecutiveRole  EmployeeRole\n" +
		"       (type Leader)              (type Executive) (type Empoloyee)\n" +
		"      /        |     \\___                            |\n" +
		"      |        |          \\                          |\n" +
		"   Manager1   Manager2    HR1~~~~~~softlink~~~~~~~~~~~\n" +
		"  (type Mgr) (type Mgr) (type HR)\n")
	prompt()
	fmt.Println("----------REMOVE SOFTLINK OPERATION-------------------")
	fmt.Println("Now let's remove role from HR")
	err = hr.UnlinkRole(context.TODO())
	if err != nil {
		panic(err)
	}
	hrRole, err = hr.GetRole(context.TODO())
	if err == nil {
		panic("Expected not found error, but it's nil instead")
	}
	fmt.Printf("Role of %s is: %v\n", hr.DisplayName(), hrRole)
	prompt()

	fmt.Printf("Graph looks like this:\n\n" +
		"           Root\n" +
		"             |   \\_________________________________\n" +
		"             |                         |            \\\n" +
		"             |                         |             |\n" +
		"            CEO~~~~~~~~~softlink~~~~ExecutiveRole  EmployeeRole\n" +
		"       (type Leader)              (type Executive) (type Empoloyee)\n" +
		"      /        |     \\___                            \n" +
		"      |        |          \\                          \n" +
		"   Manager1   Manager2    HR1\n" +
		"  (type Mgr) (type Mgr) (type HR)\n")
	prompt()
	fmt.Println("-----------DELETE-------------------------")
	fmt.Println("Now let's delete CEO")
	err = root.DeleteCEO(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Println("All child resources should be removed, for example HR and Managers")
	_, err = nexusClient.OrgchartRoot().CEO().GetHR(context.TODO(), "HR1")
	fmt.Println("We expect not found err:")
	fmt.Println(err)
	_, err = nexusClient.Management().GetMgrByName(context.TODO(), manager2.GetName())
	fmt.Println("We expect not found err:")
	fmt.Println(err)
	prompt()
	fmt.Printf("And now graph has only root and roles:\n\n" +
		"           Root\n" +
		"                 \\_________________________________\n" +
		"                                       |            \\\n" +
		"                                       |             |\n" +
		"                                  ExecutiveRole  EmployeeRole\n" +
		"                                (type Executive) (type Empoloyee)\n")
	prompt()
	fmt.Println("We can also delete objects directly like this")
	err = execRole.Delete(context.TODO())
	if err != nil {
		panic(err)
	}
	err = emRole.Delete(context.TODO())
	if err != nil {
		panic(err)
	}
	prompt()
	fmt.Println("So let's check if roles are nil, as we expect")
	execRole, err = root.GetExecutiveRole(context.TODO())
	fmt.Printf("Child ExecutiveRole of default root is %v\n", execRole)
	emRole, err = root.GetEmployeeRole(context.TODO())
	fmt.Printf("Child EmployeeRole of default root is %v\n", emRole)
	prompt()

	fmt.Println()
	fmt.Println("Thank you! After pressing return all objects will be cleaned up")
	prompt()
	cleanup(nexusClient)
}

// getK8sAPIEndpointConfig determines K8s API Server endpoint
//
// If "host" is specified in command line argument, connect to API server pointed to by it.
// If not, if "kubeconfig" file is provided as input, connect to API server pointed to by it.
// If not, attempt to read kubeconfig file from home directory and connect to API server pointed to by it.
// If none of these are available, then exit with error.
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

func cleanup(nexusClient *nexus_client.Clientset) {
	err := nexusClient.DeleteOrgchartRoot(context.TODO())
	if err != nil && !errors.IsNotFound(err) {
		panic(err)
	}
	return
}

func prompt() {
	fmt.Printf("\n-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}
