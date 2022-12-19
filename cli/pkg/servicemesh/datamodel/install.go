package datamodel

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/common"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/log"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/prereq"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/utils"
)

var Namespace string
var DatamodelImage string
var Title string
var GraphqlPath string
var Force bool
var installPrerequisites []prereq.Prerequiste = []prereq.Prerequiste{
	prereq.KUBERNETES,
	prereq.KUBERNETES_VERSION,
}
var ImagePullSecret string

const (
	DatamodelJobSpecConfig string = "datamodel-job-spec"
	ToolsImage             string = "gcr.io/nsx-sm/tools:latest"
	NamespaceFlag          string = "namespace"
	ForceFlag              string = "force"
)

func CheckSpecAvailable(JobSpecConfigmap, namespace string) error {
	err := exec.Command("kubectl", "get", "cm", JobSpecConfigmap, "-n", namespace).Run()
	if err != nil {
		log.Errorf("Could not find the %s configmap due to %s - please install latest runtime", JobSpecConfigmap, err)
		return err
	}
	return nil
}

func InstallJob(DatamodelImage, DatamodelName, ImagePullsecret, Namespace, skipCRDInstallation, Title string, forceInstall bool) error {
	if DatamodelName == "" {
		ImageName := strings.Split(DatamodelImage, ":")[0]
		Name := strings.Split(ImageName, "/")
		DatamodelName = Name[len(Name)-1]
		DatamodelName = strings.ToLower(DatamodelName)
	}
	var IsImagePullSecret = false
	if ImagePullSecret != "" {
		checkSecretCommand := exec.Command("kubectl", "get", "secret", ImagePullSecret, "-n", Namespace)
		err := checkSecretCommand.Run()
		if err != nil {
			fmt.Printf("Please Create Secret %s before calling nexus datamodel install with --secret/-s option on namespace", ImagePullSecret)
			return err
		}
		IsImagePullSecret = true
	}
	data := common.Datamodel{
		DatamodelInstaller: common.DatamodelInstaller{
			Image: DatamodelImage,
			Name:  DatamodelName,
			Force: forceInstall,
		},
		IsImagePullSecret:   IsImagePullSecret,
		ImagePullSecret:     ImagePullSecret,
		SkipCRDInstallation: skipCRDInstallation,
		DatamodelTitle:      Title,
		GraphqlPath:         GraphqlPath,
	}
	if err := CheckSpecAvailable(DatamodelJobSpecConfig, Namespace); err != nil {
		return err
	}
	err := utils.RunDatamodelInstaller(DatamodelJobSpecConfig, Namespace, DatamodelName, data)
	if err != nil {
		log.Errorf("could not complete datamodel install due to: %s", err)
		return err
	}
	fmt.Printf("\u2713 Datamodel %s install successful\n", DatamodelName)
	return nil
}

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install specified datamodel's generated CRDs to the specified namespace",
	//Args:  cobra.ExactArgs(1),
}

func init() {
	InstallCmd.AddCommand(ImageCmd)
	InstallCmd.AddCommand(NameCmd)
	InstallCmd.PersistentFlags().StringVarP(&Title, "title",
		"", "", "title of the swaggerDocs for rest endpoints")
	InstallCmd.PersistentFlags().StringVarP(&GraphqlPath, "graphql-url",
		"", "", "Url where graphql plugin is available if any custom storage is used")
	InstallCmd.PersistentFlags().StringVarP(&Namespace, "namespace",
		"r", "", "name of the namespace to install to")
	InstallCmd.PersistentFlags().BoolVarP(&Force, ForceFlag,
		"f", false, "forcefully install backward incompatible changes")
	err := cobra.MarkFlagRequired(InstallCmd.Flags(), "namespace")
	if err != nil {
		log.Debugf("please provide namespace: %v", err)
	}
}
