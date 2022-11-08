package utils

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	common "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	"gopkg.in/yaml.v2"
)

func GoToNexusDirectory() error {
	if _, err := os.Stat(common.NEXUS_DIR); os.IsNotExist(err) {
		return fmt.Errorf("%s directory not found", common.NEXUS_DIR)
	} else if err != nil {
		return fmt.Errorf("error %v trying to find directory %s", err, common.NEXUS_DIR)
	}

	if err := os.Chdir(common.NEXUS_DIR); err != nil {
		return fmt.Errorf("error %v trying to cd to directory %s", err, common.NEXUS_DIR)
	}
	return nil

}

func CheckDatamodelDirExists(datamodelName string) (bool, error) {
	dmDir := datamodelName
	if _, err := os.Stat(dmDir); os.IsNotExist(err) {
		return false, fmt.Errorf("datamodel directory %s not found", dmDir)
	} else if err != nil {
		return false, fmt.Errorf("error %v trying to find datamodel directory %s", err, dmDir)
	}
	return true, nil
}

func StoreCurrentDatamodel(datamodelName string) error {
	_, err := os.Stat(common.NexusConfFile)
	if err != nil {
		_, err = os.Create(common.NexusConfFile)
		if err != nil {
			return err
		}
	}

	conf := common.NexusConfig{
		Name: datamodelName,
	}
	yamlData, err := yaml.Marshal(&conf)
	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}
	err = ioutil.WriteFile(common.NexusConfFile, yamlData, 0644)
	if err != nil {
		fmt.Println("Could not store current datamodel name")
		return err
	}
	return nil
}

func GetCurrentDatamodel() (string, error) {
	_, err := os.Stat(common.NexusConfFile)
	if err != nil {
		return "", fmt.Errorf("Could not get datamodelname : %s does not exists\n", common.NexusConfFile)
	}
	data, err := ioutil.ReadFile(common.NexusConfFile)
	if err != nil {
		return "", fmt.Errorf("Could not read yamlfile ")
	}
	var yamlConfig common.NexusConfig
	err = yaml.Unmarshal(data, &yamlConfig)
	if err != nil {
		return "", fmt.Errorf("Could not unmarshal yamlfile")
	}
	return yamlConfig.Name, nil
}

func SetDatamodelDockerRepo(repoName string) error {
	_, err := os.Stat(common.NexusDMPropertiesFile)
	if err != nil {
		return fmt.Errorf("could not find nexus.yaml")
	}
	data, err := ioutil.ReadFile(common.NexusDMPropertiesFile)
	if err != nil {
		return fmt.Errorf("Could not read yamlfile ")
	}
	config := make(map[string]string)
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("Could not unmarshal yamlfile")
	}
	config["dockerRepo"] = repoName
	outData, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("Could not marshal yamlfile")
	}
	err = ioutil.WriteFile(common.NexusDMPropertiesFile, outData, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Could not write to yamlfile")
	}
	return nil

}

func RunDatamodelInstaller(DatamodelJobSpecConfig, Namespace, DatamodelName string, data common.Datamodel) error {
	fileContents, err := exec.Command("kubectl", "get", "cm", DatamodelJobSpecConfig, "-n", Namespace, "-ojsonpath={.data.jobSpec\\.yaml}").Output()
	if err != nil {
		return err
	}

	tmpl, err := template.New("template").Parse(strings.TrimLeft(string(fileContents), "'"))
	if err != nil {
		return err
	}
	JobName := fmt.Sprintf("%s-dmi", DatamodelName)
	k8sDeletecmd := exec.Command("kubectl", "delete", "-f", "-", "--ignore-not-found=true", "-n", Namespace)
	k8sApplycmd := exec.Command("kubectl", "apply", "-f", "-", "-n", Namespace)
	//data from renderedObject
	var renderedData bytes.Buffer
	err = tmpl.Execute(&renderedData, data)
	if err != nil {
		return err
	}

	renderedObject := renderedData.String()
	var DeleteStream bytes.Buffer
	var ApplyStream bytes.Buffer
	k8sDeletecmd.Stdin = &DeleteStream
	k8sApplycmd.Stdin = &ApplyStream
	go func() {
		io.WriteString(&DeleteStream, renderedObject)
		io.WriteString(&ApplyStream, renderedObject)
	}()

	var deleteOut bytes.Buffer
	k8sDeletecmd.Stderr = &deleteOut
	err = k8sDeletecmd.Run()
	if err != nil {
		return fmt.Errorf("could not delete the job objects due to %s", deleteOut.String())

	}
	maxRetries := 20
	current := 0
	for current < maxRetries {
		checkoutPut := bytes.Buffer{}
		k8schekcJobDeleted := exec.Command("kubectl", "get", "jobs", "-n", Namespace, JobName)
		k8schekcJobDeleted.Stderr = &checkoutPut
		err = k8schekcJobDeleted.Run()
		if err != nil {
			if strings.Contains(checkoutPut.String(), "not found") {
				break
			} else {
				return fmt.Errorf("could not check if existing job is deleted %s", JobName)
			}
		}
		time.Sleep(3 * time.Second)
		current += 1
	}

	var applyOut bytes.Buffer
	k8sApplycmd.Stderr = &applyOut
	err = k8sApplycmd.Run()
	if err != nil {
		return fmt.Errorf("could not create the job objects due to %s", applyOut.String())
	}

	jobCompletedCheck := exec.Command("kubectl", "wait", "job", JobName, "-n", Namespace, "--for=condition=complete", "--timeout=300s")
	err = jobCompletedCheck.Run()
	if err != nil {
		return fmt.Errorf("Datamodel installation job %s not be completed due to %s", JobName, err)
	}
	return nil
}
