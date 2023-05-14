package datamodel

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/prereq"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

func InstallRemote(cmd *cobra.Command, args []string) error {
	Namespace = cmd.Flags().Lookup(NamespaceFlag).Value.String()
	DatamodelImage = args[0]
	if DatamodelImage == "" {
		return fmt.Errorf("Please provide datamodel image path to install")
	}
	if utils.ListPrereq(cmd) {
		prereq.PreReqListOnDemand(prerequisites)
		return nil
	}
	if err := InstallJob(DatamodelImage, "", ImagePullSecret, Namespace, "false", Title); err != nil {
		return err
	}
	return nil
}

var ImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Remote datamodel installation which is pushed as docker image",
	Args:  cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if utils.ListPrereq(cmd) {
			return nil
		}

		if utils.SkipPrereqCheck(cmd) {
			return nil
		}
		return prereq.PreReqVerifyOnDemand(installPrerequisites)
	},
	RunE: InstallRemote,
}

func init() {
	ImageCmd.Flags().StringVarP(&ImagePullSecret, "secretname",
		"s", "", "secret to pull images on namespace - needs to be created by user")
}
