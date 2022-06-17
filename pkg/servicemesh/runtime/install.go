package runtime

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/prereq"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/version"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/logging"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/rand"
)

var Namespace string
var Registry string
var ImagePullSecret string
var NetworkingAPIVersion string
var IsNexusAdmin bool

var prerequisites = []prereq.Prerequiste{
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

		if err := prereq.PreReqVerifyOnDemand(prerequisites); err != nil {
			return err
		}
		NetworkingAPIVersion, err = utils.GetNetworkingIngressVersion()
		if err != nil {
			return err
		}
		return nil
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

func DownloadManifestsFile(manifest common.Manifest, values version.NexusValues) (string, error) {
	versionTo := os.Getenv(manifest.VersionEnv)
	if versionTo == "" {
		versionTo = reflect.ValueOf(values).FieldByName(manifest.VersionStrName).Field(0).String()
	}
	fmt.Printf("Version of %s is %s\n", manifest.FileName, versionTo)
	err := utils.DownloadFile(fmt.Sprintf(manifest.URL, versionTo), manifest.FileName)
	if err != nil {
		return "", err
	}
	fileObj, err := os.Open(manifest.FileName)
	if err != nil {
		return "", utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("accessing downloaded runtime manifests directory failed dwith error: %s", err)).Print().ExitIfFatalOrReturn()
	}
	defer fileObj.Close()

	fo, err := os.Stat(manifest.Directory)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
				fmt.Errorf("issues in checking runtime manifests directory due to %s", err)).Print().ExitIfFatalOrReturn()
		}
	}
	if fo != nil {
		err = os.RemoveAll(manifest.Directory)
		if err != nil {
			return "", utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
				fmt.Errorf("could not remove %s directory due to %s", manifest.Directory, err)).Print().ExitIfFatalOrReturn()
		}
	}
	err = os.Mkdir(manifest.Directory, os.ModePerm)
	if err != nil {
		return "", utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("could not create the runtime manifests directory due to %s", err)).Print().ExitIfFatalOrReturn()
	}
	err = utils.Untar(manifest.Directory, fileObj)
	if err != nil {
		return "", utils.GetCustomError(utils.RUNTIME_INSTALL_FAILED,
			fmt.Errorf("unarchive of runtime manifests directory failed with error %s", err)).Print().ExitIfFatalOrReturn()
	}
	if manifest.Templatized {
		manifest.Image.Tag = versionTo
		err = utils.RenderTemplateFiles(manifest.Image, manifest.Directory, ".git")
		if err != nil {
			return "", err
		}
	}
	return manifest.Directory, nil
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

	var IsImagePullSecret bool
	if Registry == "" {
		Registry = common.HarborRepo
	}
	Registry = strings.TrimSpace(Registry)
	Registry = strings.TrimSuffix(Registry, "/")
	if ImagePullSecret == "" {
		IsImagePullSecret = false
	} else {
		IsImagePullSecret = true
	}
	if utils.IsDebug(cmd) {
		fmt.Printf("Using Images from %s registry for api-gateway and nexus-validation\n", Registry)
	}
	for index, manifest := range common.RuntimeManifests {
		if manifest.Templatized {
			manifest.Image = common.ImageTemplate{
				Image:                fmt.Sprintf("%s/%s", Registry, manifest.ImageName),
				IsImagePullSecret:    IsImagePullSecret,
				ImagePullSecret:      ImagePullSecret,
				Namespace:            Namespace,
				NetworkingAPIVersion: NetworkingAPIVersion,
			}
		}
		common.RuntimeManifests[index] = manifest
	}
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

	// add a nexus label to differentiate this namespace from others
	if IsNexusAdmin {
		_, err = exec.Command("kubectl", "label", "ns", Namespace, "nexus=admin", "--overwrite").Output()
		if err != nil {
			return fmt.Errorf("failed to label namespace %s: %s", Namespace, err.Error())
		}
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
		fmt.Printf("Applying file: %s\n", file)
		err = utils.SystemCommand(cmd, utils.RUNTIME_INSTALL_FAILED, []string{}, "kubectl", "apply", "-f", file, "-n", Namespace)
		if err != nil {
			return err
		}
	}
	fmt.Println("Waiting for the Nexus runtime to come up...")
	for _, label := range common.RuntimePodLabels {
		utils.CheckPodRunning(cmd, utils.RUNTIME_INSTALL_FAILED, label, Namespace)
	}
	for _, manifest := range common.RuntimeManifests {
		os.RemoveAll(manifest.Directory)
		os.RemoveAll(manifest.FileName)
	}

	fmt.Println("Installing API Datamodel CRDs...")
	err = installApiDatamodel(cmd, versions)
	if err != nil {
		utils.GetCustomError(utils.RUNTIME_INSTALL_API_DATAMODEL_INSTALL_FAILED,
			fmt.Errorf("installing API datamodel on nexus-apiserver failed: %s", err)).Print().ExitIfFatalOrReturn()
	}
	for _, label := range common.ApiDmDependentPodLabels {
		utils.CheckPodRunning(cmd, utils.RUNTIME_INSTALL_FAILED, label, Namespace)
	}

	fmt.Printf("\u2713 Runtime installation successful on namespace %s\n", Namespace)
	return nil
}

