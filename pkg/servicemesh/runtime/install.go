package runtime

import (
	"bytes"
	"errors"
	"fmt"
	"io"
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
var runtimeFilename = "runtime-manifests.tar"
var validationFilename = "validation-manifets.tar"
var RuntimeManifestsDir = "manifests-nexus-runtime"
var ValidationManifestDir = "manifests-nexus-validation"

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

func GetFiles(Files []string, directory string) ([]string, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			Files, err = GetFiles(Files, filepath.Join(directory, file.Name()))
			if err != nil {
				return nil, err
			}
		} else {
			Files = append(Files, filepath.Join(directory, file.Name()))
		}
	}
	return Files, nil
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
	var files []string
	runtimeDir, validationDir, err := DownloadRuntimeFiles(cmd)
	if err != nil {
		return err
	}

	files, err = GetFiles(files, runtimeDir)
	if err != nil {
		return err
	}
	files, err = GetFiles(files, validationDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		fmt.Printf("Applying file: %s\n", file)
		err = utils.SystemCommand(cmd, utils.RUNTIME_INSTALL_FAILED, []string{}, "kubectl", "apply", "-f", file, "-n", Namespace)
		if err != nil {
			return err
		}
	}
	fmt.Println("Waiting for the Nexus runtime to come up...")
	for _, label := range common.PodLabels {
		utils.CheckPodRunning(cmd, utils.RUNTIME_INSTALL_FAILED, label, Namespace)
	}
	fmt.Printf("\u2713 Runtime installation successful on namespace %s\n", Namespace)
	os.Remove(runtimeFilename)
	os.Remove(validationFilename)
	os.RemoveAll(RuntimeManifestsDir)
	os.RemoveAll(ValidationManifestDir)
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

func DownloadRuntimeFiles(cmd *cobra.Command) (string, string, error) {
	var values version.NexusValues

	if err := version.GetNexusValues(&values); err != nil {
		return "", "", utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("could not download the runtime manifests due to %s", err)).Print().ExitIfFatalOrReturn()
	}
	runtimeVersion := os.Getenv("NEXUS_RUNTIME_TEMPLATE_VERSION")
	if runtimeVersion == "" {
		runtimeVersion = values.NexusDatamodelTemplates.Version
	}
	validationManifetsVersion := os.Getenv("NEXUS_VALIDATION_TEMPLATE_VERSION")
	if validationManifetsVersion == "" {
		validationManifetsVersion = values.NexusValidationTemplates.Version
	}
	err := utils.DownloadFile(fmt.Sprintf(common.RUNTIME_MANIFESTS_URL, runtimeVersion), runtimeFilename)

	if utils.IsDebug(cmd) {
		fmt.Printf("Using runtime manifests Version: %s\n", runtimeVersion)
	}
	if utils.IsDebug(cmd) {
		fmt.Printf("Using validation manifests Version: %s\n", validationManifetsVersion)
	}
	if err != nil {
		return "", "", utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("could not download the runtime manifests due to %s", err)).Print().ExitIfFatalOrReturn()
	}
	err = utils.DownloadFile(fmt.Sprintf(common.VALIDATION_MANIFESTS_URL, validationManifetsVersion), validationFilename)
	if err != nil {
		return "", "", utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("could not download the validation manifests due to %s", err)).Print().ExitIfFatalOrReturn()
	}
	runtimeFile, err := os.Open(runtimeFilename)
	if err != nil {
		return "", "", utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("accessing downloaded runtime manifests directory failed dwith error: %s", err)).Print().ExitIfFatalOrReturn()
	}
	defer runtimeFile.Close()
	validationFile, err := os.Open(validationFilename)
	if err != nil {
		return "", "", utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("accessing downloaded validation manifests directory failed dwith error: %s", err)).Print().ExitIfFatalOrReturn()
	}
	defer validationFile.Close()
	fo, err := os.Stat(RuntimeManifestsDir)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", "", utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
				fmt.Errorf("issues in checking runtime manifests directory due to %s", err)).Print().ExitIfFatalOrReturn()
		}
	}
	if fo != nil {
		err = os.RemoveAll(RuntimeManifestsDir)
		if err != nil {
			return "", "", utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
				fmt.Errorf("could not remove %s directory due to %s", RuntimeManifestsDir, err)).Print().ExitIfFatalOrReturn()
		}
	}
	// checking validation manifestDie
	fo, err = os.Stat(ValidationManifestDir)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", "", utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
				fmt.Errorf("issues in checking validation manifests directory due to %s", err)).Print().ExitIfFatalOrReturn()
		}
	}
	if fo != nil {
		err = os.RemoveAll(ValidationManifestDir)
		if err != nil {
			return "", "", utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
				fmt.Errorf("could not remove %s directory due to %s", ValidationManifestDir, err)).Print().ExitIfFatalOrReturn()
		}
	}
	err = os.Mkdir(RuntimeManifestsDir, os.ModePerm)
	if err != nil {
		return "", "", utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("could not create the runtime manifests directory due to %s", err)).Print().ExitIfFatalOrReturn()
	}
	err = os.Mkdir(ValidationManifestDir, os.ModePerm)
	if err != nil {
		return "", "", utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("could not create the runtime manifests directory due to %s", err)).Print().ExitIfFatalOrReturn()
	}
	err = os.Mkdir(filepath.Join(ValidationManifestDir, "manifests"), 0755)
	if err != nil {
		return "", "", utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("could not create the runtime manifests directory due to %s", err)).Print().ExitIfFatalOrReturn()
	}
	err = utils.Untar(RuntimeManifestsDir, runtimeFile)
	if err != nil {
		return "", "", utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("unarchive of runtime manifests directory failed with error %s", err)).Print().ExitIfFatalOrReturn()
	}
	err = utils.Untar(ValidationManifestDir, validationFile)
	if err != nil {
		return "", "", utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("unarchive of validation manifests directory failed with error %s", err)).Print().ExitIfFatalOrReturn()
	}

	return RuntimeManifestsDir, ValidationManifestDir, nil
}
