package app

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	. "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
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

func Init(cmd *cobra.Command, args []string) error {
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
	envList := []string{}
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
	if DOWNLOAD_APP != "false" {
		err := utils.DownloadFile(TEMPLATE_URL, Filename)
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
	err := utils.CreateNexusDirectory(NEXUS_DIR, NEXUS_TEMPLATE_URL)
	if err != nil {
		return fmt.Errorf("could not create nexus directory: %s", err)
	}

	if DatatmodelName != "" {
		if DatatmodelGroup == "" {
			DatatmodelGroup = fmt.Sprintf("%s.com", DatatmodelName)
		}
		envList = append(envList, fmt.Sprintf("DATAMODEL_GROUP=%s", DatatmodelGroup))
		err := utils.SystemCommand(cmd, utils.DATAMODEL_INIT_FAILED, envList, "make", "datamodel_init")
		if err != nil {
			return fmt.Errorf("error in creating new datamodel due to %s", err)
		}
		fmt.Println("\u2713 Application initialized successfully , Please continue to edit the datamodel and app controllers.")

		err = WriteToNexusDms(DatatmodelName, NexusDmProps{DatatmodelName, true})
		if err != nil {
			return fmt.Errorf("failed to write to nexus-dms.yaml")
		}
	} else {
		fmt.Println("\u2713 Application initialized successfully, Create Datamodel for consuming applications")
	}
	return nil
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a bare-bones Nexus application ready for operators to be added",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: Init,
}

func init() {
	InitCmd.Flags().StringVarP(&DatatmodelName, "datamodel-init",
		"d", "", "name of the datamodel to initialize")
	InitCmd.Flags().StringVarP(&DatatmodelImport, "datamodel-import",
		"i", "", "name of the datamodel to import")
	InitCmd.Flags().StringVarP(&AppName, "name",
		"n", "", "name of the application")
	InitCmd.Flags().StringVarP(&RegistryURL, "registry",
		"r", "284299419820.dkr.ecr.us-west-2.amazonaws.com/nexus/playground", "container registry url to publish docker images to")
	InitCmd.Flags().StringVarP(&DatatmodelGroup, "datamodel-group",
		"g", "", "group of the datamodel being initialized. only to be used in conjunction with --datamodel-init")
	err := cobra.MarkFlagRequired(InitCmd.Flags(), "name")

	if err != nil {
		fmt.Printf("init error: %v\n", err)
	}
}
