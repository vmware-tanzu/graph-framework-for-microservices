package upgrade

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/common"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/version"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/utils"
)

var (
	upgradeToVersion string
	nexusCliRepo     string = "github.com/vmware-tanzu/graph-framework-for-microservices/cli/cmd/plugin/nexus"
)

func upgradeCli(cmd *cobra.Command, args []string) error {
	var needUpgrade bool = false
	if upgradeToVersion == "" {
		needUpgrade, upgradeToVersion = version.IsNexusCliUpdateAvailable()
	} else {
		needUpgrade = true
	}
	if needUpgrade {
		fmt.Printf("Upgrading to version: %s\n", upgradeToVersion)
		return DoUpgradeCli(upgradeToVersion, cmd)
	} else {
		return nil
	}
}

func DoUpgradeCli(version string, cmd *cobra.Command) error {
	if version == "" {
		version = "master"
	}
	envList := common.GetEnvList()
	err := utils.SystemCommand(cmd, utils.CLI_UPGRADE_FAILED, envList, "go", "install", fmt.Sprintf("%s@%s", nexusCliRepo, version))
	if err == nil {
		fmt.Printf("\u2713 CLI successfully upgraded to version %s\n", version)
	} else {
		fmt.Printf("\u274C CLI upgrade to version %s failed with error %v\n", version, err)
	}
	return nil
}

var UpgradeCliCmd = &cobra.Command{
	Use:   "cli",
	Short: "upgrade cli",
	RunE:  upgradeCli,
}

func init() {
	UpgradeCliCmd.Flags().StringVarP(&upgradeToVersion, "version",
		"v", "", "version to upgrade to")
}
