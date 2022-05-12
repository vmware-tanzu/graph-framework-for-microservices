package version

import (
	"fmt"
	"strings"

	"github.com/coreos/go-semver/semver"
)

// IsNexusCliUpdateAvailable returns a bool indicating if a newer CLI version is available
// also returns the latest CLI version available
func IsNexusCliUpdateAvailable(enableDebugLogs bool) (bool, string) {
	latestNexusVersion, err := GetLatestNexusVersion()
	latestNexusVersionStripped := strings.TrimPrefix(latestNexusVersion, "v")
	if err != nil {
		if enableDebugLogs {
			fmt.Printf("Error while trying to fetch latest available nexus version: %s\n", err)
		}
		return false, ""
	}

	var currentValues NexusValues
	if err = GetNexusValues(&currentValues); err != nil {
		if enableDebugLogs {
			fmt.Println(err)
		}
		return false, ""
	}

	if enableDebugLogs && latestNexusVersionStripped != currentValues.NexusCli.Version {
		fmt.Printf("Latest available Nexus CLI version: %s\n", latestNexusVersion)
		fmt.Printf("Current Nexus CLI version: %s\n", currentValues.NexusCli.Version)
	}

	latestNexusVersionSemver, err := semver.NewVersion(latestNexusVersionStripped)
	if err != nil {
		if enableDebugLogs {
			fmt.Printf("Latest Nexus version is incorrectly formatted: %v\n", err)
		}
		return false, ""
	}

	currentNexusVersionSemver, err := semver.NewVersion(format(currentValues.NexusCli.Version))
	if err != nil {
		if enableDebugLogs {
			fmt.Printf("Current Nexus version is incorrectly formatted: %v\n", err)
		}
		return false, ""
	}

	return currentNexusVersionSemver.LessThan(*latestNexusVersionSemver), latestNexusVersion
}
