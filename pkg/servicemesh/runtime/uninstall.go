package runtime

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/prereq"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/logging"
)

func Uninstall(cmd *cobra.Command, args []string) error {
	directories, err := DownloadRuntimeFiles(cmd)
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
		return prereq.PreReqVerifyOnDemand(prerequisites)
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
