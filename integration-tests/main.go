package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// === PREREQUISITES ===
// nexus CLI and its prerequisites have to be installed
// GOPATH needs to be set

// workflow name -> steps
var workflows map[string]Workflow

var pwd string

func main() {
	SetupEnv()

	RunCmd("nexus version")
	RunCmd("nexus config view")

	if len(os.Args) != 2 {
		fmt.Println("`go run main workflows.yaml")
	}
	workflowsFile := os.Args[1]
	pwd = os.Getenv("PWD")

	Parse(workflowsFile)

	for name, workflow := range workflows {
		fmt.Printf("BEGIN WORKFLOW %s\n", name)
		RunWorkflow(workflow)
		fmt.Printf("END WORKFLOW %s\n", name)

		// Cleanup cluster after each workflow
		fmt.Println("Begin cluster cleanup...")
		err := cleanupCluster()
		CheckIfError(err)
		fmt.Println("Cluster cleanup complete")
	}
}

func Parse(file string) {
	data, err := ioutil.ReadFile(file)
	if os.IsNotExist(err) {
		fmt.Printf("%s doesn't exist. Aborting\n", file)
		os.Exit(1)
	}
	err = yaml.Unmarshal(data, &workflows)
	if err != nil {
		fmt.Printf("Failed to read contents of %s due to error: %s\n", file, err)
		os.Exit(1)
	}
}
