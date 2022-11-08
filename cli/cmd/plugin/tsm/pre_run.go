package main

import (
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/log"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/config"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/utils"
)

func RootPreRun(cmd *cobra.Command, args []string) error {
	nexusConfig := config.LoadNexusConfig()

	if nexusConfig.DebugAlways {
		cmd.Flags().Lookup("debug").Changed = true
	}

	// Set logging level to debug if debugging is enabled.
	if utils.IsDebug(cmd) {
		log.SetLevel(log.DebugLevel)
	}
	return nil
}
