package app

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

var (
	DatamodelName      string
	DatamodelUrl       string
	IsDefault          bool
	DatamodelBuildPath string
)

func DatamodelAdd(cmd *cobra.Command, args []string) error {

	// Hack: Assuming that datamodel url will only be provided for
	// common datamodel. Interpreting that absence of URL is an
	// indication that its a local datamodel and 'make replace' is
	// to be invoked.
	//
	// Not reworking it anymore as we will migrate 'make replace'
	// as part of add operator workflow.
	if DatamodelUrl == "" {
		envList := []string{
			fmt.Sprintf("DATAMODEL=%s", DatamodelName),
		}

		utils.SystemCommand(cmd, utils.DATAMODEL_INIT_FAILED, envList, "make", "replace")
		fmt.Println("\u2713 Datamodel added to application successfully")
	}

	err := WriteToNexusDms(DatamodelName, NexusDmProps{DatamodelUrl, IsDefault, DatamodelBuildPath})
	if err != nil {
		return fmt.Errorf("failed to write to nexus-dms.yaml")
	}

	return nil
}

var AddDatamodelCmd = &cobra.Command{
	Use:   "add-datamodel",
	Short: "Add a datamodel reference to this app",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: DatamodelAdd,
}

func init() {
	AddDatamodelCmd.Flags().StringVarP(&DatamodelName, "name", "n", "", "datamodel name")
	AddDatamodelCmd.MarkFlagRequired("name")

	AddDatamodelCmd.Flags().StringVarP(&DatamodelUrl, "package-name", "p", "", "importable name for the datamodel package")
	AddDatamodelCmd.Flags().BoolVarP(&IsDefault, "default", "", false, "determines if the DM must be used by default")
	AddDatamodelCmd.Flags().StringVarP(&DatamodelBuildPath, "build-path", "b", "", "build directory where api clients are localted")
}
