package runtime

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli/pkg/utils"
)

func Uninstall(cmd *cobra.Command, args []string) error {
	if Namespace == "" {
		Namespace = "default"
	}
	clientset, _, _, _, _, err := utils.GenerateContext("")
	if err != nil {
		return err
	}
	fmt.Printf("deleting tenant %s from current cluster", Namespace)
	err = clientset.CoreV1().Namespaces().Delete(context.Background(), Namespace, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

var UninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "uninstall a tenant from current cluster",
	//Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: Uninstall,
}

func init() {
	UninstallCmd.Flags().StringVarP(&Namespace, "namespace",
		"t", "", "name of the namespace to be created")
}
