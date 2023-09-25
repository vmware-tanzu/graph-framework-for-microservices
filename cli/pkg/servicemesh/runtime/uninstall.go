package runtime

import (
	"bytes"
	"fmt"
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

var minimalRuntimeUninstall bool

func minimalUninstall(cmd *cobra.Command, args []string) error {

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

	runtimeUninstallCmd := exec.Command("make", "uninstall.runtime.k0s")
	runtimeUninstallCmd.Stdout = os.Stdout
	runtimeUninstallCmd.Stderr = os.Stderr
	err = runtimeUninstallCmd.Run()
	if err != nil {
		return fmt.Errorf("minimal runtime uninstall failed with error: %v", err)
	}

	return nil
}

func Uninstall(cmd *cobra.Command, args []string) error {
	if minimalRuntimeUninstall {
		return minimalUninstall(cmd, args)
	}
	return UninstallHelm(cmd, args)
}
func UninstallHelm(cmd *cobra.Command, args []string) error {
	cmdlineArgs := fmt.Sprintf("--set=global.namespace=%s", Namespace)
	for resource, valueVariable := range common.Resources {
		apiVersion := utils.GetAPIGVK(resource)
		if apiVersion != "" {
			cmdlineArgs = fmt.Sprintf("%s,global.%s=%s", cmdlineArgs, valueVariable, apiVersion)
		}
	}
	cmdlineArgs = fmt.Sprintf("%s,global.registry=%s", cmdlineArgs, Registry)
	runtimeVersion, err := utils.GetTagVersion("NexusRuntime", "NEXUS_RUNTIME_MANIFESTS_VERSION")
	if err != nil {
		return fmt.Errorf("could not get runtime version: %s", err)
	}
	var IsImagePullSecret bool = false
	if ImagePullSecret != "" {
		IsImagePullSecret = true
	}
	checkNs := exec.Command("kubectl", "get", "ns", Namespace)
	err = checkNs.Run()
	if err != nil {
		log.Infof("Namespace %s is not available", Namespace)
		return nil
	}
	InstallerData := RuntimeInstallerData{
		RuntimeInstaller: common.RuntimeInstaller{
			Name:    fmt.Sprintf("%s-unins", Namespace),
			Image:   fmt.Sprintf("%s/nexus-runtime-chart:%s", Registry, runtimeVersion),
			Command: []string{"/bin/bash"},
			Args:    []string{"-c", fmt.Sprintf("helm template /chart.tgz %s | kubectl delete  -f - --ignore-not-found=true", cmdlineArgs)},
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

	err = RunJob(Namespace, InstallerData.RuntimeInstaller.Name, applyString)
	if err != nil {
		return err
	}

	err = utils.SystemCommand(cmd, utils.RUNTIME_UNINSTALL_FAILED, []string{}, "kubectl", "delete", "pvc", "-n", Namespace, "-lapp=nexus-etcd")
	if err != nil {
		return err
	}

	err = utils.SystemCommand(cmd, utils.RUNTIME_UNINSTALL_FAILED, []string{}, "kubectl", "delete", "pvc", "-n", Namespace, "-lcreated-by=nexus")
	if err != nil {
		return err
	}

	//adding additional job cleanup after uninstallation...
	err = utils.SystemCommand(cmd, utils.RUNTIME_UNINSTALL_FAILED, []string{}, "kubectl", "delete", "jobs", "-n", Namespace, "--all", "--ignore-not-found=true")
	if err != nil {
		return err
	}

	fmt.Printf("\u2713 Runtime %s uninstallation successful\n", Namespace)
	return nil
}

var UninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstalls the Nexus runtime from the specified namespace",
	//Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		err = prereq.PreReqVerifyOnDemand(installPrerequisites)
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
	UninstallCmd.Flags().StringVarP(&Registry, "registry",
		"r", common.ImageRegistry, "Registry where helm-chart is located")
	err := cobra.MarkFlagRequired(UninstallCmd.Flags(), "namespace")
	if err != nil {
		log.Debugf("Runtime uninstall err: %v", err)
	}
	UninstallCmd.Flags().BoolVarP(&minimalRuntimeUninstall, "minimal",
		"", false, "Uninstall a minimalistic runtime. Needs a git clone of source code repo")
}
