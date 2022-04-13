package app

import (
	"github.com/spf13/cobra"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

func Deploy(cmd *cobra.Command, args []string) error {
	envList := []string{}
	err := utils.SystemCommand(envList, !utils.IsDebug(cmd), "make", "app_deploy")
	if err != nil {
		return err
	}
	return nil
}

var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "(TBD) Deploy the application",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: Deploy,
}
