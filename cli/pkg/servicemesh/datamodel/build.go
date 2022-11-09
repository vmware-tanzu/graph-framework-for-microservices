package datamodel

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/log"

	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/common"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/prereq"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/utils"
)

var prerequisites []prereq.Prerequiste = []prereq.Prerequiste{
	prereq.DOCKER,
}

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Compiles the Nexus DSLs into consumable APIs, CRDs, etc.",
	RunE:  Build,
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if utils.ListPrereq(cmd) {
			return nil
		}

		if utils.SkipPrereqCheck(cmd) {
			return nil
		}

		return prereq.PreReqVerifyOnDemand(prerequisites)
	},
}

func Build(cmd *cobra.Command, args []string) error {

	if utils.ListPrereq(cmd) {
		prereq.PreReqListOnDemand(prerequisites)
		return nil
	}
	compilerVersion, err := utils.GetTagVersion("Nexus", "NEXUS_DATAMODEL_COMPILER_VERSION")
	if err != nil {
		return utils.GetCustomError(utils.DATAMODEL_BUILD_FAILED, fmt.Errorf("could not get compiler Version information due to %s", err)).Print().ExitIfFatalOrReturn()
	}
	log.Debugf("Using compiler Version: %s\n", compilerVersion)

	envList := common.GetEnvList()
	envList = append(envList, fmt.Sprintf("TAG=%s", compilerVersion))
	containerID := os.Getenv("CONTAINER_ID")
	if containerID != "" {
		envList = append(envList, fmt.Sprintf("CONTAINER_ID=%s", containerID))
	}
	// hack for running datamodel build locally
	err = utils.SystemCommand(cmd, utils.CHECK_CURRENT_DIRECTORY_IS_DATAMODEL, envList, "make", "datamodel_build", "--dry-run")
	if err == nil {
		fmt.Println("Runnig build from current directory as this is a common datamodel.")
		err = utils.SystemCommand(cmd, utils.DATAMODEL_BUILD_FAILED, envList, "make", "datamodel_build")
		if err != nil {
			return fmt.Errorf("datamodel %s build failed with error %v", DatamodelName, err)
		} else {
			fmt.Println("\u2713 Datamodel build successful\n")
			return nil
		}
	}
	if err := utils.GoToNexusDirectory(); err != nil {
		pwd, _ := os.Getwd()
		log.Debugf("%s directory not found. Assuming %s to be datamodel directory\n", common.NEXUS_DIR, pwd)
	}
	if DatamodelName != "" {
		if exists, err := utils.CheckDatamodelDirExists(DatamodelName); !exists {
			return utils.GetCustomError(utils.DATAMODEL_DIRECTORY_NOT_FOUND, err).Print().ExitIfFatalOrReturn()
		}
	} else {
		return utils.GetCustomError(utils.DATAMODEL_DIRECTORY_MISMATCH,
			fmt.Errorf("Please provide datamodel name using --name option when running from app directory")).Print().ExitIfFatalOrReturn()
	}

	err = os.Chdir(DatamodelName)
	if err != nil {
		return err
	}
	err = utils.SystemCommand(cmd, utils.DATAMODEL_BUILD_FAILED, envList, "make", "datamodel_build")
	if err != nil {
		return fmt.Errorf("datamodel %s build failed with error %v", DatamodelName, err)

	}
	fmt.Printf("\u2713 Datamodel %s build successful\n", DatamodelName)
	return nil
}

func init() {
	BuildCmd.Flags().StringVarP(&DatamodelName, "name", "n", "", "name of the datamodel to be build")
}
