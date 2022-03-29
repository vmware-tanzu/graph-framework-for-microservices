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
	ImportPath string
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

func InitOperation(cmd *cobra.Command, args []string) error {
	if DatatmodelName == "" {
		fmt.Print("Assuming datamodel name as default value: datamodel\n")
		DatatmodelName = "datamodel"
	}
	if GroupName == "" {
		fmt.Printf("Assuming group name as %s.com\n", DatatmodelName)
		GroupName = strings.TrimSuffix(fmt.Sprintf("%s.com", DatatmodelName), "\n")
	}

	fmt.Print("run this command outside of nexus home directory\n")
	if _, err := os.Stat(NEXUS_DIR); os.IsNotExist(err) {
		fmt.Printf("creating nexus home directory\n")
		manifestDir, err := utils.GitClone()
		if err != nil {
			return err
		}

		err = CopyDir(fmt.Sprintf("%s/nexus", manifestDir), NEXUS_DIR)
		if err != nil {
			return err
		}

		err = os.RemoveAll(manifestDir)
		if err != nil {
			return err
		}
	}
	os.Chdir(NEXUS_DIR)
	if _, err := os.Stat(DatatmodelName); err == nil {
		fmt.Printf("Datamodel %s already exists\n", DatatmodelName)
		return nil
	}
	err := CopyDir(".datamodel.templatedir", DatatmodelName)
	if err != nil {
		return fmt.Errorf("could not create datamodel due to %s\n", err)
	}
	err = utils.GoModInit(DatatmodelName)
	if err != nil {
		return err
	}

	importPath := "gitlab.eng.vmware.com/nsx-allspark_users/m7"
	values := TemplateValues{
		ImportPath: strings.TrimSuffix(string(importPath), "\n"),
		ModuleName: strings.TrimSuffix(DatatmodelName, "\n"),
		GroupName:  strings.TrimSuffix(GroupName, "\n"),
	}
	err = utils.RenderTemplateFiles(values, DatatmodelName)
	if err != nil {
		return fmt.Errorf("could not create datamodel due to %s\n", err)

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
}
