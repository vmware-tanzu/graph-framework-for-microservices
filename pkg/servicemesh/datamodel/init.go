package datamodel

import (
	"fmt"
	"io/ioutil"
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

func CopyDir(src string, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}

	file, err := f.Stat()
	if err != nil {
		return err
	}
	if !file.IsDir() {
		return fmt.Errorf("Source " + file.Name() + " is not a directory!")
	}

	err = os.Mkdir(dest, 0755)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, f := range files {

		if f.IsDir() {

			err = CopyDir(src+"/"+f.Name(), dest+"/"+f.Name())
			if err != nil {
				return err
			}

		}

		if !f.IsDir() {

			content, err := ioutil.ReadFile(src + "/" + f.Name())
			if err != nil {
				return err

			}

			err = ioutil.WriteFile(dest+"/"+f.Name(), content, 0755)
			if err != nil {
				return err

			}

		}

	}
	return nil
}

func createDatamodel() error {
	os.Chdir(NEXUS_DIR)
	if _, err := os.Stat(DatatmodelName); err == nil {
		fmt.Printf("Datamodel %s already exists\n", DatatmodelName)
		return nil
	}

	fmt.Printf("creating %s datamodel as part of initialization\n", DatatmodelName)
	err := os.Mkdir(DatatmodelName, 0755)
	if err != nil {
		return err
	}
	os.Chdir(DatatmodelName)
	err = utils.DownloadFile(DATAMODEL_TEMPLATE_URL, "datamodel.tar")
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
	err = utils.GoModInit(DatatmodelName)
	if err != nil {
		return err
	}

	values := TemplateValues{
		ModuleName: strings.TrimSuffix(DatatmodelName, "\n"),
		GroupName:  strings.TrimSuffix(GroupName, "\n"),
	}
	err = utils.RenderTemplateFiles(values, DatatmodelName, "")
	if err != nil {
		return fmt.Errorf("could not create datamodel due to %s\n", err)

	}
	os.Chdir("..")
	return nil
}

func InitOperation(cmd *cobra.Command, args []string) error {
	if GroupName == "" {
		fmt.Printf("Assuming group name as %s.com\n", DatatmodelName)
		GroupName = strings.TrimSuffix(fmt.Sprintf("%s.com", DatatmodelName), "\n")
	}
	err := utils.CreateNexusDirectory(NEXUS_DIR, NEXUS_TEMPLATE_URL)
	if err != nil {
		return fmt.Errorf("could not create nexus directory..")
	}
	os.Chdir(NEXUS_DIR)
	if _, err := os.Stat("helloworld"); os.IsNotExist(err) {
		fmt.Printf("creating helloworld example datamodel as part of initialization\n")
		err = os.Mkdir("helloworld", 0755)
		if err != nil {
			return err
		}
		os.Chdir("helloworld")
		err := utils.DownloadFile(HELLOWORLD_URL, "helloworld.tar")
		if err != nil {
			return fmt.Errorf("could not download template files due to %s\n", err)
		}

		file, err := os.Open("helloworld.tar")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer file.Close()
		err = utils.Untar(".", file)
		if err != nil {
			return fmt.Errorf("could not unarchive template files due to %s\n", err)
		}
		os.Remove("helloworld.tar")
		os.Chdir("..")
	}
	err = createDatamodel()
	if err != nil {
		return err
	}
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
	err := cobra.MarkFlagRequired(InitCmd.Flags(), "name")
	if err != nil {
		fmt.Printf("name is mandatory")
	}
}
