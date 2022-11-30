package upgrade

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/version"
)

var (
	upgradeToVersion string
)

func upgradeCli(cmd *cobra.Command, args []string) error {
	var needUpgrade bool = false
	if upgradeToVersion == "" {
		needUpgrade, upgradeToVersion = version.IsNexusCliUpdateAvailable()
	} else {
		needUpgrade = true
	}
	if needUpgrade {
		return DoUpgradeCli(upgradeToVersion, cmd)
	} else {
		return nil
	}
}

func DoUpgradeCli(version string, cmd *cobra.Command) error {
	const nexusInstallScriptUrl = "https://raw.githubusercontent.com/vmware-tanzu/graph-framework-for-microservices/main/cli/get-nexus-cli.sh"
	if version == "" {
		version = "latest"
	}
	fmt.Printf("\u2713 Please upgrade nexus CLI version to %s with below steps\n", version)
	fmt.Printf("\u2794 curl -fsSL %s -o get-nexus-cli.sh \n\u2794 bash get-nexus-cli.sh --version %s\n", nexusInstallScriptUrl, version)
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
