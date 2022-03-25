package app

import (
	"fmt"

	"github.com/spf13/cobra"

	"gitlab.eng.vmware.com/nexus/cli/pkg/utils"
)

var (
	DMDir          string
	AppDir         string
	DatatmodelName string
	AppName        string
	RegistryURL    string
)

func Init(cmd *cobra.Command, args []string) error {
	envList := []string{}

	fmt.Println("XXXX:", args)

	if AppDir != "" {
		envList = append(envList, fmt.Sprintf("APP_DIR=%s", AppDir))
	}
	if DMDir != "" {
		envList = append(envList, fmt.Sprintf("DATAMODEL_DIR=%s", DMDir))
	}
	if DatatmodelName != "" {
		envList = append(envList, fmt.Sprintf("DATAMODEL=%s", DatatmodelName))
	}
	if AppName != "" {
		envList = append(envList, fmt.Sprintf("APPNAME=%s", AppName))
	}
	if RegistryURL != "" {
		envList = append(envList, fmt.Sprintf("REGISTRY=%s", RegistryURL))
	}

	// cd nexus/
	err := utils.SystemCommand(envList, "make", "app_init")
	if err != nil {
		return err
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
	InitCmd.Flags().StringVarP(&DMDir, "datamodel-dir",
		"m", "", "datamodel directory location.")
	InitCmd.Flags().StringVarP(&AppDir, "app-dir",
		"p", "", "app directory location.")
	InitCmd.Flags().StringVarP(&DatatmodelName, "datamodel",
		"d", "", "name of the datamodel")
	InitCmd.Flags().StringVarP(&AppName, "app",
		"a", "", "name of the application")
	InitCmd.Flags().StringVarP(&RegistryURL, "registry",
		"r", "", "registry url to publish docker image")
}
