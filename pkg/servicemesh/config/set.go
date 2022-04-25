package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

var disableUpgradePrompt bool

const flagDisableUpgradePrompt = "disable-upgrade-prompt"

func Set(cmd *cobra.Command, args []string) error {
	config := LoadNexusConfig()

	if cmd.Flags().Lookup(flagDisableUpgradePrompt).Changed {
		config.UpgradePromptDisable = disableUpgradePrompt
	}

	err := writeNexusConfig(config)
	if err != nil {
		return utils.GetCustomError(utils.CONFIG_SET_FAILED,
			fmt.Errorf("writing nexus config failed with error %v", err)).
			Print().ExitIfFatalOrReturn()
	}
	return nil
}

var SetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a Nexus config property",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: Set,
}

func init() {
	// each property gets a flag
	SetCmd.Flags().BoolVarP(&disableUpgradePrompt, flagDisableUpgradePrompt,
		"", false, "Enable/disable functionality in nexus CLI to check if a newer version is available")
}
