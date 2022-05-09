package runtime

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/prereq"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/logging"
)

func Uninstall(cmd *cobra.Command, args []string) error {
	runtimeDir, validationDir, err := DownloadRuntimeFiles(cmd)
	if err != nil {
		return err
	}
	var files []string
	files, err = GetFiles(files, runtimeDir)
	if err != nil {
		return err
	}
	files, err = GetFiles(files, validationDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		fmt.Printf("Deleting objects from file: %s\n", file)
		utils.SystemCommand(cmd, utils.RUNTIME_UNINSTALL_FAILED, []string{}, "kubectl", "delete", "-f", file, "-n", Namespace, "--ignore-not-found=true")
	}
	utils.SystemCommand(cmd, utils.RUNTIME_UNINSTALL_FAILED, []string{}, "kubectl", "delete", "pvc", "-n", Namespace, "-lapp=nexus-etcd")
	fmt.Printf("\u2713 Runtime %s uninstallation successful\n", Namespace)
	os.Remove(runtimeFilename)
	os.Remove(validationFilename)
	os.RemoveAll(RuntimeManifestsDir)
	os.RemoveAll(ValidationManifestDir)
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
