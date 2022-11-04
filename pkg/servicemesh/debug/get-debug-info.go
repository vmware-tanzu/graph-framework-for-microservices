package debug

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/log"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/auth"
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
