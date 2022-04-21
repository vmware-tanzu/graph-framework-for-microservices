package operator

import (
	"fmt"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/app"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/prereq"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

var prerequisites []prereq.Prerequiste = []prereq.Prerequiste{
	prereq.GOLANG_VERSION,
}
var (
	CrdGroup     string
	CrdVersion   string
	CrdKind      string
	CrdDatamodel string
)

func Create(cmd *cobra.Command, args []string) error {
	if utils.ListPrereq(cmd) {
		prereq.PreReqListOnDemand(prerequisites)
		return nil
	}

	if CrdDatamodel == "" {
		// check if we have a default datamodel configured in nexus-dms.yaml
		defaultDM, err := app.GetDefaultDm()
		if err != nil {
			fmt.Print("Please provide name of datamodel with --datamodel option or set a default DM using `nexus app add-datamodel`\n")
			return err
		}
		CrdDatamodel = defaultDM.Location
		fmt.Printf("Using default DM %v\n", defaultDM)
	}
	envList := common.GetEnvList()
	envList = append(envList, fmt.Sprintf("CRD_GROUP=%s", CrdGroup))
	envList = append(envList, fmt.Sprintf("CRD_VERSION=%s", CrdVersion))
	envList = append(envList, fmt.Sprintf("CRD_KIND=%s", CrdKind))
	envList = append(envList, fmt.Sprintf("CRD_DATAMODEL_NAME=%s", CrdDatamodel))

	// check if we are in the correct directory
	// TBD. for now, we run from PWD
	fmt.Println("Running add_operator from current directory")
	err := utils.SystemCommand(cmd, utils.APPLICATION_OPERATOR_CREATE_FAILED, envList, "make", "add_operator")
	if err != nil {
		return err
	}
	fmt.Printf("Successfully created operator for type %s/%s/%s\n", CrdGroup, CrdVersion, CrdKind)
	return nil
}

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "creates an operator that subscribes to changes to the specified resource",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if utils.ListPrereq(cmd) {
			return nil
		}

		if utils.SkipPrereqCheck(cmd) {
			return nil
		}
		return prereq.PreReqVerifyOnDemand(prerequisites)
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

	CreateCmd.Flags().StringVarP(&CrdDatamodel, "datamodel",
		"d", "", "Datamodel that contains the specified resource")

	if err != nil {
		fmt.Printf("init error: %v\n", err)
	}
}
