package main

import (
	"os"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

var enableDebug = false

var rootCmd = &cobra.Command{
	Use:               "tsm",
	Short:             "TSM CLI",
	Long:              "TSM CLI helps with features configurations and policies via Declarative Config",
	PersistentPreRunE: RootPreRun,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(
		servicemesh.ApplyCmd,
		servicemesh.DeleteCmd,
		servicemesh.LoginCmd,
		servicemesh.GetCmd,
		servicemesh.ConfigCmd,
		servicemesh.TSMVersionCmd,
	)

	rootCmd.PersistentFlags().BoolVarP(&enableDebug, utils.EnableDebugFlag, "", false, "Enables extra logging")
}
