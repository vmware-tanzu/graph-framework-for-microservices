package gns

import (
	"fmt"

	"github.com/spf13/cobra"
)

// DeleteCmd ...Deletes a global namespace
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a global namespace",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: delete,
}

var delResourceFile string

func init() {
	DeleteCmd.Flags().StringVarP(&delResourceFile, "file", "f", "", "Resource file to be applied.")
}

func delete(cmd *cobra.Command, args []string) error {
	fmt.Println("TBD: servicemesh delete()")
	return nil
}
