package cluster

import (
	"context"
	"fmt"
	"time"

	//"github.com/ghodss/yaml"

	"github.com/spf13/cobra"
	pb "gitlab.eng.vmware.com/nsx-allspark_users/go-protos/pkg/tsm-cli/cli"
	"gitlab.eng.vmware.com/nexus/cli/pkg/utils"
	"google.golang.org/grpc"
)

// ListCmd ...Lists all the clusters
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List clusters",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: list,
}

type listClusterOptions struct {
	outputFormat string
}

var lc = listClusterOptions{}

func init() {
	ListCmd.Flags().StringVarP(&lc.outputFormat, "output", "o", "", "Output formart. Supported formats: json|yaml")
}

func list(cmd *cobra.Command, args []string) error {

	conn, err := grpc.Dial(utils.GetTestURL(), grpc.WithInsecure())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	defer conn.Close()
	c := pb.NewCliClient(conn)

	// Contact Tenant API GRPC Server
	func() {
		//log.Info("Connecting: ", conn.Target())
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()

		p := pb.Path{}
		p.Kind = clusterKind
		p.Id = "*"
		p.Path = nil

		resp, err := c.List(ctx, &pb.ListRequest{Path: &p})
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		paths := resp.GetResponse()
		if lc.outputFormat == "" {
			for _, path := range paths {
				fmt.Println(path.Id)
			}
		}

		// Render as List Kind TBD
	}()

	return nil
}
