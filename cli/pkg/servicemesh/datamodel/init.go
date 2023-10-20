package datamodel

import (
	"fmt"
	"os"
	"strings"

	"github.com/otiai10/copy"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/log"

	"github.com/spf13/cobra"
	. "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
)

const (
	localDatamodelFlag = "local"
)

type TemplateValues struct {
	ModuleName string
	GroupName  string
}

var DatamodelName string
var GroupName string
var localDatamodel bool
var dockerRepo string
var BuildDockerImg bool
var localDir string

func createDatamodel(dmName string, DatamodelTarballUrl string, Render bool, standalone bool) error {
	var Directory string
	var err error
	if !standalone {
		err = utils.GoToNexusDirectory()
		if err != nil {
			fmt.Printf("Could not locate nexusDirectory\n")
			return err
		}
		if _, err = os.Stat(dmName); err == nil {
			fmt.Printf("Datamodel %s already exists\n", dmName)
			return nil
		}

		fmt.Printf("creating %s datamodel as part of initialization\n", dmName)
		err = os.Mkdir(dmName, 0755)
		if err != nil {
			return err
		}
		err = os.Chdir(dmName)
		if err != nil {
			return err
		}
	}

	if localDir != "" {
		// completely arbitrary paths
		sourceDir := localDir + "/datamodel-templates/nexus/.datamodel.templatedir"
		fmt.Println("sourceDir:", sourceDir)
		destDir, _ := os.Getwd()
		fmt.Println("Working directory:", destDir)
		fmt.Println("copy using library:", copy.Copy(sourceDir, destDir))
	} else {
		err = utils.DownloadFile(DatamodelTarballUrl, "datamodel.tar")
		if err != nil {
			return fmt.Errorf("could not download template files due to %s", err)
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
	}

	if !standalone {
		if dmName != "helloworld" {
			err = utils.GoModInit(dmName, false)
			if err != nil {
				return err
			}
		} else {
			err = os.Chdir("..")
			if err != nil {
				return err
			}
		}
		Directory = dmName
	} else {
		err = utils.GoModInit(dmName, true)
		if err != nil {
			return err
		}
		fmt.Printf("Current Datamodel name: %s\n", dmName)
		Directory, _ = os.Getwd()
	}
	if Render {
		values := TemplateValues{
			ModuleName: strings.TrimSuffix(dmName, "\n"),
			GroupName:  strings.TrimSuffix(GroupName, "\n"),
		}
		err = utils.RenderTemplateFiles(values, Directory, ".git")
		if err != nil {
			return fmt.Errorf("could not create datamodel due to %s", err)
		}
	}
	if standalone {
		fmt.Printf("Storing current datamodel as default datamodel\n")
		err = utils.StoreCurrentDatamodel(dmName)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkIfDirectoryEmpty(standalone bool, DatamodelName string) {

	// TODO: Check if directory with datamodelname under nexus directory already exists.
	if standalone {
		empty, _ := utils.IsDirEmpty(".")
		if !empty {
			_, err := os.Stat("go.mod")
			if err == nil {
				// TODO: standard error logs
				fmt.Println("Datamodel already initialized with go.mod file, Please delete go.mod file or create a empty folder")
				os.Exit(1)
			}

			var input string
			fmt.Println("Current Directory is not empty do you want to continue to initialize datamodel [y/n]: ")
			fmt.Scanln(&input)
			if input == "n" {
				fmt.Println("Aborting datamodel initialization operation.")
				os.Exit(0)
			}
		}
	}
}

func InitOperation(cmd *cobra.Command, args []string) error {
	datamodelVersion, err := utils.GetTagVersion("NexusDatamodelTemplates", "NEXUS_DATAMODEL_TEMPLATE_VERSION")
	if err != nil {
		return utils.GetCustomError(utils.DATAMODEL_INIT_FAILED,
			fmt.Errorf("could not download the datamodel manifests due to %s", err)).Print().ExitIfFatalOrReturn()
	}

	log.Debugf("Using datamodel template Version: %s\n", datamodelVersion)

	dmName := DatamodelName
	fmt.Printf("Datamodel name: %s\n", dmName)
	checkIfDirectoryEmpty(!localDatamodel, DatamodelName)

	if localDatamodel {
		err := utils.CreateNexusDirectory(NEXUS_DIR, fmt.Sprintf(NEXUS_TEMPLATE_URL, datamodelVersion))
		if err != nil {
			// TODO standard log library error
			return fmt.Errorf("could not create nexus directory")
		}
	}

	if dmName == "helloworld" {
		err := createDatamodel(dmName, fmt.Sprintf(HELLOWORLD_URL, datamodelVersion), false, !localDatamodel)
		if err != nil {
			return err
		}
	} else {
		err := createDatamodel(dmName, fmt.Sprintf(DATAMODEL_TEMPLATE_URL, datamodelVersion), true, !localDatamodel)
		if err != nil {
			return err
		}
	}
	if dockerRepo != "" {
		if localDatamodel {
			err = os.Chdir(dmName)
			if err != nil {
				return err
			}
		}
		err := utils.SetDatamodelDockerRepo(dockerRepo)
		if err != nil {
			return err
		}
	}
	fmt.Printf("\u2713 Datamodel %s initialized successfully\n", dmName)
	return nil
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a Nexus Datamodel",
	//Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: InitOperation,
}

func init() {
	name := "name"
	group := "group"
	InitCmd.Flags().StringVarP(&DatamodelName, name, "n", "", "name of the datamodel")
	InitCmd.Flags().StringVarP(&GroupName, group, "g", "", "subdomain for the datamodel resources")
	InitCmd.Flags().BoolVarP(&localDatamodel, localDatamodelFlag, "", false, "initializes a app local datamodel")
	InitCmd.Flags().StringVarP(&dockerRepo, "docker-repo", "d", "", "docker repo to publish image")
	InitCmd.Flags().StringVarP(&localDir, "local-dir", "l", "", "directory of the nexus repo")

	err := InitCmd.MarkFlagRequired(name)
	if err != nil {
		fmt.Printf("Failed to mark flag %s required: %v", name, err)
	}

	err = InitCmd.MarkFlagRequired(group)
	if err != nil {
		fmt.Printf("Failed to mark flag %s required: %v", group, err)
	}
}
