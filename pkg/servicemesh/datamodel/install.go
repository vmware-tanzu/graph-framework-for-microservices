package datamodel

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nexus/cli/pkg/utils"
)

var Namespace string

func Install(cmd *cobra.Command, args []string) error {
	envList := []string{}
	fmt.Print("Checking if the tenant-apiserver is reachable for installing datamodel crds\n")

	if DatatmodelName != "" {
		envList = append(envList, fmt.Sprintf("DATAMODEL=%s", DatatmodelName))
		if err := utils.CheckDatamodelDirExists(DatatmodelName); err != nil {
			return err
		}
	}

	if Namespace != "" {
		envList = append(envList, fmt.Sprintf("NAMESPACE=%s", Namespace))
	}

	fmt.Print("run this command outside of nexus home directory\n")
	if err := utils.GoToNexusDirectory(); err != nil {
		return err
	}

	err := utils.SystemCommand(envList, "make", "datamodel_install")
	if err != nil {
		return err
	}
	return nil
}

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "install namespace from directory",
	//Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: Install,
}

func init() {
	InstallCmd.Flags().StringVarP(&Namespace, "namespace",
		"n", "", "name of the namespace to be created")
	InstallCmd.Flags().StringVarP(&DatatmodelName, "datamodel",
		"d", "", "name of the datamodel to be build")

}
