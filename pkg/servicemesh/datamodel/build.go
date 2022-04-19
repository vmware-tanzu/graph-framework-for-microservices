package datamodel

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/version"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Compiles the Nexus DSLs into consumable APIs, CRDs, etc.",
	RunE:  Build,
}

func Build(cmd *cobra.Command, args []string) error {
	envList := []string{}

	var values version.NexusValues
	err := utils.IsDockerRunning(cmd)
	if err != nil {
		return fmt.Errorf("docker daemon doesn't seem to be running. Please retry after starting Docker\n")
	}

	yamlFile, err := common.TemplateFs.ReadFile("values.yaml")
	if err != nil {
		return fmt.Errorf("error while reading version yamlFile %v", err)

	}

	err = yaml.Unmarshal(yamlFile, &values)
	if err != nil {
		return fmt.Errorf("error while unmarshal version yaml data %v", err)
	}
	envList = append(envList, fmt.Sprintf("TAG=%s", values.NexusCompiler.Version))
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
		if utils.IsDebug(cmd) {
			pwd, _ := os.Getwd()
			fmt.Printf("%s directory not found. Assuming %s to be datamodel directory\n", common.NEXUS_DIR, pwd)
		}
	}
	if DatamodelName != "" {
		if exists, err := utils.CheckDatamodelDirExists(DatamodelName); !exists {
			return utils.GetCustomError(utils.DATAMODEL_DIRECTORY_NOT_FOUND, err).Print().ExitIfFatalOrReturn()
		}
	} else {
		DatamodelName, err = utils.GetCurrentDatamodel()
		if err != nil {
			return err
		}
		fmt.Printf("Running build for datamodel %s\n", DatamodelName)
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
