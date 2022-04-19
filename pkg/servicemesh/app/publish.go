package app

import (
	"github.com/spf13/cobra"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

func Publish(cmd *cobra.Command, args []string) error {
	envList := []string{}
	err := utils.SystemCommand(cmd, utils.APPLICATION_PUBLISH_FAILED, envList, "make", "app_publish")
	if err != nil {
		return err
	}
	return nil
}

var PublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "(TBD) Publish the Nexus application as a docker image",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: Publish,
}
