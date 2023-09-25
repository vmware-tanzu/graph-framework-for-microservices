package datamodel

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/log"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/prereq"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
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
	compilerVersion, err := utils.GetTagVersion("NexusCompiler", "NEXUS_DATAMODEL_COMPILER_VERSION")
	if err != nil {
		return utils.GetCustomError(utils.DATAMODEL_BUILD_FAILED, fmt.Errorf("could not get compiler Version information due to %s", err)).Print().ExitIfFatalOrReturn()
	}
	log.Debugf("Using compiler Version: %s\n", compilerVersion)

	envList := common.GetEnvList()
	envList = append(envList, fmt.Sprintf("COMPILER_TAG=%s", compilerVersion))
	containerID := os.Getenv("CONTAINER_ID")
	if containerID != "" {
		envList = append(envList, fmt.Sprintf("CONTAINER_ID=%s", containerID))
	}

	// hack for running datamodel build locally
	err = utils.SystemCommand(cmd, utils.CHECK_CURRENT_DIRECTORY_IS_DATAMODEL, envList, "make", "datamodel_build", "--dry-run")
	if err == nil {
		fmt.Println("Running build from current directory as this is a common datamodel.")
		err = utils.SystemCommand(cmd, utils.DATAMODEL_BUILD_FAILED, envList, "make", "datamodel_build")
		if err != nil {
			return fmt.Errorf("datamodel %s build failed with error %v", DatamodelName, err)
		} else {
			fmt.Println("\u2713 Datamodel build successful\n")

			err = updateNexusYaml()
			if err != nil {
				return fmt.Errorf("could not write to nexus.yaml: %v", err)
			}
			publishDockerImage(cmd)
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
	err = updateNexusYaml()
	if err != nil {
		return fmt.Errorf("could not write to nexus.yaml: %v", err)
	}
	publishDockerImage(cmd)
	return nil
}

func publishDockerImage(cmd *cobra.Command) error {
	if BuildDockerImg {
		envList := common.GetEnvList()
		envList = append(envList, fmt.Sprintf("DOCKER_REPO=%s", DatamodelName))
		envList = append(envList, fmt.Sprintf("VERSION=%s", "latest"))

		err := utils.SystemCommand(cmd, utils.DATAMODEL_DOCKER_IMAGE_BUILD_FAILED, envList, "make", "docker_build")
		if err != nil {
			return fmt.Errorf("docker image %s build failed with error %v", DatamodelName, err)

		}
		fmt.Printf("\u2713 Datamodel docker image %s:latest built successfully\n", DatamodelName)
	}
	return nil
}

// store the datamodel name in nexus.yaml and use it to figure out the datamodel name
func updateNexusYaml() error {
	if DatamodelName != "" {
		err := utils.SetDatamodelName(DatamodelName)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	BuildCmd.Flags().StringVarP(&DatamodelName, "name", "n", "", "name of the datamodel to be build")
	BuildCmd.Flags().BoolVarP(&BuildDockerImg, "dockerbuild", "b", true, "build docker image")
}
