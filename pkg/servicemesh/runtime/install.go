package runtime

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/log"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/prereq"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

var Namespace string
var Registry string
var ImagePullSecret string
var IsNexusAdmin bool
var DryRunOutputFile string

type RuntimeInstallerData struct {
	RuntimeInstaller  common.RuntimeInstaller
	Namespace         string
	IsImagePullSecret bool
	ImagePullSecret   string
}

var installPrerequisites = []prereq.Prerequiste{
	prereq.KUBERNETES,
	prereq.KUBERNETES_VERSION,
}

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs the Nexus runtime on the specified namespace using helm",
	//Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if utils.ListPrereq(cmd) {
			return nil
		}

		if utils.SkipPrereqCheck(cmd) {
			return nil
		}

		if err := prereq.PreReqVerifyOnDemand(installPrerequisites); err != nil {
			return err
		}

		return nil
	},
	RunE: HelmInstall,
}

func CreateNs(Namespace string) error {
	createCmd := exec.Command("kubectl", "create", "namespace", Namespace, "--dry-run", "-oyaml")
	applyCmd := exec.Command("kubectl", "apply", "-f", "-")
	labelCmd := exec.Command("kubectl", "label", "namespace", Namespace, fmt.Sprintf("name=%s", Namespace))
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

	err = labelCmd.Start()
	if err != nil {
		return err
	}

	err = labelCmd.Wait()
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, &b2)
	return nil
}

func GetCustomTags(cmdlinArgs string) string {
	for _, value := range common.TagsList {
		if customTag := os.Getenv(value.VersionEnv); customTag != "" {
			cmdlinArgs = fmt.Sprintf("%s,global.%s.tag=%s", cmdlinArgs, value.VersionEnv, customTag)
		}
	}
	return cmdlinArgs
}

func HelmInstall(cmd *cobra.Command, args []string) error {
	Registry = strings.TrimSuffix(strings.TrimSpace(Registry), "/")
	cmdlineArgs := fmt.Sprintf("--set=global.namespace=%s", Namespace)
	cmdlineArgs = fmt.Sprintf("%s,global.registry=%s", cmdlineArgs, Registry)
	if ImagePullSecret != "" {
		cmdlineArgs = fmt.Sprintf("%s,global.imagepullsecret=%s", cmdlineArgs, ImagePullSecret)
	}
	if IsNexusAdmin {
		cmdlineArgs = fmt.Sprintf("%s,global.nexusAdmin=%t", cmdlineArgs, IsNexusAdmin)
	}
	cmdlineArgs = GetCustomTags(cmdlineArgs)
	runtimeVersion, err := utils.GetTagVersion("NexusRuntime", "NEXUS_RUNTIME_MANIFESTS_VERSION")
	if err != nil {
		return fmt.Errorf("could not get runtime version: %s", err)
	}
	checkNs := exec.Command("kubectl", "get", "ns", Namespace)
	err = checkNs.Run()
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
	var Args []string
	Args = []string{"upgrade", "--install", Namespace, "/chart.tgz", cmdlineArgs, "--wait", "--wait-for-jobs", "--timeout=10m"}

	var IsImagePullSecret bool = false
	if ImagePullSecret != "" {
		IsImagePullSecret = true
	}
	InstallerData := RuntimeInstallerData{
		RuntimeInstaller: common.RuntimeInstaller{
			Name:    fmt.Sprintf("%s-ins", Namespace),
			Image:   fmt.Sprintf("%s/nexus-runtime-chart:%s", Registry, runtimeVersion),
			Command: []string{"helm"},
			Args:    Args,
		},
		Namespace:         Namespace,
		IsImagePullSecret: IsImagePullSecret,
		ImagePullSecret:   ImagePullSecret,
	}

	yamlFile, err := common.RuntimeTemplate.ReadFile("runtime_installer.yaml")
	if err != nil {
		return fmt.Errorf("error while reading version yamlFile %v", err)

	}

	tmpl, err := template.New("template").Parse(strings.TrimLeft(string(yamlFile), "'"))
	if err != nil {
		return err
	}
	var applyString bytes.Buffer
	tmpl.Execute(&applyString, InstallerData)

	err = RunJob(Namespace, InstallerData.RuntimeInstaller.Name, applyString)
	if err != nil {
		return err
	}
	fmt.Printf("\u2713 Runtime installation successful on namespace %s\n", Namespace)

	return nil
}

func RunJob(Namespace, jobName string, applyString bytes.Buffer) error {
	var data []byte = applyString.Bytes()

	// this cmd is for ensuring the previous incomplete jobs can be deleted and re-ran again
	clearPreviouscmd := exec.Command("kubectl", "delete", "-f", "-", "-n", Namespace, "--ignore-not-found=true")
	clearPreviouscmd.Stdin = bytes.NewBuffer(data)
	clearPreviouscmd.Stdout = os.Stdout
	clearPreviouscmd.Stderr = os.Stderr

	if err := clearPreviouscmd.Start(); err != nil {
		return fmt.Errorf("Could not delete the existing installation job: %s on %s", jobName, Namespace)
	}

	if err := clearPreviouscmd.Wait(); err != nil {
		return fmt.Errorf("Could not delete the existing installation job: %s on %s", jobName, Namespace)
	}

	applyCmd := exec.Command("kubectl", "apply", "-f", "-", "-n", Namespace)
	applyCmd.Stdin = bytes.NewBuffer(data)
	applyCmd.Stdout = os.Stdout
	applyCmd.Stderr = os.Stderr

	if err := applyCmd.Start(); err != nil {
		return fmt.Errorf("Could not start the installation job: %s on %s", jobName, Namespace)
	}

	if err := applyCmd.Wait(); err != nil {
		return fmt.Errorf("Could not apply the installation job: %s on %s", jobName, Namespace)
	}

	err := exec.Command("kubectl", "wait", "--for=condition=complete", fmt.Sprintf("job/%s", jobName), "--timeout=10m", "-n", Namespace).Run()
	if err != nil {
		return fmt.Errorf("could not complete the installation job: %s on %s", jobName, Namespace)
	}

	return nil
}

func init() {
	InstallCmd.Flags().StringVarP(&Namespace, "namespace",
		"n", "", "name of the namespace to be created")
	InstallCmd.Flags().StringVarP(&Registry, "registry",
		"r", common.HarborRepo, "Registry where validation webhook and api-gw is located")
	InstallCmd.Flags().StringVarP(&ImagePullSecret, "secretname",
		"s", "", "Registry where validation webhook and api-gw is located")
	InstallCmd.Flags().BoolVarP(&IsNexusAdmin, "admin",
		"", false, "Install the Nexus Admin runtime")
	InstallCmd.Flags().StringVarP(&DryRunOutputFile, "output",
		"o", "", "Save genrated manifests to file")

	err := cobra.MarkFlagRequired(InstallCmd.Flags(), "namespace")
	if err != nil {
		log.Debugf("Runtime install err: %v", err)
	}

}
