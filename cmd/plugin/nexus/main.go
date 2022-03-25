package main

import (
	"os"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nexus/cli/pkg/servicemesh"
)

// var descriptor = cli.PluginDescriptor{
// 	Name:        "nexus",
// 	Description: "nexus features",
// 	Version:     "v0.0.1",
// 	BuildSHA:    "",
// 	Group:       cli.ManageCmdGroup,
// 	DocURL:      "",
// }

var rootCmd = &cobra.Command{
	Use:   "nexus",
	Short: "nexus cli",
	Long:  "nexus cli to execute tsm operations",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(
		servicemesh.ClusterCmd,
		servicemesh.GnsCmd,
		servicemesh.ConfigCmd,
		servicemesh.ApplyCmd,
		servicemesh.DeleteCmd,
		servicemesh.LoginCmd,
		servicemesh.RuntimeCmd,
		servicemesh.DataModelCmd,
		servicemesh.AppCmd,
	)
}
