package app

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"gitlab.eng.vmware.com/nexus/cli/pkg/utils"
)

type TemplateValues struct {
	AppName       string
	ImageRegistry string
}

var (
	DMDir            string
	AppDir           string
	DatatmodelName   string
	DatatmodelGroup  string
	AppName          string
	RegistryURL      string
	DatatmodelImport string
)

var TEMPLATE_URL = "https://storage.googleapis.com/nexus-template-downloads/app-template.tar"
var filename = "app-template.tar"

func Init(cmd *cobra.Command, args []string) error {
	envList := []string{}
	fmt.Println("XXXX:", args)
	if DatatmodelName != "" {
		envList = append(envList, fmt.Sprintf("DATAMODEL=%s", DatatmodelName))
	}
	var DOWNLOAD_APP string = "true"
	files, _ := ioutil.ReadDir(".")
	for _, file := range files {
		if file.Name() == "PROJECT" {
			fmt.Printf("Skipping template download and rendering the app directory..")
			DOWNLOAD_APP = "false"
		}
	}
	if DOWNLOAD_APP != "false" {
		err := utils.DownloadFile(TEMPLATE_URL, filename)
		if err != nil {
			return fmt.Errorf("could not download template files due to %s", err)
		}

		file, err := os.Open(filename)
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
		_ = os.RemoveAll(filename)
	}
	if DatatmodelName != "" {
		err := utils.SystemCommand(envList, "make", "datamodel_init")
		if err != nil {
			return fmt.Errorf("error in creating new datamodel due to %s", err)
		}
	}
	_, err := os.Stat("nexus")
	if os.IsNotExist(err) {
		os.Mkdir("nexus", 0755)
	}
	return nil
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "intalls a sample application",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: Init,
}

func init() {
	InitCmd.Flags().StringVarP(&DatatmodelName, "datamodel-init",
		"d", "", "name of the datamodel")
	InitCmd.Flags().StringVarP(&DatatmodelImport, "datamodel-import",
		"i", "", "name of the datamodel")
	InitCmd.Flags().StringVarP(&AppName, "name",
		"n", "", "name of the application")
	InitCmd.Flags().StringVarP(&RegistryURL, "registry",
		"r", "284299419820.dkr.ecr.us-west-2.amazonaws.com/nexus/playground", "registry url to publish docker image")
	InitCmd.Flags().StringVarP(&DatatmodelGroup, "datamodel-group",
		"g", "", "group in datamodel")
	err := cobra.MarkFlagRequired(InitCmd.Flags(), "name")

	if err != nil {
		fmt.Printf("init error: %v\n", err)
	}
}
