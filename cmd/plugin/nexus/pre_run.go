package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/log"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/config"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/upgrade"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/version"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

func RootPreRun(cmd *cobra.Command, args []string) error {
	nexusConfig := config.LoadNexusConfig()

	if nexusConfig.DebugAlways {
		cmd.Flags().Lookup("debug").Changed = true
	}

	// Set logging level to debug if debugging is enabled.
	if utils.IsDebug(cmd) {
		log.SetLevel(log.DebugLevel)
	}

	if !nexusConfig.SkipUpgradeCheck {
		if isNewerVersionAvailable, latestVersion := version.IsNexusCliUpdateAvailable(); isNewerVersionAvailable {
			fmt.Printf("A new version of Nexus CLI (%s) is available\n", latestVersion)
			if !nexusConfig.UpgradePromptDisable {
				var input string
				fmt.Println("Would you like to upgrade? [y/n]")
				fmt.Scanln(&input)
				if input == "y" || input == "Y" || input == "yes" || input == "YES" {
					fmt.Println("Please specify the Nexus CLI version you'd like to upgrade to (press RETURN to upgrade to latest): ")
					n, err := fmt.Scanln(&input)
					if n == 0 {
						input = "latest"
					}
					err = upgrade.DoUpgradeCli(input, cmd)
					if err != nil {
						fmt.Printf("Failed to upgrade Nexus CLI. ")
						if !utils.IsDebug(cmd) {
							fmt.Println("Please retry upgrading by enabling the --debug flag and contact support with the output of the retry attached")
						} else {
							fmt.Println("Please contact support with the error output attached")
						}
						os.Exit(1)
					} else {
						fmt.Printf("Successfully upgraded Nexus CLI to %s\n", latestVersion)
						os.Exit(0)
					}
				}
			}
		}
	}
	return nil
}
