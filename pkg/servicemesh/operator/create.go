package operator

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

var (
	CrdGroup     string
	CrdVersion   string
	CrdKind      string
	CrdDatamodel string
)

func Create(cmd *cobra.Command, args []string) error {
	if CrdDatamodel == "" {
		fmt.Print("Please provide name of datamodel with --import option\n")
		return nil
	}
	envList := append([]string{}, fmt.Sprintf("CRD_GROUP=%s", CrdGroup))
	envList = append(envList, fmt.Sprintf("CRD_VERSION=%s", CrdVersion))
	envList = append(envList, fmt.Sprintf("CRD_KIND=%s", CrdKind))
	envList = append(envList, fmt.Sprintf("CRD_DATAMODEL_NAME=%s", CrdDatamodel))

	// check if we are in the correct directory
	// TBD. for now, we run from PWD
	fmt.Println("Running add_operator from current directory")
	err := utils.SystemCommand(envList, true, "make", "add_operator")
	if err != nil {
		return err
	}
	fmt.Printf("Successfully created operator for type %s/%s/%s", CrdGroup, CrdVersion, CrdKind)
	return nil
}

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "creates an operator that subscribes to changes to the specified resource",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: Create,
}

func init() {
	var err error
	CreateCmd.Flags().StringVarP(&CrdGroup, "group",
		"g", "", "group of the CRD")
	err = cobra.MarkFlagRequired(CreateCmd.Flags(), "group")

	CreateCmd.Flags().StringVarP(&CrdVersion, "version",
		"v", "", "version of the CRD")
	err = cobra.MarkFlagRequired(CreateCmd.Flags(), "version")

	CreateCmd.Flags().StringVarP(&CrdKind, "kind",
		"k", "", "'kind' of the CRD")
	err = cobra.MarkFlagRequired(CreateCmd.Flags(), "kind")

	CreateCmd.Flags().StringVarP(&CrdDatamodel, "import",
		"r", "", "Location of the CRD (DM git repo URL)")

	if err != nil {
		fmt.Printf("init error: %v\n", err)
	}
}
