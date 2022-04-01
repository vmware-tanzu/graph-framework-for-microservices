package utils

import (
	"fmt"
	"os"

	common "gitlab.eng.vmware.com/nexus/cli/pkg/common"
)

func GoToNexusDirectory() error {
	if _, err := os.Stat(common.NEXUS_DIR); os.IsNotExist(err) {
		return fmt.Errorf("%s directory not found", common.NEXUS_DIR)
	} else if err != nil {
		return fmt.Errorf("error %v trying to find directory %s", err, common.NEXUS_DIR)
	}

	if err := os.Chdir(common.NEXUS_DIR); err != nil {
		return fmt.Errorf("error %v trying to cd to directory %s", err, common.NEXUS_DIR)
	}
	return nil

}

func CheckDatamodelDirExists(datamodelName string) error {
	dmDir := datamodelName
	if _, err := os.Stat(dmDir); os.IsNotExist(err) {
		return fmt.Errorf("datamodel directory %s not found", dmDir)
	} else if err != nil {
		return fmt.Errorf("error %v trying to find datamodel directory %s", err, dmDir)
	}
	return nil
}

func StoreCurrentDatamodel(datamodelName string) error {
	f, err := os.OpenFile("NEXUSDATAMODEL", os.O_RDWR, os.ModeAppend)
	if err != nil {
		f, err = os.Create("NEXUSDATAMODEL")
		if err != nil {
			return err
		}
	}
	_, err = f.WriteString(datamodelName)
	if err != nil {
		fmt.Println("Could not store current datamodel name")
		return err
	}
	return nil
}
