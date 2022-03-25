package gns

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
	gnsKind = "GNS"
)

// GetCmd ...Gets a global namespace identified by name
var GetCmd = &cobra.Command{
	Use:   "get GlobalNamespace_NAME",
	Short: "Get a global namespace",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: get,
}

type getGnsOptions struct {
	gnsName      string
	outputFormat string
}

var cd = getGnsOptions{}

func init() {
	GetCmd.Flags().StringVarP(&cd.outputFormat, "output", "o", "", "Output formart. Supported formats: json|yaml")
}

func get(cmd *cobra.Command, args []string) error {

	cd.gnsName = args[0]

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
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second) // TODO check into desired timeout
		defer cancel()

		p := pb.Path{}
		p.Kind = gnsKind
		p.Id = cd.gnsName
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

		resourceNotFoundErr := fmt.Errorf("error: the server doesn't have a resource type: %s", cd.gnsName)
		fmt.Println(resourceNotFoundErr.Error())

	}()

	return nil
}
