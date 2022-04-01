package app

import (
	"github.com/spf13/cobra"

	"gitlab.eng.vmware.com/nexus/cli/pkg/utils"
)

func Package(cmd *cobra.Command, args []string) error {
	envList := []string{}
	err := utils.SystemCommand(envList, "make", "app_package")
	if err != nil {
		return err
	}
	return nil
}

var PackageCmd = &cobra.Command{
	Use:   "package",
	Short: "installs application package",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: Package,
}
