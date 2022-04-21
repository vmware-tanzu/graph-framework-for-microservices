package datamodel

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/prereq"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

var Namespace string

var installPrerequisites []prereq.Prerequiste = []prereq.Prerequiste{
	prereq.KUBERNETES,
	prereq.KUBERNETES_VERSION,
}

func Install(cmd *cobra.Command, args []string) error {
	if utils.ListPrereq(cmd) {
		prereq.PreReqListOnDemand(prerequisites)
		return nil
	}

	if err := utils.GoToNexusDirectory(); err != nil {
		return err
	}
	envList := common.EnvList
	if DatamodelName != "" {
		envList = append(envList, fmt.Sprintf("DATAMODEL=%s", DatamodelName))
		if exists, err := utils.CheckDatamodelDirExists(DatamodelName); !exists {
			return err
		}
	} else {
		DatamodelName, err := utils.GetCurrentDatamodel()
		if err != nil {
			return err
		}
		fmt.Printf("Installing datamodel %s\n", DatamodelName)
		envList = append(envList, fmt.Sprintf("DATAMODEL=%s", DatamodelName))
	}

	if Namespace != "" {
		envList = append(envList, fmt.Sprintf("NAMESPACE=%s", Namespace))
	}

	err := utils.SystemCommand(cmd, utils.DATAMODEL_INSTALL_FAILED, envList, "make", "datamodel_install")
	if err != nil {
		return err
	}

	fmt.Printf("\u2713 Datamodel %s install successful\n", DatamodelName)
	return nil
}

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install specified datamodel's generated CRDs to the specified namespace",
	//Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if utils.ListPrereq(cmd) {
			return nil
		}

		if utils.SkipPrereqCheck(cmd) {
			return nil
		}
		return prereq.PreReqVerifyOnDemand(installPrerequisites)
	},
	RunE: Install,
}

func init() {
	InstallCmd.Flags().StringVarP(&Namespace, "namespace",
		"r", "", "name of the namespace to install to")
	InstallCmd.Flags().StringVarP(&DatamodelName, "name",
		"n", "", "name of the datamodel to install")

}
