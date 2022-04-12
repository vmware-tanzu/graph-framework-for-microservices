package runtime

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

var Namespace string

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs the Nexus runtime on the specified namespace",
	//Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: Install,
}

func Install(cmd *cobra.Command, args []string) error {
	envList := []string{}

	if Namespace != "" {
		envList = append(envList, fmt.Sprintf("NAMESPACE=%s", Namespace))
	}

	if err := utils.GoToNexusDirectory(); err != nil {
		return err
	}

	err := utils.SystemCommand(envList, !utils.IsDebug(cmd), "make", "runtime_install")
	if err != nil {
		return fmt.Errorf("runtime install failed with error %v", err)

	}
	fmt.Printf("\u2713 Runtime %s install successful\n", Namespace)

	return nil
}

func init() {
	InstallCmd.Flags().StringVarP(&Namespace, "namespace",
		"n", "", "name of the namespace to be created")
}
