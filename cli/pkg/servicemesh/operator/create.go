package operator

import (
	"fmt"

	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/common"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/app"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/prereq"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/utils"
)

var prerequisites []prereq.Prerequiste = []prereq.Prerequiste{
	prereq.GOLANG_VERSION,
}
var (
	CrdGroup              string
	CrdVersion            string
	CrdKind               string
	crdDatamodelBuildPath string
	CrdDatamodel          string
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
		crdDatamodelBuildPath = defaultDM.BuildDirectory
		fmt.Printf("Using default DM %v\n", defaultDM)
	}
	envList := common.GetEnvList()
	envList = append(envList, fmt.Sprintf("CRD_GROUP=%s", CrdGroup))
	envList = append(envList, fmt.Sprintf("CRD_VERSION=%s", CrdVersion))
	envList = append(envList, fmt.Sprintf("CRD_KIND=%s", CrdKind))
	envList = append(envList, fmt.Sprintf("CRD_DATAMODEL_NAME=%s", CrdDatamodel))
	if crdDatamodelBuildPath != "" {
		envList = append(envList, fmt.Sprintf("CRD_DATAMODEL_BUILD_DIRECTORY=%s", crdDatamodelBuildPath))
	}

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

	CreateCmd.Flags().StringVarP(&crdDatamodelBuildPath, "build-path",
		"b", "", "Build directory of CRDs and clients")
	if err != nil {
		fmt.Printf("init error: %v\n", err)
	}
}
