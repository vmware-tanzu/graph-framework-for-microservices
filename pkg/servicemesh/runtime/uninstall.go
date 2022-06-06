package runtime

import (
	"fmt"
	"os"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/version"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/prereq"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/logging"
)

func Uninstall(cmd *cobra.Command, args []string) error {
	for index, manifest := range common.RuntimeManifests {
		if manifest.Templatized {
			manifest.Image = common.ImageTemplate{
				Image:                fmt.Sprintf("%s/%s", Registry, manifest.ImageName),
				IsImagePullSecret:    false,
				ImagePullSecret:      ImagePullSecret,
				Namespace:            Namespace,
				NetworkingAPIVersion: NetworkingAPIVersion,
			}
		}
		common.RuntimeManifests[index] = manifest
	}
	var versions version.NexusValues
	if err := version.GetNexusValues(&versions); err != nil {
		return utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("could not download the runtime manifests due to %s", err)).Print().ExitIfFatalOrReturn()
	}

	directories, err := DownloadRuntimeFiles(cmd, versions)
	if err != nil {
		return err
	}
	var files []string
	for _, dir := range directories {
		files, err = GetFiles(files, dir)
		if err != nil {
			return err
		}
	}
	for _, file := range files {
		fmt.Printf("Deleting objects from file: %s\n", file)
		utils.SystemCommand(cmd, utils.RUNTIME_UNINSTALL_FAILED, []string{}, "kubectl", "delete", "-f", file, "-n", Namespace, "--ignore-not-found=true")
	}
	utils.SystemCommand(cmd, utils.RUNTIME_UNINSTALL_FAILED, []string{}, "kubectl", "delete", "pvc", "-n", Namespace, "-lapp=nexus-etcd")
	fmt.Printf("\u2713 Runtime %s uninstallation successful\n", Namespace)
	for _, manifest := range common.RuntimeManifests {
		os.RemoveAll(manifest.Directory)
		os.RemoveAll(manifest.FileName)
	}
	return nil

}

var UninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstalls the Nexus runtime from the specified namespace",
	//Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		err = prereq.PreReqVerifyOnDemand(prerequisites)
		if err != nil {
			return err
		}
		NetworkingAPIVersion, err = utils.GetNetworkingIngressVersion()
		if err != nil {
			return err
		}
		return nil
	},
	RunE: Uninstall,
}

func init() {
	UninstallCmd.Flags().StringVarP(&Namespace, "namespace",
		"n", "", "name of the namespace to be created")
	err := cobra.MarkFlagRequired(UninstallCmd.Flags(), "namespace")
	if err != nil {
		logging.Debugf("Runtime uninstall err: %v", err)
	}
}
