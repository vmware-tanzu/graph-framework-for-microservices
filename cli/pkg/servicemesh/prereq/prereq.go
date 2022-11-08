package prereq

import (
	"bytes"
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/common"
	nexusCommon "github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/common"
	nexusVersion "github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/version"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/utils"
)

const goMinVersion = "1.17"

const (
	k8sMinVersion  = "1.19"
	helmMinVersion = "3.4"
)

var All bool
var push bool
var registry string

type Prerequiste int

const (
	// Adding validation enum here to add multiple Prerequistes.
	DOCKER Prerequiste = iota
	KUBERNETES
	KUBERNETES_VERSION
	GOLANG_VERSION
	GOPATH
)

type PrerequisteMeta struct {
	what                  string
	verify                func() (bool, error)
	AdditionalDescription string
	Version               string
	Validator             Prerequiste
	Always                bool
}

var preReqs = map[Prerequiste]PrerequisteMeta{
	DOCKER: {
		what:                  "docker",
		Always:                true,
		AdditionalDescription: "docker daemon should be running on the host",
		verify: func() (bool, error) {
			_, err := exec.Command("docker", "ps").Output()
			if err != nil {
				return false, fmt.Errorf("verify if docker is running failed with error %v", err)
			}
			return true, nil
		},
	},
	GOLANG_VERSION: {
		what:    "go",
		Always:  true,
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
	KUBERNETES: {
		what:                  "kubernetes",
		Always:                false,
		AdditionalDescription: "kubernetes cluster should be reachable via kubectl",
		verify: func() (bool, error) {
			_, err := exec.Command("kubectl", "get", "ns").Output()
			if err != nil {
				return false, fmt.Errorf("verifying running kubernetes cluster failed with error %v", err)
			}
			return true, nil
		},
	},
	KUBERNETES_VERSION: {
		what:                  "kubernetes version",
		Always:                false,
		AdditionalDescription: fmt.Sprintf("Kubernetes version should be atleast %s", k8sMinVersion),
		verify: func() (bool, error) {
			versionStringBytes, _ := exec.Command("kubectl", "version", "--short=true").Output()
			cmd := exec.Command("tail", "-n", "1")
			cmd.Stdin = strings.NewReader(string(versionStringBytes))
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				return false, fmt.Errorf("could not get k8s version string")
			}

			re := regexp.MustCompile(`Server Version: ([a-z0-9][^\s]*)`)
			serverVersion := re.FindStringSubmatch(out.String())
			if len(serverVersion) == 0 {
				return false, fmt.Errorf("unable to get k8s version from output: %v", out.String())
			}
			v1min, errMinVersion := version.NewVersion(k8sMinVersion)
			if errMinVersion != nil {
				return false, fmt.Errorf("verify version of kubernetes cluster is running failed with error on minVersion formation %v", errMinVersion)
			}
			v1, errVersion := version.NewVersion(strings.TrimPrefix(serverVersion[1], "v"))
			if errVersion != nil {
				return false, fmt.Errorf("verify version of kubernetes cluster is running failed with error current version formation %v", errVersion)
			}
			if v1.LessThan(v1min) {
				return false, fmt.Errorf("K8s Version should be atleast %s, current Version is %s", k8sMinVersion, v1)
			}
			return true, nil
		},
	},
	GOPATH: {
		what:                  "GOPATH",
		AdditionalDescription: "app workspace should be in GOPATH",
		verify: func() (bool, error) {

			gopath := os.Getenv("GOPATH")
			if gopath == "" {
				gopath = build.Default.GOPATH
			}

			workspacePath, err := os.Getwd()
			if err != nil {
				log.Println(err)
			}

			up := ".." + string(os.PathSeparator)
			// path-comparisons using filepath.Abs don't work reliably according to docs (no unique representation).
			if rel, err := filepath.Rel(gopath, workspacePath); err == nil {
				if !strings.HasPrefix(rel, up) && rel != ".." {
					return true, nil
				} else {
					return false, fmt.Errorf("workspace %s is not in GOPATH %s", workspacePath, gopath)
				}
			}
			fmt.Printf("Verifying if workspace is in GOPATH failed with error %s. Ignoring verification.", err)
			return true, nil
		},
	},
}

