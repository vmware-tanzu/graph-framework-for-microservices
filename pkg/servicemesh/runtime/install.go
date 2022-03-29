package runtime

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nexus/cli/pkg/utils"
)

const TENANT_INSTALLATION_MANIFEST = "runtime-manifests/deployment/"
const CRD_FOLDER = "runtime-manifests/crds/"

var Namespace string

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "install tenant from directory",
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

	fmt.Print("run this command outside of nexus home directory\n")
	if err := utils.GoToNexusDirectory(); err != nil {
		return err
	}

	err := utils.SystemCommand(envList, "make", "runtime_install")
	if err != nil {
		return fmt.Errorf("runtime install failed with error %v", err)

	}
	return nil
}

func init() {
	InstallCmd.Flags().StringVarP(&Namespace, "namespace",
		"n", "", "name of the namespace to be created")
}
