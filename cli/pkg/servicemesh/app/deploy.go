package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/prereq"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

var (
	isKindCluster bool
	namespace     string
)

var appDeployPrereqs = []prereq.Prerequiste{
	prereq.KUBERNETES,
	prereq.KUBERNETES_VERSION,
}

func Deploy(cmd *cobra.Command, args []string) error {
	envList := common.GetEnvList()
	if isKindCluster {
		envList = append(envList, fmt.Sprintf("CLUSTER=kind"))
	}
	if namespace != "" {
		envList = append(envList, fmt.Sprintf("NAMESPACE=%s", namespace))
	}
	if imageRegistry != "" {
		envList = append(envList, fmt.Sprintf("IMAGE_REGISTRY=%s", imageRegistry))
	}
	envList = append(envList, fmt.Sprintf("IMAGE_TAG=%s", imageTag))

	err := utils.SystemCommand(cmd, utils.APPLICATION_DEPLOY_FAILED, envList, "make", "deploy")
	if err != nil {
		return err
	}
	return nil
}

var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the application",
	RunE:  Deploy,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if utils.ListPrereq(cmd) {
			prereq.PreReqListOnDemand(appDeployPrereqs)
			return nil
		}
		if utils.SkipPrereqCheck(cmd) {
			return nil
		}
		return prereq.PreReqVerifyOnDemand(appDeployPrereqs)
	},
}

func init() {
	DeployCmd.Flags().BoolVarP(&isKindCluster, "kind", "", false, "indicates deployment to a kind cluster")

	DeployCmd.Flags().StringVarP(&imageRegistry, "registry", "r", "", "the image registry to deploy from")

	DeployCmd.Flags().StringVarP(&imageTag, "tag", "t", "", "the tag of the image to deploy")
	_ = cobra.MarkFlagRequired(DeployCmd.Flags(), "tag")

	DeployCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "the namespace to deploy to")
}
