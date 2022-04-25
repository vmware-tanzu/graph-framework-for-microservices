package version

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	"gopkg.in/yaml.v2"
)

type NexusValues struct {
	NexusCli          versionFields `yaml:"nexusCli"`
	NexusCompiler     versionFields `yaml:"nexusCompiler"`
	NexusAppTemplates versionFields `yaml:"nexusAppTemplates"`
}

type versionFields struct {
	Version string `yaml:"version"`
}

func Version(cmd *cobra.Command, args []string) error {
	var values NexusValues

	if err := GetNexusValues(&values); err != nil {
		return err
	}

	fmt.Printf("NexusCli: %s\n", values.NexusCli.Version)
	fmt.Printf("NexusCompiler: %s\n", values.NexusCompiler.Version)
	fmt.Printf("NexusAppTemplates: %s\n", values.NexusAppTemplates.Version)
	return nil
}

func GetNexusValues(values *NexusValues) error {
	yamlFile, err := common.TemplateFs.ReadFile("values.yaml")
	if err != nil {
		return fmt.Errorf("error while reading version yamlFile %v", err)

	}

	err = yaml.Unmarshal(yamlFile, &values)
	if err != nil {
		return fmt.Errorf("error while unmarshal version yaml data %v", err)
	}

	return nil
}

func GetLatestNexusVersion() (string, error) {
	const cliRepo = "git@gitlab.eng.vmware.com:nsx-allspark_users/nexus-sdk/cli"
	output, err := exec.Command("git", "ls-remote", "-t", "--sort", "-v:refname", cliRepo).Output()
	if err != nil {
		errMsg := fmt.Sprintf("Nexus CLI Upgrade check: failed to fetch remote tags from Nexus CLI repo. Please ensure you are able to clone this repo: `git clone %s`", cliRepo)
		fmt.Println(errMsg)
		return "", fmt.Errorf(errMsg)
	}
	if len(output) == 0 {
		return "", fmt.Errorf("No tags found")
	}
	strOutput := strings.Split(string(output), "\n")
	if len(strOutput) == 0 {
		return "", fmt.Errorf("No tags found")
	}
	line := strOutput[0] // because of the sort order (descending), the first one would be the latest
	// an example of what 'line' looks like
	// e2e3bf7de9fcda76d0d1f647fcb92a9d9451b11d	refs/tags/v7.3.7
	lineParts := strings.Fields(line)
	if len(lineParts) != 2 {
		return "", fmt.Errorf("output format different from expected")
	}
	// we're interested in just the second part of 'line'
	tagString := lineParts[len(lineParts)-1]
	tagsRegex := regexp.MustCompile(`refs/tags/([a-z0-9][^\s]*)`)
	versionString := tagsRegex.FindStringSubmatch(tagString)
	if len(versionString) != 2 {
		return "", fmt.Errorf("version string output format different from expected")
	} else {
		return versionString[1], nil
	}
}

func format(versionString string) string {
	if strings.HasPrefix(versionString, "v") {
		versionString = strings.TrimPrefix(versionString, "v")
	}
	if strings.HasSuffix(versionString, "\n") {
		versionString = strings.TrimSuffix(versionString, "\n")
	}
	return versionString
}
