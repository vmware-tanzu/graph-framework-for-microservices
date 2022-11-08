package app

import (
	"github.com/spf13/cobra"

	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/common"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/utils"
)

func Package(cmd *cobra.Command, args []string) error {
	envList := common.GetEnvList()
	err := utils.SystemCommand(cmd, utils.APPLICATION_PACKAGE_FAILED, envList, "make", "app_package")
	if err != nil {
		return err
	}
	return nil
}

var PackageCmd = &cobra.Command{
	Use:   "package",
	Short: "(TBD) Package the Nexus application ready for deployment",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: Package,
}
