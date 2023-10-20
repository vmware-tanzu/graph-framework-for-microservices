package runtime

import (
	"bytes"
	"fmt"

	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

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
var clientId string
var clientSecret string
var oAuthIssuerUrl string
var oAuthRedirectUrl string
var jwtClaim string
var jwtClaimValue string
var skipAdminBootstrap bool
var cpuResources *[]string
var memoryResources *[]string
var additionalOptions *[]string
var minimalRuntime bool

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

		if minimalRuntime == false {
			if err := prereq.PreReqVerifyOnDemand(installPrerequisites); err != nil {
				return err
			}
		}

		return nil
	},
	RunE: Install,
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
	_, err = io.Copy(os.Stdout, &b2)
	return err
}

func GetCustomTags(cmdlinArgs string) string {
	for _, value := range common.TagsList {
		if customTag := os.Getenv(value.VersionEnv); customTag != "" {
			cmdlinArgs = fmt.Sprintf("%s,global.%s.tag=%s", cmdlinArgs, value.FieldName, customTag)
		}
	}
	return cmdlinArgs
}

func HelmInstall(cmd *cobra.Command, args []string) error {
	Registry = strings.TrimSuffix(strings.TrimSpace(Registry), "/")
	cmdlineArgs := "--set="
	cmdlineArgs = fmt.Sprintf("%sglobal.namespace=%s", cmdlineArgs, Namespace)
	for resource, valueVariable := range common.Resources {
		apiVersion := utils.GetAPIGVK(resource)
		if apiVersion != "" {
			cmdlineArgs = fmt.Sprintf("%s,global.%s=%s", cmdlineArgs, valueVariable, apiVersion)
		}
	}
	for _, value := range *cpuResources {
		cmdlineArgs = fmt.Sprintf("%s,global.resources.%s.cpu=%s", cmdlineArgs, strings.Split(value, "=")[0], strings.Split(value, "=")[1])
	}

	for _, value := range *memoryResources {
		cmdlineArgs = fmt.Sprintf("%s,global.resources.%s.memory=%s", cmdlineArgs, strings.Split(value, "=")[0], strings.Split(value, "=")[1])
	}
	for _, value := range *additionalOptions {
		cmdlineArgs = fmt.Sprintf("%s,global.%s=%s", cmdlineArgs, strings.Split(value, "=")[0], strings.Split(value, "=")[1])
	}

	cmdlineArgs = fmt.Sprintf("%s,global.registry=%s", cmdlineArgs, Registry)
	if ImagePullSecret != "" {
		cmdlineArgs = fmt.Sprintf("%s,global.imagepullsecret=%s", cmdlineArgs, ImagePullSecret)
	}
	if IsNexusAdmin {
		cmdlineArgs = fmt.Sprintf("%s,global.nexusAdmin=%t", cmdlineArgs, IsNexusAdmin)
		cmdlineArgs = fmt.Sprintf("%s,global.skipAdminBootstrap=%t", cmdlineArgs, skipAdminBootstrap)
		if !skipAdminBootstrap {
			if clientId == "" || clientSecret == "" || oAuthIssuerUrl == "" || oAuthRedirectUrl == "" || jwtClaim == "" || jwtClaimValue == "" {
				return fmt.Errorf("at least one mandatory arg (client-id, client-secret, oauth-issuer-url, oauth-redirect-url, jwt-clain, jwt-claim-value) missing in admin runtime install")
			} else {
				cmdlineArgs = fmt.Sprintf("%s,global.clientId=%s", cmdlineArgs, clientId)
				cmdlineArgs = fmt.Sprintf("%s,global.clientSecret=%s", cmdlineArgs, clientSecret)
				cmdlineArgs = fmt.Sprintf("%s,global.oAuthIssuerUrl=%s", cmdlineArgs, oAuthIssuerUrl)
				cmdlineArgs = fmt.Sprintf("%s,global.oAuthRedirectUrl=%s", cmdlineArgs, oAuthRedirectUrl)
				cmdlineArgs = fmt.Sprintf("%s,global.jwtClaim=%s", cmdlineArgs, jwtClaim)
				cmdlineArgs = fmt.Sprintf("%s,global.jwtClaimValue=%s", cmdlineArgs, jwtClaimValue)
			}
		}
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
	Args = []string{"upgrade", "--install", Namespace, "/chart.tgz", cmdlineArgs, "--wait", "--wait-for-jobs", "--timeout=15m"}

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
	err = tmpl.Execute(&applyString, InstallerData)
	if err != nil {
		return err
	}

	fmt.Printf("Install job starting at %s\n", time.Now())
	err = RunJob(Namespace, InstallerData.RuntimeInstaller.Name, applyString)
	fmt.Printf("Install job ended at %s\n", time.Now())
	if err != nil {
		return err
	}
	fmt.Printf("\u2713 Runtime installation successful on namespace %s\n", Namespace)

	return nil
}

func minimalInstall(cmd *cobra.Command, args []string) error {

	res, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return fmt.Errorf("unable to determine if the CWD is a git repo. Error: %v", err)
	}

	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return fmt.Errorf("unable to get CWD to dertermine if CWD is a git repo. Error: %v", err)
	}

	if string(res[:len(res)-1]) != cwd {
		return fmt.Errorf("Current directory %s is not the root directory of nexus repo. Retry the command after cd to the nexus repo root dir.\n", cwd)
	}

	runtimeInstallCmd := exec.Command("make", "install.runtime.k0s")
	runtimeInstallCmd.Stdout = os.Stdout
	runtimeInstallCmd.Stderr = os.Stderr
	err = runtimeInstallCmd.Run()
	if err != nil {
		return fmt.Errorf("minimal runtime install failed with error: %v", err)
	}

	return nil
}

