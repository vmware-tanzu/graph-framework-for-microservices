package graphqlcalls

import (
	"context"
	"log"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-calibration/gqlclient"
)

func (g *GraphqlFuncs) GetManagers() {
	ctx := context.Background()
	_, err := gqlclient.Managers(ctx, g.Gclient)
	if err != nil {
		log.Printf("Failed to build request %v", err)
	}
}

func (g *GraphqlFuncs) GetEmployeeRole() {
	ctx := context.Background()
	_, err := gqlclient.Employees(ctx, g.Gclient)
	if err != nil {
		log.Printf("Failed to build request %v", err)
	}
}
