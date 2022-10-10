package app

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/log"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	. "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/prereq"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

type TemplateValues struct {
	AppName       string
	ImageRegistry string
}

var (
	DatatmodelName   string
	DatatmodelGroup  string
	AppName          string
	RegistryURL      string
	DatatmodelImport string
)

var prerequisites []prereq.Prerequiste = []prereq.Prerequiste{
	prereq.GOLANG_VERSION,
}

func Init(cmd *cobra.Command, args []string) error {

	if utils.ListPrereq(cmd) {
		prereq.PreReqListOnDemand(prerequisites)
		return nil
	}

	empty, _ := utils.IsDirEmpty(".")
	if empty == false {
		var input string
		fmt.Println("Current Directory is not empty do you want to continue[y/n]: ")
		fmt.Scanln(&input)
		if input == "n" {
			fmt.Println("Aborting nexus app init operation.")
			return nil
		}
	}
	envList := common.GetEnvList()
	if DatatmodelName != "" {
		envList = append(envList, fmt.Sprintf("DATAMODEL=%s", DatatmodelName))
	}
	var DOWNLOAD_APP string = "true"
	files, _ := ioutil.ReadDir(".")
	for _, file := range files {
		if file.Name() == "PROJECT" {
			fmt.Printf("Skipping template download and rendering the app directory..\n")
			DOWNLOAD_APP = "false"
		}
	}
	appVersion, err := utils.GetTagVersion("NexusAppTemplates", "NEXUS_APP_TEMPLATE_VERSION")
	if err != nil {
		return utils.GetCustomError(utils.APPLICATION_INIT_PREREQ_FAILED,
			fmt.Errorf("could not download the app manifests due to %s", err)).Print().ExitIfFatalOrReturn()
	}

	log.Debugf("Using App template Version: %s\n", appVersion)

	if DOWNLOAD_APP != "false" {

		err := utils.DownloadFile(fmt.Sprintf(TEMPLATE_URL, appVersion), Filename)
		if err != nil {
			return fmt.Errorf("could not download template files due to %s", err)
		}

		file, err := os.Open(Filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer file.Close()
		err = utils.Untar(".", file)
		if err != nil {
			return fmt.Errorf("could not unarchive template files due to %s", err)
		}
		data := TemplateValues{
			AppName:       AppName,
			ImageRegistry: RegistryURL,
		}
		err = utils.RenderTemplateFiles(data, ".", "nexus")
		if err != nil {
			return fmt.Errorf("error in rendering template files due to %s", err)
		}
		_ = os.RemoveAll(Filename)
	}

	nexusVersion, err := utils.GetTagVersion("NexusDatamodelTemplates", "NEXUS_DATAMODEL_TEMPLATE_VERSION")
	if err != nil {
		return utils.GetCustomError(utils.APPLICATION_INIT_PREREQ_FAILED,
			fmt.Errorf("could not download the app manifests due to %s", err)).Print().ExitIfFatalOrReturn()
	}

	fmt.Printf("Using Nexus template version: %s\n", nexusVersion)
	err = utils.CreateNexusDirectory(NEXUS_DIR, fmt.Sprintf(NEXUS_TEMPLATE_URL, nexusVersion))
	if err != nil {
		return fmt.Errorf("could not create nexus directory: %s", err)
	}
	fmt.Println("\u2713 Application initialized successfully")
	return nil
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a bare-bones Nexus application ready for operators to be added",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if utils.ListPrereq(cmd) {
			return nil
		}

		if utils.SkipPrereqCheck(cmd) {
			return nil
		}
		return prereq.PreReqVerifyOnDemand(prerequisites)
	},
	RunE: Init,
}

func init() {
	InitCmd.Flags().StringVarP(&AppName, "name",
		"n", "", "name of the application")
	InitCmd.Flags().StringVarP(&RegistryURL, "registry",
		"r", "284299419820.dkr.ecr.us-west-2.amazonaws.com/nexus/playground", "container registry url to publish docker images to")
	err := cobra.MarkFlagRequired(InitCmd.Flags(), "name")

	if err != nil {
		fmt.Printf("init error: %v\n", err)
	}
}
