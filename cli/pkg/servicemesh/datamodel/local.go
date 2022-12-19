package datamodel

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/common"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/prereq"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/utils"
)

func InstallLocal(cmd *cobra.Command, args []string) error {
	Namespace = cmd.Flags().Lookup(NamespaceFlag).Value.String()
	DatamodelName = args[0]
	if DatamodelName == "" {
		return fmt.Errorf("Please provide datamodel name to install")
	}
	if utils.ListPrereq(cmd) {
		prereq.PreReqListOnDemand(prerequisites)
		return nil
	}

	if err := utils.GoToNexusDirectory(); err != nil {
		return err
	}
	envList := common.GetEnvList()
	if DatamodelName != "" {
		envList = append(envList, fmt.Sprintf("DATAMODEL=%s", DatamodelName))
		if exists, err := utils.CheckDatamodelDirExists(DatamodelName); !exists {
			return utils.GetCustomError(utils.DATAMODEL_INSTALL_FAILED, err).Print().ExitIfFatalOrReturn()
		}
	} else {
		return utils.GetCustomError(utils.DATAMODEL_DIRECTORY_MISMATCH,
			fmt.Errorf("Please provide datamodel name using --name option when running from app directory")).Print().ExitIfFatalOrReturn()
	}
	envList = append(envList, fmt.Sprintf("NAMESPACE=%s", Namespace))

	err := utils.SystemCommand(cmd, utils.DATAMODEL_INSTALL_FAILED, envList, "make", "datamodel_install")
	if err != nil {
		return err
	}
	if err := InstallJob(ToolsImage, DatamodelName, "", Namespace, "true", Title, Force); err != nil {
		return err
	}
	return nil
}

var NameCmd = &cobra.Command{
	Use:   "name",
	Short: "Installing app local datamodel present in nexus/ folder",
	Args:  cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if utils.ListPrereq(cmd) {
			return nil
		}

		if utils.SkipPrereqCheck(cmd) {
			return nil
		}
		return prereq.PreReqVerifyOnDemand(installPrerequisites)
	},
	RunE: InstallLocal,
}
