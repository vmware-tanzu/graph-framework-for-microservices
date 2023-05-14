package upgrade

import (
	"github.com/spf13/cobra"
)

var UpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "upgrade components in nexus-sdk",
}

func init() {
	UpgradeCmd.AddCommand(UpgradeCliCmd)
}
