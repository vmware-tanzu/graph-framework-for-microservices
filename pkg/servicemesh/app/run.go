package app

import (
	"fmt"

	"github.com/spf13/cobra"

	"gitlab.eng.vmware.com/nexus/cli/pkg/utils"
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
	err := utils.SystemCommand(envList, "make", "app_run")
	if err != nil {
		return err
	}
	return nil
}

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "run the application",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: Run,
}

func init() {
	RunCmd.Flags().StringVarP(&Namespace, "namespace",
		"n", "", "namespace name")
}
