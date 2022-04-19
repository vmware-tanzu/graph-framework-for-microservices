package app

import (
	"github.com/spf13/cobra"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

func Deploy(cmd *cobra.Command, args []string) error {
	err := utils.SystemCommand(cmd, utils.APPLICATION_DEPLOY_FAILED, common.EnvList, "make", "app_deploy")
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
