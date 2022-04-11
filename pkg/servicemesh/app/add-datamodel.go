package app

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	DatamodelName string
	DatamodelUrl  string
	IsDefault     bool
)

func DatamodelAdd(cmd *cobra.Command, args []string) error {
	err := WriteToNexusDms(DatamodelName, NexusDmProps{DatamodelUrl, IsDefault})
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
	AddDatamodelCmd.Flags().StringVarP(&DatamodelUrl, "location", "l", "", "datamodel location (git URL)")
	AddDatamodelCmd.Flags().BoolVarP(&IsDefault, "default", "", false, "determines if the DM must be used by default")
}