func installApiDatamodel(cmd *cobra.Command, values version.NexusValues) error {
	dir, err := DownloadManifestsFile(common.NexusApiDatamodelManifest, values)
	if err != nil {
		fmt.Println("Error downloading API datamodel manifest")
		return err
	}

	// get a nexus-proxy-container pod
	labels := "app=nexus-proxy-container"
	nexusProxyContainerPod := utils.GetPodByLabelAndState(Namespace, labels, v1.PodRunning)
	if nexusProxyContainerPod == nil {
		return fmt.Errorf("no running pod with label %s found", labels)
	}

	// channels to manage lifecycle of the port-forward session
	var stopCh = make(chan struct{}, 1)
	var readyCh = make(chan struct{})

	// managing termination signal from the terminal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-stopCh:
			fmt.Println("Received a stop signal, closing stopCh")
			close(stopCh)
		case <-sigs:
			fmt.Println("Received a process termination signal, closing stopCh")
			close(stopCh)
		}
	}()

	// prepare for port-forward
	localPort := rand.IntnRange(40000, 45000)
	nexusProxyContainerPort := 8001
	go func() {
		err = utils.StartPortforward(nexusProxyContainerPod.Name, Namespace, localPort, nexusProxyContainerPort, stopCh, readyCh, os.Stdout, os.Stdout)
		if err != nil {
			fmt.Printf("error initiating port-forward to %s: %s\n", nexusProxyContainerPod.Name, err.Error())
			os.Exit(1)
		}
	}()
	fmt.Printf("Started port-forward to pod %s on port %d\n", nexusProxyContainerPod.Name, localPort)

	// give the port-forward a few secs
	err = utils.CheckLocalAPIServer(fmt.Sprintf("%s:%d", "localhost", localPort), 30, 2*time.Second)
	if err != nil {
		return fmt.Errorf("checking local apiserver started failed ")
	}

	// apply CRDs
	err = utils.SystemCommand(cmd, utils.RUNTIME_INSTALL_FAILED, []string{}, "kubectl", "apply", "-f", dir, "--recursive", "-s", fmt.Sprintf("%s:%d", "localhost", localPort))
	if err != nil {
		return fmt.Errorf("applying API datamodel failed due to error: %s", err.Error())
	}

	// stop port-forward
	stopCh <- struct{}{}

	// cleanup
	os.RemoveAll(dir)
	os.RemoveAll(common.NexusApiDatamodelManifest.FileName)

	return nil
}

func init() {
	InstallCmd.Flags().StringVarP(&Namespace, "namespace",
		"n", "", "name of the namespace to be created")
	InstallCmd.Flags().StringVarP(&Registry, "registry",
		"r", "", "Registry where validation webhook and api-gw is located")
	InstallCmd.Flags().StringVarP(&ImagePullSecret, "secretname",
		"s", "", "Registry where validation webhook and api-gw is located")
	InstallCmd.Flags().BoolVarP(&IsNexusAdmin, "admin",
		"", false, "Install the Nexus Admin runtime")

	err := cobra.MarkFlagRequired(InstallCmd.Flags(), "namespace")
	if err != nil {
		logging.Debugf("Runtime install err: %v", err)
	}

}

func DownloadRuntimeFiles(cmd *cobra.Command, versions version.NexusValues) ([]string, error) {
	var files []string
	for _, manifest := range common.RuntimeManifests {
		file, err := DownloadManifestsFile(manifest, versions)
		if err != nil {
			return files, err
		}
		files = append(files, file)
	}
	return files, nil

}
