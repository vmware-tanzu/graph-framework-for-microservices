package prereq

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
)

const goMinVersion = "1.17"

var preReqs = []struct {
	What                  string
	Version               string
	AdditionalDescription string
	verify                func() (bool, error)
}{
	{
		What:    "go",
		Version: goMinVersion,
		verify: func() (bool, error) {
			out, err := exec.Command("go", "version").Output()
			if err != nil {
				return false, fmt.Errorf("verify go version failed with error %v", err)
			}
			re := regexp.MustCompile(`go[0-9][^\s]*`)
			match := re.FindStringSubmatch(string(out))
			if len(match) == 0 {
				return false, fmt.Errorf("unable to get go version from output: %v", string(out))
			}

			v1, errMinVersion := version.NewVersion(goMinVersion)
			if errMinVersion != nil {
				return false, fmt.Errorf("parse min go version failed with error %v", errMinVersion)
			}
			v2, errCurrVersion := version.NewVersion(strings.Trim(match[0], "go"))
			if errCurrVersion != nil {
				return false, fmt.Errorf("parse current go version failed with error %v", errCurrVersion)
			}

			if v2.LessThan(v1) {
				return false, fmt.Errorf("go version %s is less than %s", string(match[0]), goMinVersion)
			}

			return true, nil
		},
	},
	{
		What:                  "docker",
		AdditionalDescription: "docker daemon should be running on the host",
		verify: func() (bool, error) {
			_, err := exec.Command("docker", "ps").Output()
			if err != nil {
				return false, fmt.Errorf("verify if docker is running failed with error %v", err)
			}
			return true, nil
		},
	},
}

func PreReqVerify(cmd *cobra.Command, args []string) error {
	for _, util := range preReqs {
		if ok, err := util.verify(); ok {
			fmt.Printf("\u2705 %s %s\n", util.What, util.Version)
		} else {
			fmt.Printf("\u274C %s %s verify failed with err: %v\n", util.What, util.Version, err)
		}
	}
	return nil
}

func PreReqList(cmd *cobra.Command, args []string) error {
	for _, util := range preReqs {
		fmt.Printf("\u2023 %s", util.What)
		if len(util.Version) > 0 {
			fmt.Printf(" (version: %s)", util.Version)
		}
		if len(util.AdditionalDescription) > 0 {
			fmt.Printf(" [ %s ]", util.AdditionalDescription)
		}
		fmt.Println()
	}
	return nil
}

var PreReqVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "verify all pre-requisites",
	RunE:  PreReqVerify,
}

var PreReqListCmd = &cobra.Command{
	Use:   "list",
	Short: "list all pre-requisites",
	RunE:  PreReqList,
}

var PreReqCmd = &cobra.Command{
	Use:   "prereq",
	Short: "pre-requisites for a successful nexus-sdk experience",
}

func init() {
	PreReqCmd.AddCommand(PreReqListCmd)
	PreReqCmd.AddCommand(PreReqVerifyCmd)
}
