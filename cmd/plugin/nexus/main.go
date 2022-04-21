package main

import (
	"os"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/prereq"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/upgrade"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

var enableDebug = false
var prereqShow = false
var skipPrereqCheck = false

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
	)

	rootCmd.PersistentFlags().BoolVarP(&enableDebug, utils.EnableDebugFlag, "", false, "Enables extra logging")
	rootCmd.PersistentFlags().BoolVarP(&prereqShow, utils.ListPrereqFlag, "", false, "List prerequisites")
	rootCmd.PersistentFlags().BoolVarP(&skipPrereqCheck, utils.SkipPrereqCheckFlag, "", false, "Skip prerequisites check")
}