func PreReqVerify(cmd *cobra.Command, args []string) error {
	for _, util := range preReqs {
		if utils.VerifyAll(cmd) || util.Always {
			if ok, err := util.verify(); ok {
				fmt.Printf("\u2705 %s %s %s\n", util.what, util.Version, util.AdditionalDescription)
			} else {
				fmt.Printf("\u274C %s %s verify failed with err: %v\n", util.what, util.Version, err)
			}
		}
	}
	return nil
}

func printPreReq(req PrerequisteMeta) {
	fmt.Printf("\u2023 %s", req.what)
	if len(req.Version) > 0 {
		fmt.Printf(" (version: %s)", req.Version)
	}
	if len(req.AdditionalDescription) > 0 {
		fmt.Printf(" [ %s ]", req.AdditionalDescription)
	}
	fmt.Println()
}

func PreReqList(cmd *cobra.Command, args []string) error {
	for _, util := range preReqs {
		printPreReq(util)
	}
	return nil
}

func PreReqVerifyOnDemand(reqs []Prerequiste) error {
	for _, current := range reqs {
		util := preReqs[current]
		if ok, err := util.verify(); !ok {
			return err
		}
	}
	return nil
}

// PreReqListOnDemand lists the prerequisites and exits
func PreReqListOnDemand(reqs []Prerequiste) {
	for _, current := range reqs {
		if util, found := preReqs[current]; found {
			printPreReq(util)
		}
	}
	os.Exit(0)
}

func PreReqImages(cmd *cobra.Command, args []string) error {
	var values nexusVersion.NexusValues
	if err := nexusVersion.GetNexusValues(&values); err != nil {
		return utils.GetCustomError(utils.RUNTIME_PREREQUISITE_IMAGE_PREP_FAILED,
			fmt.Errorf("could not pull runtime deps images %s", err)).Print().ExitIfFatalOrReturn()
	}
	for _, manifest := range nexusCommon.TagsList {
		if manifest.ImageName != "" {
			Image := fmt.Sprintf("%s/%s", common.HarborRepo, manifest.ImageName)
			versionToStr := reflect.ValueOf(values).FieldByName(manifest.FieldName).Field(0).String()
			if os.Getenv(versionToStr) != "" {
				versionTo := os.Getenv(versionToStr)
				Image = fmt.Sprintf("%s:%s", Image, versionTo)
				fmt.Printf("Pulling image: %s\n", Image)
				err := utils.SystemCommand(cmd, utils.RUNTIME_PREREQUISITE_IMAGE_PREP_FAILED, []string{}, "docker", "pull", Image)
				if err != nil {
					fmt.Printf("could not pull image: %s\n", Image)
					return err
				}
				if cmd.Flags().Lookup("push").Changed {
					if registry == "" {
						fmt.Println("provide registry with --registry/-r to push the images to.")
					} else {
						newImage := fmt.Sprintf("%s/%s:%s", registry, manifest.ImageName, versionTo)
						fmt.Printf("Pushing Image: %s\n", newImage)
						err := utils.SystemCommand(cmd, utils.RUNTIME_PREREQUISITE_IMAGE_PREP_FAILED, []string{}, "docker", "tag", Image, newImage)
						if err != nil {
							fmt.Printf("could not tag image: %s\n", newImage)
							return err
						}
						err = utils.SystemCommand(cmd, utils.RUNTIME_PREREQUISITE_IMAGE_PREP_FAILED, []string{}, "docker", "push", newImage)
						if err != nil {
							fmt.Printf("could not push image: %s\n", newImage)
							return err
						}
					}
				}
			}
		}
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
var PreReqImageCmd = &cobra.Command{
	Use:   "images",
	Short: "pull all images needed for runtime",
	RunE:  PreReqImages,
}

var PreReqCmd = &cobra.Command{
	Use:   "prereq",
	Short: "pre-requisites for a successful nexus-sdk experience",
}

func init() {
	PreReqCmd.AddCommand(PreReqListCmd)
	PreReqCmd.AddCommand(PreReqVerifyCmd)
	PreReqCmd.AddCommand(PreReqImageCmd)
	PreReqCmd.PersistentFlags().BoolVarP(&All, "all", "", false, "For validation check")
	PreReqImageCmd.PersistentFlags().BoolVarP(&push, "push", "p", false, "For publishing images")
	PreReqImageCmd.PersistentFlags().StringVarP(&registry, "registry", "r", "", "registry to push images")
}
