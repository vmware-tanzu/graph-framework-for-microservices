package runtime

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var RunApiCmd = &cobra.Command{
	Use:   "start-api",
	Short: "start an api gateway serving api installed in the runtime",
	RunE:  RunApi,
}

var runtimekubeconfig string

func RunApi(cmd *cobra.Command, args []string) error {

	res, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return fmt.Errorf("unable to determine if the CWD is the nexus git repo. Error: %v", err)
	}

	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return fmt.Errorf("unable to get CWD to dertermine if CWD is the nexut git repo. Error: %v", err)
	}

	if string(res[:len(res)-1]) != cwd {
		return fmt.Errorf("Current directory %s is not the root directory of nexus repo. Retry the command after cd to the nexus repo root dir.\n", cwd)
	}

	runApiGWCmd := exec.Command("make", "api-gw.run")
	runApiGWCmd.Env = os.Environ()
	runApiGWCmd.Env = append(runApiGWCmd.Env, fmt.Sprintf("HOST_KUBECONFIG=%s", runtimekubeconfig))
	runApiGWCmd.Stdout = os.Stdout
	runApiGWCmd.Stderr = os.Stderr
	err = runApiGWCmd.Run()
	if err != nil {
		return fmt.Errorf("running api gateway failed with error: %v", err)
	}
	return nil
}

func init() {
	RunApiCmd.Flags().StringVarP(&runtimekubeconfig, "kubeconfig", "", "", "kubeconfig file")
}
