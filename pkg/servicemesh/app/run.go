package app

import (
	"fmt"

	"github.com/spf13/cobra"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

var (
	Namespace string
)

func Run(cmd *cobra.Command, args []string) error {
	envList := []string{}
	if Namespace != "" {
		envList = append(envList, fmt.Sprintf("NAMESPACE=%s", Namespace))
	}
	// cd nexus/
	err := utils.SystemCommand(cmd, utils.APPLICATION_RUN_FAILED, envList, "make", "app_run")
	if err != nil {
		return err
	}
	return nil
}

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the Nexus application after deploying to the specified namespace",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: Run,
}

func init() {
	RunCmd.Flags().StringVarP(&Namespace, "namespace",
		"n", "", "namespace name")
}
