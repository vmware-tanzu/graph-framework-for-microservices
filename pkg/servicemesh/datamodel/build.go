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

	if DatatmodelName != "" {
		envList = append(envList, fmt.Sprintf("DATAMODEL=%s", DatatmodelName))
		if err := utils.CheckDatamodelDirExists(DatatmodelName); err != nil {
			return err
		}
	}

	if err := utils.GoToNexusDirectory(); err != nil {
		return err
	}
	err := os.Chdir(DatatmodelName)
	if err != nil {
		return err
	}
	err = utils.SystemCommand(envList, "make", "datamodel_build")
	if err != nil {
		return fmt.Errorf("datamodel %s build failed with error %v", DatatmodelName, err)

	}
	return nil
}

func init() {
	BuildCmd.Flags().StringVarP(&DatatmodelName, "name", "n", "", "name of the datamodel to be build")
}
