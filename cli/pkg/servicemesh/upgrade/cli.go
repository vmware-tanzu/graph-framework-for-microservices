package upgrade

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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

func DoUpgradeCli(versionstr string, cmd *cobra.Command) error {
	// Check a new tag is released
	// download the binary and copy to execpath
	var cliV bool
	if versionstr == "" {
		cliV, versionstr = version.IsNexusCliUpdateAvailable()
		if !cliV {
			fmt.Printf("Skipping upgrade")
			return nil
		}
	}
	envList := common.GetEnvList()
	OS := runtime.GOOS
	ARCH := runtime.GOARCH
	URL := fmt.Sprintf("https://github.com/vmware-tanzu/graph-framework-for-microservices/releases/download/%s/nexus-%s_%s", versionstr, OS, ARCH)
	err := utils.DownloadFile(URL, "nexusbin")
	if err != nil {
		fmt.Printf("\u274C CLI Upgrade to version %s Failed with error %v \n", versionstr, err)
		return err
	}
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)
	err = utils.SystemCommand(cmd, utils.CLI_UPGRADE_FAILED, envList, "mv", "nexusbin", fmt.Sprintf("%s/nexus", exPath))
	if err == nil {
		err = utils.SystemCommand(cmd, utils.CLI_UPGRADE_FAILED, envList, "chmod", "+x", fmt.Sprintf("%s/nexus", exPath))
		if err != nil {
			fmt.Printf("\u274C CLI upgrade to version %s failed with error %v\n", versionstr, err)
			return err
		}
		fmt.Printf("\u2713 CLI successfully upgraded to version %s\n", versionstr)
	} else {
		fmt.Printf("\u274C CLI upgrade to version %s failed with error %v\n", versionstr, err)
		return err
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
