package app

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	datamodelName string
)

func SetDefaultDatamodel(cmd *cobra.Command, args []string) error {
	err := SetDefaultDm(datamodelName)
	if err != nil {
		return fmt.Errorf("failed to set default DM")
	}
	return nil
}

var SetDefaultDatamodelCmd = &cobra.Command{
	Use:   "set-default-datamodel",
	Short: "Sets a default datamodel for the app. Any controllers created without specifying a DM explicitly will use the default DM",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: SetDefaultDatamodel,
}

func init() {
	SetDefaultDatamodelCmd.Flags().StringVarP(&datamodelName, "name", "n", "", "datamodel name")
}
