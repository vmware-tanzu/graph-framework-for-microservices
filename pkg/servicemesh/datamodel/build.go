package datamodel

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Compiles the Nexus DSLs into consumable APIs, CRDs, etc.",
	RunE:  Build,
}

func Build(cmd *cobra.Command, args []string) error {
	envList := []string{}

	// check if build can be run from current directory, if not proceed to next steps..
	err := utils.SystemCommand(envList, true, "make", "datamodel_build", "-n")
	if err == nil {
		fmt.Printf("Running build from current directory.\n")
		err = utils.SystemCommand(envList, false, "make", "datamodel_build")
		if err != nil {
			fmt.Printf("Error in building datamodel\n")
			return err
		}
		return nil
	}
	if err := utils.GoToNexusDirectory(); err != nil {
		return err
	}
	if DatamodelName != "" {
		if err := utils.CheckDatamodelDirExists(DatamodelName); err != nil {
			return err
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
	err = utils.SystemCommand(envList, false, "make", "datamodel_build")
	if err != nil {
		return fmt.Errorf("datamodel %s build failed with error %v", DatamodelName, err)

	}
	fmt.Printf("\u2713 Datamodel %s build successful\n", DatamodelName)
	return nil
}

func init() {
	BuildCmd.Flags().StringVarP(&DatamodelName, "name", "n", "", "name of the datamodel to be build")
}
