package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
)

type K8sVersionObject struct {
	Version string `json:"gitVersion"`
}

type K8sVersionMap struct {
	ClientVersion K8sVersionObject
	ServerVersion K8sVersionObject
}

func GetNetworkingIngressVersion() (string, error) {
	versionStringBytes := exec.Command("kubectl", "version", "-o", "json")
	var out bytes.Buffer
	versionStringBytes.Stdout = &out
	err := versionStringBytes.Run()
	if err != nil {
		return "", fmt.Errorf("could not get version string")
	}
	versionObj := &K8sVersionMap{}
	err = json.Unmarshal(out.Bytes(), versionObj)
	if err != nil {
		return "", fmt.Errorf("Json unmarshal error in get version due to %s", err)
	}
	serverVersion := versionObj.ServerVersion.Version

	if len(serverVersion) == 0 {
		return "", fmt.Errorf("unable to get k8s version from output: %v", out.String())
	}

	v1min, _ := version.NewVersion("1.22.0")
	v1, errVersion := version.NewVersion(strings.TrimPrefix(serverVersion, "v"))
	if errVersion != nil {
		return "", fmt.Errorf("could not get network ingress version %s", errVersion)
	}
	if v1.LessThan(v1min) {
		return "v1beta1", nil
	} else {
		return "v1", nil
	}

}
func SystemCommand(cmd *cobra.Command, customErr ClientErrorCode, envList []string, name string, args ...string) error {
	// Get silent flag status.
	debugEnabled := IsDebug(cmd)

	if debugEnabled {
		if len(envList) != 0 {
			fmt.Printf("envList: %v\n", envList)
		}
		fmt.Printf("command: %v\n", name)
		fmt.Printf("args: %v\n", args)
	}
	command := exec.Command(name, args...)
	command.Env = os.Environ()

	if len(envList) > 0 {
		command.Env = append(command.Env, envList...)
	}

	stdout, err := command.StdoutPipe()
	if err != nil {
		return GetCustomError(INTERNAL_ERROR, err).Print().ExitIfFatalOrReturn()
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		return GetCustomError(INTERNAL_ERROR, err).Print().ExitIfFatalOrReturn()
	}
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			if debugEnabled {
				fmt.Printf("\t > %s\n", scanner.Text())
			}
		}
	}()

	errScanner := bufio.NewScanner(stderr)
	go func() {
		for errScanner.Scan() {
			if debugEnabled {
				fmt.Printf("\t > %s\n", errScanner.Text())
			}
		}
	}()
	err = command.Start()
	if err != nil {
		return GetCustomError(customErr,
			fmt.Errorf("starting cmd %s %v failed with error %v", name, args, err)).
			Print().ExitIfFatalOrReturn()
	}

	err = command.Wait()
	if err != nil {
		return GetCustomError(customErr,
			fmt.Errorf("waiting for cmd %s %v failed with error %v", name, args, err)).
			Print().ExitIfFatalOrReturn()
	}

	return nil
}
