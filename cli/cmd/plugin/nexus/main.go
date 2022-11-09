package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/prereq"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/upgrade"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/utils"
)

var enableDebug = false
var prereqShow = false
var skipPrereqCheck = false

var rootCmd = &cobra.Command{
	Use:               "nexus",
	Short:             "Nexus CLI",
	Long:              "The Nexus CLI to create and deploy Nexus datamodels and applications. Learn about the Nexus platform here - https://github.com/vmware-tanzu/graph-framework-for-microservices/docs",
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
		servicemesh.RuntimeCmd,
		servicemesh.DataModelCmd,
		servicemesh.AppCmd,
		servicemesh.OperatorCmd,
		servicemesh.VersionCmd,
		upgrade.UpgradeCmd,
		prereq.PreReqCmd,
		servicemesh.ConfigCmd,
		servicemesh.DebugCmd,
	)

	rootCmd.PersistentFlags().BoolVarP(&enableDebug, utils.EnableDebugFlag, "", false, "Enables extra logging")
	rootCmd.PersistentFlags().BoolVarP(&prereqShow, utils.ListPrereqFlag, "", false, "List prerequisites")
	rootCmd.PersistentFlags().BoolVarP(&skipPrereqCheck, utils.SkipPrereqCheckFlag, "", false, "Skip prerequisites check")
}
