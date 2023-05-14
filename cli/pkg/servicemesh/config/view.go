package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func View(cmd *cobra.Command, args []string) error {
	nexusConfig, err := LoadNexusConfig()
	if err != nil {
		return err
	}

	data, _ := yaml.Marshal(&nexusConfig)
	fmt.Println(string(data))
	return nil
}

var ViewCmd = &cobra.Command{
	Use:   "view",
	Short: "Displays the current Nexus config",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: View,
}
