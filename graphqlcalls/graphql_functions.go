package graphqlcalls

import (
	"context"
	"log"

	"github.com/Khan/genqlient/graphql"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-calibration/gqlclient"
)

func GetManagers(ctx context.Context, gclient graphql.Client) {
	_, err := gqlclient.Managers(ctx, gclient)
	if err != nil {
		log.Printf("Failed to build request %v", err)
	}
}

func GetEmployeeRole(ctx context.Context, gclient graphql.Client) {
	_, err := gqlclient.Employees(ctx, gclient)
	if err != nil {
		log.Printf("Failed to build request %v", err)
	}
}
