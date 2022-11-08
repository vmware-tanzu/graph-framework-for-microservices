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
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/log"
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
	if len(envList) != 0 {
		log.Debugf("envList: %v\n", envList)
	}
	log.Debugf("command: %v\n", name)
	log.Debugf("args: %v\n", args)

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
			text := scanner.Text()
			log.Debugf("\t > %s\n", text)
			if !IsDebug(cmd) && strings.Contains(text, "Nexus Compiler") {
				log.Infof("%s", text)
			}
		}
	}()

	errScanner := bufio.NewScanner(stderr)
	go func() {
		for errScanner.Scan() {
			log.Debugf("\t > %s\n", errScanner.Text())
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

func WriteToFile(filename string, data []byte) error {
	var OutputFileObject *os.File
	_, err := os.Stat(filename)
	if err != nil {
		OutputFileObject, err = os.Create(filename)
		if err != nil {
			return fmt.Errorf("Could not create output file: %s due to %s", filename, err)
		}
		OutputFileObject.Close()
	}
	OutputFileObject, err = os.OpenFile(filename, os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Could not open output file: %s due to %s", filename, err)
	}
	defer OutputFileObject.Close()
	_, err = OutputFileObject.Write(data)
	if err != nil {
		return fmt.Errorf("err: %s", err)
	}
	return nil
}

func AddHelmRepo(RepoName, RepoUrl string) error {
	_ = exec.Command("helm", "repo", "remove", RepoName).Run()
	CreateOutput, err := exec.Command("helm", "repo", "add", RepoName, RepoUrl).CombinedOutput()
	if err != nil {
		return fmt.Errorf("could not add %s to helm repos due to %s", RepoUrl, CreateOutput)
	}
	return nil
}
