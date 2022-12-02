package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/log"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/config"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/upgrade"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/version"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/utils"
)

func RootPreRun(cmd *cobra.Command, args []string) error {
	nexusConfig, err := config.LoadNexusConfig()
	if err != nil {
		return err
	}

	if nexusConfig.DebugAlways {
		cmd.Flags().Lookup("debug").Changed = true
	}

	// Set logging level to debug if debugging is enabled.
	if utils.IsDebug(cmd) {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
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
					if err != nil {
						return err
					}

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
