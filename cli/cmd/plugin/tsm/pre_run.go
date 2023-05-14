package main

import (
	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/log"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/config"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

func RootPreRun(cmd *cobra.Command, args []string) error {
	nexusConfig, err := config.LoadNexusConfig()
	if err != nil {
		return err
	}

	if nexusConfig.DebugAlways {
		cmd.Flags().Lookup("debug").Changed = true
	}

	// Set logging level to debug if debugging is enabled.
	if utils.IsDebug(cmd) {
		log.SetLevel(log.DebugLevel)
	}
	return nil
}
