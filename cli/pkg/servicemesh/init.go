package servicemesh

import (
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/log"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/app"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/apply"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/config"
	servicemesh_datamodel "github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/datamodel"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/debug"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/login"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/operator"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/runtime"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/version"
)

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Declarative configuration from file or with type",
	Args:  cobra.RangeArgs(0, 3),
	RunE:  apply.GetResource,
}

var GetSpecCmd = &cobra.Command{
	Use:   "spec [short/long crd name]",
	Short: "Get YAML spec for given object",
	Args:  cobra.RangeArgs(1, 1),
	RunE:  apply.GetSpec,
}

var ApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply Declarative configuration from file",
	RunE:  apply.ApplyResource,
}

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Declarative configuration from file",
	RunE:  apply.DeleteResource,
}

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Saas",
	PreRun: func(cmd *cobra.Command, args []string) {
		insecure, _ := cmd.Flags().GetBool("in-secure")
		if !insecure {
			cmd.MarkFlagRequired("token")
		}
	},
	RunE: login.Login,
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

var DebugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Collect required debug info",
	RunE:  debug.Debug,
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

	GetCmd.Flags().StringVarP(&apply.GetResourceFile, "file",
		"f", "", "Resource file from which cluster is fetched.")
	apply.DefaultGetHelpFunc = GetCmd.HelpFunc()
	GetCmd.SetHelpFunc(apply.GetHelp)

	GetCmd.Flags().StringVarP(&apply.Labels, "labels",
		"l", "", "labels required for the resource to fetch.")

	LoginCmd.Flags().StringVarP(&login.ApiToken, "token",
		"t", "", "token for api access")

	LoginCmd.Flags().StringVarP(&login.Server, "server",
		"s", "", "saas server fqdn")

	err = cobra.MarkFlagRequired(LoginCmd.Flags(), "server")
	if err != nil {
		log.Debugf("saas server fqdn name is mandatory for login")
	}

	LoginCmd.Flags().BoolVarP(&login.IsPrivateSaas, "private-saas",
		"p", false, "private saas cluster")

	LoginCmd.Flags().BoolVarP(&login.IsInSecure, "in-secure",
		"k", false, "local/kind cluster")

	err = cobra.MarkFlagRequired(LoginCmd.Flags(), "server")
	if err != nil {
		log.Debugf("saas server fqdn name is mandatory for login")
	}

	DebugCmd.Flags().BoolVarP(&debug.IsDatamodelObjs, "datamodel-objs",
		"d", false, "dump all the datamodel crd and objects")

	GetCmd.AddCommand(GetSpecCmd)

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
