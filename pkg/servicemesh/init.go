package servicemesh

import (
	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/log"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/app"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/apply"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/config"
	servicemesh_datamodel "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/datamodel"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/login"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/operator"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/runtime"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/version"
)

// ApplyCmd ... Apply command
var ApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply Servicemesh configuration from file",
	RunE:  apply.ApplyResource,
}

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Servicemesh configuration from file",
	RunE:  apply.DeleteResource,
}

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to csp",
	RunE:  login.Login,
}

var RuntimeCmd = &cobra.Command{
	Use:   "runtime",
	Short: "Perform Nexus Runtime operations",
}

var DataModelCmd = &cobra.Command{
	Use:   "datamodel",
	Short: "Perform Nexus Datamodel operations",
}

var AppCmd = &cobra.Command{
	Use:   "app",
	Short: "Perform Nexus Application operations",
}

var OperatorCmd = &cobra.Command{
	Use:   "operator",
	Short: "Create, update or delete operators within Nexus apps",
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Provides Nexus CLI, compiler, app-template and runtime versions",
	RunE:  version.Version,
}

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "set nexus CLI preferences",
}

func initCommands() {
	ApplyCmd.Flags().StringVarP(&apply.CreateResourceFile, "file",
		"f", "", "Resource file from which cluster is created.")

	err := cobra.MarkFlagRequired(ApplyCmd.Flags(), "file")
	if err != nil {
		log.Debugf("init error: %v", err)
	}

	DeleteCmd.Flags().StringVarP(&apply.DeleteResourceFile, "file",
		"f", "", "Resource file from which cluster is created.")

	err = cobra.MarkFlagRequired(DeleteCmd.Flags(), "file")
	if err != nil {
		log.Debugf("init error: %v", err)
	}

	LoginCmd.Flags().StringVarP(&login.ApiToken, "token",
		"t", "", "token for api access")

	err = cobra.MarkFlagRequired(LoginCmd.Flags(), "token")
	if err != nil {
		log.Debugf("api token is mandatory for login")
	}

	LoginCmd.Flags().StringVarP(&login.Server, "server",
		"s", "", "saas server fqdn")

	err = cobra.MarkFlagRequired(LoginCmd.Flags(), "server")
	if err != nil {
		log.Debugf("saas server fqdn name is mandatory for login")
	}
	RuntimeCmd.AddCommand(runtime.InstallCmd)
	RuntimeCmd.AddCommand(runtime.UninstallCmd)

	DataModelCmd.AddCommand(servicemesh_datamodel.InitCmd)
	DataModelCmd.AddCommand(servicemesh_datamodel.InstallCmd)
	DataModelCmd.AddCommand(servicemesh_datamodel.BuildCmd)
	DataModelCmd.AddCommand(servicemesh_datamodel.ConfigureCmd)

	AppCmd.AddCommand(app.InitCmd)
	AppCmd.AddCommand(app.PackageCmd)
	AppCmd.AddCommand(app.PublishCmd)
	AppCmd.AddCommand(app.DeployCmd)
	AppCmd.AddCommand(app.RunCmd)
	AppCmd.AddCommand(app.AddDatamodelCmd)
	AppCmd.AddCommand(app.SetDefaultDatamodelCmd)
	AppCmd.AddCommand(app.BuildCmd)

	OperatorCmd.AddCommand(operator.CreateCmd)

	ConfigCmd.AddCommand(config.SetCmd)

	ConfigCmd.AddCommand(config.ViewCmd)
}

func init() {
	initCommands()
}
