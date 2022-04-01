package datamodel

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"gitlab.eng.vmware.com/nexus/cli/pkg/utils"
)

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "clones policymodel and creates crds",
	RunE:  Build,
}

func Build(cmd *cobra.Command, args []string) error {
	envList := []string{}

	// check if build can be run from current directory, if not proceed to next steps..
	err := utils.SystemCommand(envList, true, "make", "datamodel_build", "-n")
	if err == nil {
		fmt.Printf("Running build from current directory.")
		err = utils.SystemCommand(envList, false, "make", "datamodel_build")
		if err != nil {
			fmt.Printf("Error in building datamodel")
			return err
		}
		return nil
	}
	if err := utils.GoToNexusDirectory(); err != nil {
		return err
	}
	if DatatmodelName != "" {
		if err := utils.CheckDatamodelDirExists(DatatmodelName); err != nil {
			return err
		}
	} else {
		_, err := os.Stat("NEXUSDATAMODEL")
		if err != nil {
			fmt.Printf("Please provdide --name of your datamodel to build")
			return nil
		}
		DatamodelNamebytes, _ := os.ReadFile("NEXUDATAMODEL")
		DatatmodelName = string(DatamodelNamebytes)
	}

	err = os.Chdir(DatatmodelName)
	if err != nil {
		return err
	}
	err = utils.SystemCommand(envList, false, "make", "datamodel_build")
	if err != nil {
		return fmt.Errorf("datamodel %s build failed with error %v", DatatmodelName, err)

	}
	return nil
}

func init() {
	BuildCmd.Flags().StringVarP(&DatatmodelName, "name", "n", "", "name of the datamodel to be build")
}
