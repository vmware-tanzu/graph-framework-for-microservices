package main

import (
	"flag"

	nexuscompare "github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus-compare"
	"github.com/vmware-tanzu/graph-framework-for-microservices/install-validator/internal/dir"
	"github.com/vmware-tanzu/graph-framework-for-microservices/install-validator/internal/kubernetes"
)

func main() {
	// get flags from cmd
	directory := flag.String("dir", "", "Directory with crds to install")
	force := flag.Bool("force", false, "Force install crds")
	flag.Parse()
	if *directory == "" {
		panic("No directory with crds to install. Provide it with -dir=$DIRECTORY ")
	}

	// setup k8s client
	c, _ := kubernetes.NewClient()
	err := c.ListCrds()
	if err != nil {
		panic(err)
	}

	err = dir.ApplyDir(*directory, *force, &c, nexuscompare.CompareFiles)
	if err != nil {
		panic(err)
	}

}
