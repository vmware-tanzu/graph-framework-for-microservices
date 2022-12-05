package main

import (
	"flag"
	"fmt"
	"graph-framework-for-microservices/install-validator/internal/dir"
	"graph-framework-for-microservices/install-validator/internal/kubernetes"
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

	// check for incompatible models and not installed. Panic if any and force != true
	incNames, _, text, err := dir.CheckDir(*directory, c)
	if err != nil {
		panic(err)
	}
	if len(incNames) > 0 && *force == false {
		fmt.Println("Changes detected. If you want to install models anyway, run with -force=true")
		panic(text)
	}

	// check if any data for incompatible models and panic if so
	var dataExist []string
	for _, n := range incNames {
		res, err := c.ListResources(*c.GetCrd(n))
		fmt.Println(res)
		if err != nil {
			panic(err)
		}
		if len(res) > 0 {
			dataExist = append(dataExist, n)
		}
	}
	if len(dataExist) > 0 {
		panic(fmt.Sprintf("There are some data that exist in datamodels that are backward incompatible: %v. Please remove them manually to force upgrade CRDs", dataExist))
	}

	// upsert all the models
	err = dir.InstallDir(*directory, c)
	if err != nil {
		panic(err)
	}

}
