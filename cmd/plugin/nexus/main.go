package main

import (
	"os"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nexus/cli/pkg/servicemesh"
)

var rootCmd = &cobra.Command{
	Use:   "nexus",
	Short: "nexus cli",
	Long:  "nexus cli to execute datamodel, runtime and application related operations",
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
		//servicemesh.ConfigCmd,
		servicemesh.ApplyCmd,
		servicemesh.DeleteCmd,
		servicemesh.LoginCmd,
		servicemesh.RuntimeCmd,
		servicemesh.DataModelCmd,
		servicemesh.AppCmd,
		servicemesh.VersionCmd,
	)
}
