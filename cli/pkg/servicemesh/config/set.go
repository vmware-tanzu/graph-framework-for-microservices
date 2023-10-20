package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

var disableUpgradePrompt bool
var debugAlways bool
var skipUpgradeCheck bool

const flagDisableUpgradePrompt = "disable-upgrade-prompt"
const flagDebugAlways = "debug-always"
const flagSkipUpgradeCheck = "skip-upgrade-check"

func noFlagsSet(cmd *cobra.Command) bool {
	return !cmd.Flags().Lookup(flagDisableUpgradePrompt).Changed &&
		!cmd.Flags().Lookup(flagDebugAlways).Changed &&
		!cmd.Flags().Lookup(flagSkipUpgradeCheck).Changed
}

func Set(cmd *cobra.Command, args []string) error {
	if noFlagsSet(cmd) {
		return utils.GetCustomError(utils.CONFIG_SET_FAILED,
			fmt.Errorf("`nexus config set` expects at least one flag to be set")).
			Print().ExitIfFatalOrReturn()
	}
	config, err := LoadNexusConfig()
	if err != nil {
		return err
	}

	if cmd.Flags().Lookup(flagDisableUpgradePrompt).Changed {
		config.UpgradePromptDisable = disableUpgradePrompt
	}

	if cmd.Flags().Lookup(flagDebugAlways).Changed {
		config.DebugAlways = debugAlways
	}

	if cmd.Flags().Lookup(flagSkipUpgradeCheck).Changed {
		config.SkipUpgradeCheck = skipUpgradeCheck
	}

	err = writeNexusConfig(config)
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
	SetCmd.Flags().BoolVarP(&debugAlways, flagDebugAlways,
		"", false, "Print debug output even without specifying the --debug flag")
	SetCmd.Flags().BoolVarP(&skipUpgradeCheck, flagSkipUpgradeCheck,
		"", false, "Skip checking for the latest available CLI version")
}
