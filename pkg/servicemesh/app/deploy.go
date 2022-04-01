package app

import (
	"github.com/spf13/cobra"

	"gitlab.eng.vmware.com/nexus/cli/pkg/utils"
)

func Deploy(cmd *cobra.Command, args []string) error {
	envList := []string{}
	err := utils.SystemCommand(envList, "make", "app_deploy")
	if err != nil {
		return err
	}
	return nil
}

var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploys the application",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: Deploy,
}
