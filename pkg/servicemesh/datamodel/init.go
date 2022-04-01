package datamodel

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	. "gitlab.eng.vmware.com/nexus/cli/pkg/common"
	"gitlab.eng.vmware.com/nexus/cli/pkg/utils"
)

type TemplateValues struct {
	ModuleName string
	GroupName  string
}

var DatatmodelName string
var GroupName string

func createDatamodel(DatatmodelName string, URL string, Render bool, standalone bool) error {
	var Directory string
	if standalone == false {
		err := utils.GoToNexusDirectory()
		if err != nil {
			fmt.Printf("Could not locate nexusDirectory\n")
			return err
		}
		if _, err := os.Stat(DatatmodelName); err == nil {
			fmt.Printf("Datamodel %s already exists\n", DatatmodelName)
			return nil
		}

		fmt.Printf("creating %s datamodel as part of initialization\n", DatatmodelName)
		err = os.Mkdir(DatatmodelName, 0755)
		if err != nil {
			return err
		}
		os.Chdir(DatatmodelName)
	}
	err := utils.DownloadFile(URL, "datamodel.tar")
	if err != nil {
		return fmt.Errorf("could not download template files due to %s\n", err)
	}

	file, err := os.Open("datamodel.tar")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()
	err = utils.Untar(".", file)
	if err != nil {
		return fmt.Errorf("could not unarchive template files due to %s", err)
	}
	os.Remove("datamodel.tar")
	if DatatmodelName != "" {
		if DatatmodelName != "helloworld" {
			err := utils.GoModInit(DatatmodelName)
			if err != nil {
				return err
			}
		} else {
			os.Chdir("..")
		}
		Directory = DatatmodelName
	} else {
		err := utils.GoModInit(DatatmodelName)
		if err != nil {
			return err
		}
		DatatmodelName, err = utils.GetModuleName("")
		if err != nil {
			return err
		}
		fmt.Printf("Current Datamodel name: %s\n", DatatmodelName)
		Directory, _ = os.Getwd()
	}
	if Render {
		values := TemplateValues{
			ModuleName: strings.TrimSuffix(DatatmodelName, "\n"),
			GroupName:  strings.TrimSuffix(GroupName, "\n"),
		}
		err = utils.RenderTemplateFiles(values, Directory, "")
		if err != nil {
			return fmt.Errorf("could not create datamodel due to %s\n", err)
		}
	}
	if standalone != true {
		fmt.Printf("Storing current datamodel as default datamodel\n")
		err = utils.StoreCurrentDatamodel(DatatmodelName)
		if err != nil {
			return err
		}
	}
	return nil
}

func InitOperation(cmd *cobra.Command, args []string) error {
	var standalone bool = false
	if DatatmodelName == "" {
		if GroupName == "" {
			fmt.Println("You need to provide a groupname if datamodelname is not provided : ")
			fmt.Scanln(&GroupName)
			if GroupName == "" {
				fmt.Println("Please provide a non-empty groupname")
				return nil
			}
			empty, _ := utils.IsDirEmpty(".")
			if empty == false {
				var input string
				fmt.Println("Current Directory is not empty do you want to continue to initialize datamodel [y/n]: ")
				fmt.Scanln(&input)
				if input == "n" {
					fmt.Println("Aborting datamodel initialization operation.")
					return nil
				}
			}
			standalone = true
		}
	}
	if standalone == false {
		if GroupName == "" {
			fmt.Printf("Assuming group name as %s.com\n", DatatmodelName)
			GroupName = strings.TrimSuffix(fmt.Sprintf("%s.com", DatatmodelName), "\n")
		}
		err := utils.CreateNexusDirectory(NEXUS_DIR, NEXUS_TEMPLATE_URL)
		if err != nil {
			return fmt.Errorf("could not create nexus directory..")
		}
	}
	if DatatmodelName == "helloworld" {
		err := createDatamodel(DatatmodelName, HELLOWORLD_URL, false, standalone)
		if err != nil {
			return err
		}
	} else {
		err := createDatamodel(DatatmodelName, DATAMODEL_TEMPLATE_URL, true, standalone)
		if err != nil {
			return err
		}
	}
	fmt.Printf("\u2713 Datamodel %s initialized successfully\n", DatatmodelName)
	return nil

}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize the tenant-manifests directory",
	//Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: InitOperation,
}

func init() {
	InitCmd.Flags().StringVarP(&DatatmodelName, "name", "n", "", "name of the datamodel to be created")
	InitCmd.Flags().StringVarP(&GroupName, "group", "g", "", "group for the datamodel to be created")
}
