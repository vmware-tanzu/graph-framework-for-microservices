package app

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/common"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/prereq"

	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/utils"
)

var appBuildPrereqs = []prereq.Prerequiste{
	prereq.GOLANG_VERSION,
	prereq.DOCKER,
}

func Build(cmd *cobra.Command, args []string) error {
	envList := common.GetEnvList()
	if imageRegistry != "" {
		envList = append(envList, fmt.Sprintf("IMAGE_REGISTRY=%s", imageRegistry))
	}
	envList = append(envList, fmt.Sprintf("IMAGE_TAG=%s", imageTag))

	err := utils.SystemCommand(cmd, utils.APPLICATION_BUILD_FAILED, envList, "make", "build")
	if err != nil {
		return err
	}
	return nil
}

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds the application",
	RunE:  Build,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if utils.ListPrereq(cmd) {
			prereq.PreReqListOnDemand(appBuildPrereqs)
			return nil
		}
		if utils.SkipPrereqCheck(cmd) {
			return nil
		}
		return prereq.PreReqVerifyOnDemand(appBuildPrereqs)
	},
}

func init() {
	BuildCmd.Flags().StringVarP(&imageRegistry, "registry", "r", "", "the image registry used to name the image")

	BuildCmd.Flags().StringVarP(&imageTag, "tag", "t", "", "the tag to be given to the built image")
	_ = cobra.MarkFlagRequired(BuildCmd.Flags(), "tag")
}