func RunJob(Namespace, jobName string, applyString bytes.Buffer) error {
	var data []byte = applyString.Bytes()

	err := DeleteJob(data, jobName, Namespace)
	if err != nil {
		return err
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

	err = exec.Command("kubectl", "wait", "--for=condition=complete", fmt.Sprintf("job/%s", jobName), "--timeout=15m", "-n", Namespace).Run()
	if err != nil {
		// dump some info to debug
		log.Debugf("Pod Status")
		if pods, err := exec.Command("kubectl", "get", "pods", "-n", Namespace).Output(); err == nil {
			log.Debugf(string(pods))
		}

		log.Debugf("Job Status")
		if jobs, err := exec.Command("kubectl", "get", "jobs", "-n", Namespace).Output(); err == nil {
			log.Debugf(string(jobs))
		}
		return fmt.Errorf("could not complete the installation job: %s on %s", jobName, Namespace)
	}

	err = DeleteJob(data, jobName, Namespace)
	if err != nil {
		return err
	}
	return nil
}

func DeleteJob(data []byte, jobName, Namespace string) error {
	deleteCmd := exec.Command("kubectl", "delete", "-f", "-", "-n", Namespace, "--ignore-not-found=true")
	deleteCmd.Stdin = bytes.NewBuffer(data)
	deleteCmd.Stdout = os.Stdout
	deleteCmd.Stderr = os.Stderr

	// add delete job after installation or uninstallation completed
	if err := deleteCmd.Start(); err != nil {
		return fmt.Errorf("Could not delete the existing installation job objects: %s on %s", jobName, Namespace)
	}

	if err := deleteCmd.Wait(); err != nil {
		return fmt.Errorf("Could not delete the existing installation job objects: %s on %s", jobName, Namespace)
	}
	return nil
}

func Install(cmd *cobra.Command, args []string) error {
	if minimalRuntime {
		return minimalInstall(cmd, args)
	}
	return HelmInstall(cmd, args)
}

func init() {
	InstallCmd.Flags().StringVarP(&Namespace, "namespace",
		"n", "", "name of the namespace to be created")
	InstallCmd.Flags().StringVarP(&Registry, "registry",
		"r", common.ImageRegistry, "Registry where validation webhook and api-gw is located")
	InstallCmd.Flags().StringVarP(&ImagePullSecret, "secretname",
		"s", "", "Registry where validation webhook and api-gw is located")
	InstallCmd.Flags().BoolVarP(&IsNexusAdmin, "admin",
		"", false, "Install the Nexus Admin runtime")
	InstallCmd.Flags().StringVarP(&clientId, "client-id",
		"", "", "client id of the OIDC application. ignored if not --admin runtime")
	InstallCmd.Flags().StringVarP(&clientSecret, "client-secret",
		"", "", "client secret of the OIDC application. ignored if not --admin runtime")
	InstallCmd.Flags().StringVarP(&oAuthIssuerUrl, "oauth-issuer-url",
		"", "", "OAuth Issuer URL of the identity provider. ignored if not --admin runtime")
	InstallCmd.Flags().StringVarP(&oAuthRedirectUrl, "oauth-redirect-url",
		"", "", "OAuth Redirect/Callback URL. ignored if not --admin runtime")
	InstallCmd.Flags().StringVarP(&jwtClaim, "jwt-claim",
		"", "", "the JWT claim to be used as part of the admin match condition. ignored if not --admin runtime")
	InstallCmd.Flags().StringVarP(&jwtClaimValue, "jwt-claim-value",
		"", "", "the JWT claim to be used as part of the admin match condition. ignored if not --admin runtime")
	InstallCmd.Flags().BoolVarP(&skipAdminBootstrap, "skip-bootstrap",
		"", false, "skips the bootstrap step (only relevant for admin-runtime)")
	cpuResources = InstallCmd.Flags().StringArrayP("cpuResources", "",
		[]string{}, "for configuring cpu resources")
	memoryResources = InstallCmd.Flags().StringArrayP("memoryResources", "",
		[]string{}, "for configuring memory resources")
	additionalOptions = InstallCmd.Flags().StringArrayP("options", "",
		[]string{}, "for configuring additional helm values")
	InstallCmd.Flags().BoolVarP(&minimalRuntime, "minimal",
		"", false, "Install a minimalistic runtime. Needs a git clone of source code repo")

	err := cobra.MarkFlagRequired(InstallCmd.Flags(), "namespace")
	if err != nil {
		log.Debugf("Runtime install err: %v", err)
	}

}
