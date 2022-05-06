package runtime

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/prereq"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/version"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/logging"
)

var Namespace string
var Filename = "runtime-manifests.tar"
var ManifestsDir = "manifets-nexus-runtime"

var prerequisites []prereq.Prerequiste = []prereq.Prerequiste{
	prereq.KUBERNETES,
	prereq.KUBERNETES_VERSION,
}

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs the Nexus runtime on the specified namespace",
	//Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if utils.ListPrereq(cmd) {
			return nil
		}

		if utils.SkipPrereqCheck(cmd) {
			return nil
		}
		return prereq.PreReqVerifyOnDemand(prerequisites)
	},
	RunE: Install,
}

func CreateNs(Namespace string) error {
	createCmd := exec.Command("kubectl", "create", "namespace", Namespace, "--dry-run", "-oyaml")
	applyCmd := exec.Command("kubectl", "apply", "-f", "-")

	r, w := io.Pipe()
	createCmd.Stdout = w
	applyCmd.Stdin = r

	var b2 bytes.Buffer
	applyCmd.Stdout = &b2

	err := createCmd.Start()
	if err != nil {
		return err
	}
	err = applyCmd.Start()
	if err != nil {
		return err
	}
	err = createCmd.Wait()
	if err != nil {
		return err
	}
	w.Close()
	err = applyCmd.Wait()
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, &b2)
	return nil
}
func Install(cmd *cobra.Command, args []string) error {

	if utils.ListPrereq(cmd) {
		prereq.PreReqListOnDemand(prerequisites)
		return nil
	}

	checkNs := exec.Command("kubectl", "get", "ns", Namespace)
	err := checkNs.Run()
	if err != nil {
		err := CreateNs(Namespace)
		if err != nil {
			fmt.Printf("Namespace %s creation failure due to %s", Namespace, err)
			return err
		}
	}
	files, err := DownloadRuntimeFiles(cmd)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			utils.SystemCommand(cmd, utils.RUNTIME_INSTALL_FAILED, []string{}, "kubectl", "apply", "-f", filepath.Join(ManifestsDir, "runtime-manifests", file.Name()), "-n", Namespace)
		}
	}
	fmt.Println("Waiting for the Nexus runtime to come up...")
	for _, label := range common.PodLabels {
		utils.CheckPodRunning(cmd, utils.RUNTIME_INSTALL_FAILED, label, Namespace)
	}
	fmt.Printf("\u2713 Runtime installation successful on namespace %s\n", Namespace)
	os.Remove(Filename)
	os.RemoveAll(ManifestsDir)
	return nil
}

func init() {
	InstallCmd.Flags().StringVarP(&Namespace, "namespace",
		"n", "", "name of the namespace to be created")
	err := cobra.MarkFlagRequired(InstallCmd.Flags(), "namespace")
	if err != nil {
		logging.Debugf("Runtime install err: %v", err)
	}
}

func DownloadRuntimeFiles(cmd *cobra.Command) ([]fs.FileInfo, error) {
	var values version.NexusValues

	if err := version.GetNexusValues(&values); err != nil {
		return nil, utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("could not download the runtime manifests due to %s", err)).Print().ExitIfFatalOrReturn()
	}
	runtimeVersion := os.Getenv("NEXUS_RUNTIME_TEMPLATE_VERSION")
	if runtimeVersion == "" {
		runtimeVersion = values.NexusDatamodelTemplates.Version
	}
	if utils.IsDebug(cmd) {
		fmt.Printf("Using runtime manifests Version: %s\n", runtimeVersion)
	}
	err := utils.DownloadFile(fmt.Sprintf(common.RUNTIME_MANIFESTS_URL, runtimeVersion), Filename)
	if err != nil {
		return nil, utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("could not download the runtime manifests due to %s", err)).Print().ExitIfFatalOrReturn()
	}
	file, err := os.Open(Filename)
	if err != nil {
		return nil, utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("accessing downloaded runtime manifests directory failed dwith error: %s", err)).Print().ExitIfFatalOrReturn()
	}
	defer file.Close()
	fo, err := os.Stat(ManifestsDir)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
				fmt.Errorf("issues in checking runtime manifests directory due to %s", err)).Print().ExitIfFatalOrReturn()
		}
	}
	if fo != nil {
		err = os.RemoveAll(ManifestsDir)
		if err != nil {
			return nil, utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
				fmt.Errorf("could not remove %s directory due to %s", ManifestsDir, err)).Print().ExitIfFatalOrReturn()
		}
	}
	err = os.Mkdir(ManifestsDir, os.ModePerm)
	if err != nil {
		return nil, utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("could not create the runtime manifests directory due to %s", err)).Print().ExitIfFatalOrReturn()
	}
	err = utils.Untar(ManifestsDir, file)
	if err != nil {
		return nil, utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("unarchive of runtime manifests directory failed with error %s", err)).Print().ExitIfFatalOrReturn()
	}

	files, err := ioutil.ReadDir(filepath.Join(ManifestsDir, "runtime-manifests"))
	if err != nil {
		return nil, utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("accessing runtime manifests directory failed with error to %s", err)).Print().ExitIfFatalOrReturn()
	}
	return files, nil
}
