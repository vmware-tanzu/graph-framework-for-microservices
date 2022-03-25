package datamodel

import (
	"fmt"

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
	}

	err := utils.SystemCommand(envList, "make", "datamodel_build")
	if err != nil {
		return err
	}
	return nil
}

func init() {
	BuildCmd.Flags().StringVarP(&DatatmodelName, "name", "n", "", "name of the datamodel to be build")
}
