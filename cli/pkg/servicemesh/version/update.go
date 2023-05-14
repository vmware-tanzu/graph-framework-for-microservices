package version

import (
	"strings"

	"github.com/coreos/go-semver/semver"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/log"
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

	if latestNexusVersionStripped != currentValues.NexusCli.Version {
		log.Debugf("Latest available Nexus CLI version: %s\n", latestNexusVersion)
		log.Debugf("Current Nexus CLI version: %s\n", currentValues.NexusCli.Version)
	}

	latestNexusVersionSemver, err := semver.NewVersion(latestNexusVersionStripped)
	if err != nil {
		log.Debugf("Latest Nexus version is incorrectly formatted: %v\n", err)
		return false, ""
	}

	currentNexusVersionSemver, err := semver.NewVersion(format(currentValues.NexusCli.Version))
	if err != nil {
		log.Debugf("Current Nexus version is incorrectly formatted: %v\n", err)
		return false, ""
	}

	return currentNexusVersionSemver.LessThan(*latestNexusVersionSemver), latestNexusVersion
}
