package gns

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	pb "gitlab.eng.vmware.com/nsx-allspark_users/go-protos/pkg/tsm-cli/cli"
	"gitlab.eng.vmware.com/nsx-allspark_users/lib-go/logging"
	"gitlab.eng.vmware.com/nexus/cli/pkg/utils"
	"google.golang.org/grpc"
)

// CreateCmd ...Creates a global namespace
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a global namespace object",
	//Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: create,
}

var createResourceFile string
var requiredResourceTags = []string{"spec", "metadata", "kind"}

func init() {
	CreateCmd.Flags().StringVarP(&createResourceFile, "file",
		"f", "", "Resource file from which global namespace is created.")

	err := cobra.MarkFlagRequired(CreateCmd.Flags(), "file")
	if err != nil {
		logging.Debugf("init error: %v", err)
	}
}

func create(cmd *cobra.Command, args []string) error {
	file, err := utils.ValidateAndGetJSON(createResourceFile, requiredResourceTags)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

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
		p := pb.File{}
		p.File = file
		_, err := c.Upsert(ctx, &pb.UpsertRequest{File: &p})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}()

	return nil
}
