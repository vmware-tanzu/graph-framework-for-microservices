package upgrade

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

var (
	upgradeToVersion string
	nexusCliRepo     string = "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/cmd/plugin/nexus"
)

func UpgradeCli(cmd *cobra.Command, args []string) error {

	if upgradeToVersion == "" {
		upgradeToVersion = "latest"
	}
	envList := common.GetEnvList()
	err := utils.SystemCommand(cmd, utils.CLI_UPGRADE_FAILED, envList, "go", "install", fmt.Sprintf("%s@%s", nexusCliRepo, upgradeToVersion))
	if err == nil {
		fmt.Printf("\u2713 CLI successfully upgraded to version %s\n", upgradeToVersion)
	} else {
		fmt.Printf("\u274C CLI upgrade to version %s failed with error %v\n", upgradeToVersion, err)
	}

	return nil
}

var UpgraceCliCmd = &cobra.Command{
	Use:   "cli",
	Short: "upgrade cli",
	RunE:  UpgradeCli,
}

func init() {
	UpgraceCliCmd.Flags().StringVarP(&upgradeToVersion, "version",
		"v", "", "version to upgrade to")
}
