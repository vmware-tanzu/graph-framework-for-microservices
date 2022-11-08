package version

import (
	"regexp"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/common"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/log"
)

// IsNexusCliUpdateAvailable returns a bool indicating if a newer CLI version is available
// also returns the latest CLI version available
func IsNexusCliUpdateAvailable() (bool, string) {
	latestNexusVersion, err := GetLatestNexusVersion()
	latestNexusVersionStripped := strings.TrimPrefix(latestNexusVersion, "v")
	if err != nil {
		log.Debugf("Error while trying to fetch latest available nexus version: %s\n", err)
		return false, ""
	}

	var currentValues NexusValues
	if err = GetNexusValues(&currentValues); err != nil {
		log.Debugf("Get current Nexus artifacts versions failed with error %v", err)
		return false, ""
	}

	if latestNexusVersionStripped != common.VERSION {
		log.Debugf("Latest available Nexus CLI version: %s\n", latestNexusVersion)
		log.Debugf("Current Nexus CLI version: %s\n", common.VERSION)
	}
	latestNexusVersionSemver, err := semver.NewVersion(latestNexusVersionStripped)
	if err != nil {
		log.Debugf("Latest Nexus version is incorrectly formatted: %v\n", err)
		return false, ""
	}
	regexM := regexp.MustCompile(`(v?\d+.\d+.\d+$)`)
	versionString := regexM.FindStringSubmatch(common.VERSION)
	if len(versionString) != 2 {
		log.Debugf("Current Nexus Version is not official version: %s\n, Please use --version to force upgrade", common.VERSION)
		return false, ""

	}

	currentNexusVersionSemver, err := semver.NewVersion(format(common.VERSION))
	if err != nil {
		log.Debugf("Current Nexus version is incorrectly formatted: %v\n", err)
		return false, ""
	}

	return currentNexusVersionSemver.LessThan(*latestNexusVersionSemver), latestNexusVersion
}
