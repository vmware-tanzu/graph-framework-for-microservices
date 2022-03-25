package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nexus/cli/pkg/utils"
)

// ViewCmd ... View Config
var ViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View Configuration",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: view,
}

type getConfigOptions struct {
	//serverName     string
	outputFormat   string
	allContext     bool
	currentContext string
}

var cd = getConfigOptions{}

func init() {
	ViewCmd.Flags().StringVarP(&cd.outputFormat, "output", "o", "", "Output formart. Supported formats: json|yaml")
	ViewCmd.Flags().StringVarP(&cd.currentContext, "context", "c", "", "View current config")
	ViewCmd.Flags().BoolVar(&cd.allContext, "all-contexts", true, "View all the contexts")
}

func view(cmd *cobra.Command, args []string) error {
	s, _ := utils.GetCurrentServer()
	fmt.Println("\nACCESS TOKEN:", s.GlobalOpts.Auth.AccessToken)
	fmt.Println("\nSAAS URL:", utils.GetSaasURL())
	fmt.Println("\nToken Expired:", utils.IsExpired(s.GlobalOpts.Auth.Expiration.Time))
	return nil
}
