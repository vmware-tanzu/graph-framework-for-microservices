package app

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/common"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/prereq"

	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/utils"
)

var appPublishPrereqs = []prereq.Prerequiste{
	prereq.DOCKER,
}

func Publish(cmd *cobra.Command, args []string) error {
	envList := common.GetEnvList()
	if imageRegistry != "" {
		envList = append(envList, fmt.Sprintf("IMAGE_REGISTRY=%s", imageRegistry))
	}
	envList = append(envList, fmt.Sprintf("IMAGE_TAG=%s", imageTag))
	err := utils.SystemCommand(cmd, utils.APPLICATION_PUBLISH_FAILED, envList, "make", "publish")
	if err != nil {
		return err
	}
	return nil
}

var PublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish the Nexus application as a docker image",
	RunE:  Publish,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if utils.ListPrereq(cmd) {
			prereq.PreReqListOnDemand(appPublishPrereqs)
			return nil
		}
		if utils.SkipPrereqCheck(cmd) {
			return nil
		}
		return prereq.PreReqVerifyOnDemand(appPublishPrereqs)
	},
}

var (
	imageRegistry string
	imageTag      string
)

func init() {
	PublishCmd.Flags().StringVarP(&imageRegistry, "registry", "r", "", "the image registry to publish to")

	PublishCmd.Flags().StringVarP(&imageTag, "tag", "t", "", "the tag of the image to be published")
	_ = cobra.MarkFlagRequired(PublishCmd.Flags(), "tag")
}
