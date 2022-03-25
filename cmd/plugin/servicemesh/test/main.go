package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu-private/core/pkg/v1/cli"
	"github.com/vmware-tanzu-private/core/pkg/v1/cli/command/plugin"
)

var pluginName = "servicemesh"

var descriptor = cli.NewTestFor(pluginName)

func main() {
	p, err := plugin.NewPlugin(descriptor)
	if err != nil {
		log.Fatal(err)
	}
	p.Cmd.RunE = test
	if err := p.Execute(); err != nil {
		os.Exit(1)
	}
}

func test(c *cobra.Command, _ []string) error {
	// TBD
	// create
	// delete
	// get
	// list

	return nil
}

// Cleanup the test.
// func Cleanup() error {
// 	return nil
// }
