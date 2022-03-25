package cluster

import (
	//"errors"

	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	pb "gitlab.eng.vmware.com/nsx-allspark_users/go-protos/pkg/tsm-cli/cli"
	"gitlab.eng.vmware.com/nexus/cli/pkg/utils"
	"google.golang.org/grpc"
)

const (
	clusterKind = "Cluster"
)

// GetCmd ...Gets a cluster identified by name
var GetCmd = &cobra.Command{
	Use:   "get CLUSTER_NAME",
	Short: "Get a cluster",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: get,
}

type getClusterOptions struct {
	clusterName  string
	outputFormat string
}

var cd = getClusterOptions{}

func init() {
	GetCmd.Flags().StringVarP(&cd.outputFormat, "output", "o", "", "Output formart. Supported formats: json|yaml")
}

func get(cmd *cobra.Command, args []string) error {

	cd.clusterName = args[0]

	conn, err := grpc.Dial(utils.GetTestURL(), grpc.WithInsecure())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	defer conn.Close()
	c := pb.NewCliClient(conn)

	func() {
		//log.Info("Connecting: ", conn.Target())
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()

		p := pb.Path{}
		p.Kind = clusterKind
		p.Id = cd.clusterName
		p.Path = nil

		r, err := c.Get(ctx, &pb.GetRequest{Path: &p})
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		payload := r.File.GetFile()

		if len(payload) > 0 {
			out := utils.RenderOutput(payload, cd.outputFormat)
			fmt.Printf("%s", out)
			return
		}

		resourceNotFoundErr := fmt.Errorf("error: the server doesn't have a resource type: %s", cd.clusterName)
		fmt.Println(resourceNotFoundErr.Error())

	}()

	return nil
}
