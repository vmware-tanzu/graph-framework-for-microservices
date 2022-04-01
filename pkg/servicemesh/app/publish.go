package app

import (
	"github.com/spf13/cobra"

	"gitlab.eng.vmware.com/nexus/cli/pkg/utils"
)

func Publish(cmd *cobra.Command, args []string) error {
	envList := []string{}
	err := utils.SystemCommand(envList, "make", "app_publish")
	if err != nil {
		return err
	}
	return nil
}

var PublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish the application",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: Publish,
}
