package datamodel

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/log"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/prereq"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

var Namespace string
var DatamodelImage string
var Title string
var GraphqlPath string
var installPrerequisites []prereq.Prerequiste = []prereq.Prerequiste{
	prereq.KUBERNETES,
	prereq.KUBERNETES_VERSION,
}
var ImagePullSecret string

const (
	DatamodelJobSpecConfig string = "datamodel-job-spec"
	ToolsImage             string = "gcr.io/nsx-sm/tools:latest"
	NamespaceFlag          string = "namespace"
)

func CheckSpecAvailable(JobSpecConfigmap, namespace string) error {
	err := exec.Command("kubectl", "get", "cm", JobSpecConfigmap, "-n", namespace).Run()
	if err != nil {
		log.Errorf("Could not find the %s configmap due to %s - please install latest runtime", JobSpecConfigmap, err)
		return err
	}
	return nil
}

func InstallJob(DatamodelImage, DatamodelName, ImagePullsecret, Namespace, skipCRDInstallation, Title string) error {
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
	InstallCmd.AddCommand(DirCmd)
	InstallCmd.PersistentFlags().StringVarP(&Title, "title",
		"", "", "title of the swaggerDocs for rest endpoints")
	InstallCmd.PersistentFlags().StringVarP(&GraphqlPath, "graphql-url",
		"", "", "Url where graphql plugin is available if any custom storage is used")
	InstallCmd.PersistentFlags().StringVarP(&Namespace, "namespace",
		"r", "", "name of the namespace to install to")
	err := cobra.MarkFlagRequired(InstallCmd.Flags(), "namespace")
	if err != nil {
		log.Debugf("please provide namespace: %v", err)
	}
}
