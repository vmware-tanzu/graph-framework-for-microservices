package debug

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/common"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/log"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/auth"
)

var IsDatamodelObjs bool

func Debug(cmd *cobra.Command, args []string) error {

	serverInfo, err := auth.ServerInfo()
	if err != nil {
		log.Errorf("Get serverInfo failed with error %v", err)
		return err
	}
	log.Debugf("ServerInfo: %v", serverInfo)

	if IsDatamodelObjs {
		err = DumpDatamodelObjects(serverInfo.Name)
		if err != nil {
			log.Errorf("Error while dumping all the datamodel objects: %v", err)
			return err
		}
	} else { // collect all the debug info if no option is given
		err = DumpDatamodelObjects(serverInfo.Name)
		if err != nil {
			log.Errorf("Error while dumping all the datamodel objects: %v", err)
			return err
		}
	}
	return nil
}

func DumpDatamodelObjects(serverName string) error {
	scriptFile, err := common.ScriptFs.Open("dump_datamodel_objects.sh")
	if err != nil {
		return fmt.Errorf("error while reading debug scriptFile %v", err)

	}

	c := exec.Command("/bin/bash", "/dev/stdin", serverName)
	c.Stdin = scriptFile

	resp, err := c.Output()
	if err != nil {
		return err
	}

	fmt.Println(string(resp))

	return nil
}
