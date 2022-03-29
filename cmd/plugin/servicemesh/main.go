package main

import (
	"os"

	"github.com/vmware-tanzu-private/core/pkg/v1/cli"
	"github.com/vmware-tanzu-private/core/pkg/v1/cli/command/plugin"
	"gitlab.eng.vmware.com/nsx-allspark_users/lib-go/logging"
	"gitlab.eng.vmware.com/nexus/cli/pkg/servicemesh"
)

var descriptor = cli.PluginDescriptor{
	Name:        "servicemesh",
	Description: "service mesh features",
	Version:     "v0.0.1",
	BuildSHA:    "",
	Group:       cli.ManageCmdGroup,
	DocURL:      "",
}

func main() {
	p, err := plugin.NewPlugin(&descriptor)
	if err != nil {
		logging.Fatalf("Plugin failed to load %v", err.Error())
	}

	p.AddCommands(
		servicemesh.ClusterCmd,
		servicemesh.GnsCmd,
		servicemesh.ConfigCmd,
		servicemesh.ApplyCmd,
		servicemesh.DeleteCmd,
		servicemesh.LoginCmd,
	)

	if err := p.Execute(); err != nil {
		os.Exit(1)
	}
}
