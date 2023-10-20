package datamodel

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var kubeconfig string

func InstallFromDirectory(cmd *cobra.Command, args []string) error {

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

	dmDirectory := args[0]
	if dmDirectory == "" {
		return fmt.Errorf("a valid datamodel directory is mandatory")
	}

	dmInstallCmd := exec.Command("make", "dm.install")
	dmInstallCmd.Env = os.Environ()
	dmInstallCmd.Env = append(dmInstallCmd.Env, fmt.Sprintf("DATAMODEL_DIR=%s", dmDirectory))
	dmInstallCmd.Env = append(dmInstallCmd.Env, fmt.Sprintf("HOST_KUBECONFIG=%s", kubeconfig))
	dmInstallCmd.Stdout = os.Stdout
	dmInstallCmd.Stderr = os.Stderr
	err = dmInstallCmd.Run()
	if err != nil {
		return fmt.Errorf("datamodel installed from directory %s failed with error: %v", dmDirectory, err)
	}
	return nil
}

var DirCmd = &cobra.Command{
	Use:   "dir",
	Short: "datamodel directory with built artifacts",
	Args:  cobra.MinimumNArgs(1),
	RunE:  InstallFromDirectory,
}

func init() {
	DirCmd.Flags().StringVarP(&kubeconfig, "kubeconfig", "", "", "kubeconfig file")
}
