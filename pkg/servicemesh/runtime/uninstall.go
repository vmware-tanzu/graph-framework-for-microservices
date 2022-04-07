package runtime

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

func Uninstall(cmd *cobra.Command, args []string) error {
	envList := []string{}

	if Namespace != "" {
		envList = append(envList, fmt.Sprintf("NAMESPACE=%s", Namespace))
	}

	if err := utils.GoToNexusDirectory(); err != nil {
		return err
	}

	err := utils.SystemCommand(envList, false, "make", "runtime_uninstall")
	if err != nil {
		return fmt.Errorf("runtime install failed with error %v", err)

	}
	fmt.Printf("\u2713 Runtime %s install successful\n", Namespace)

	return nil
}

var UninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstalls the Nexus runtime from the specified namespace",
	//Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: Uninstall,
}

func init() {
	UninstallCmd.Flags().StringVarP(&Namespace, "namespace",
		"n", "", "name of the namespace to be created")
}
