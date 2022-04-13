package version

import (
	"fmt"

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

	yamlFile, err := common.TemplateFs.ReadFile("values.yaml")
	if err != nil {
		return fmt.Errorf("error while reading version yamlFile %v", err)

	}

	err = yaml.Unmarshal(yamlFile, &values)
	if err != nil {
		return fmt.Errorf("error while unmarshal version yaml data %v", err)
	}

	fmt.Printf("NexusCli: %s\n", values.NexusCli.Version)
	fmt.Printf("NexusCompiler: %s\n", values.NexusCompiler.Version)
	fmt.Printf("NexusAppTemplates: %s\n", values.NexusAppTemplates.Version)
	return nil
}
