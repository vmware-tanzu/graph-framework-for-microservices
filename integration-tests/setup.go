package main

import "os/exec"

func SetupEnv() {
	// check if nexus prereqs are satisfied
	_, err := exec.Command("nexus", "prereq", "verify").Output()
	CheckIfError(err)

	_, err = exec.Command("nexus", "config", "set", "--debug-always", "--disable-upgrade-prompt", "--skip-upgrade-check").Output()
	CheckIfError(err)
}
