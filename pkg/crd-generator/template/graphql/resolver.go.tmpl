package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"context"
	"fmt"
	"sync"

	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/generated/graphql"
	"google.golang.org/grpc"
	"nexustempmodule/nexus-gql/graph/model"
)

type Resolver struct {
	CustomQueryHandler
}

type CustomQueryHandler struct {
	mtx     sync.Mutex
	Clients map[string]graphql.ServerClient
}

func (c *CustomQueryHandler) Connect(endpoint string) (graphql.ServerClient, error) {
	conn, err := grpc.Dial(endpoint)
	if err != nil {
		return nil, err
	}
	cl := graphql.NewServerClient(conn)
	c.mtx.Lock()
	c.Clients[endpoint] = cl
	c.mtx.Unlock()
	return cl, nil
}

func (c *CustomQueryHandler) GetClient(endpoint string) (graphql.ServerClient, error) {
	c.mtx.Lock()
	cl, ok := c.Clients[endpoint]
	c.mtx.Unlock()
	if ok {
		return cl, nil
	}
	return c.Connect(endpoint)
}

func (c *CustomQueryHandler) Query(endpoint string, query *graphql.GraphQLQuery) (*model.NexusGraphqlResponse, error) {
	client, err := c.GetClient(endpoint)
	if err != nil {
		return nil, err
	}
	resp, err := client.Query(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	return grpcResToGraphQl(resp), nil
}

func grpcResToGraphQl(response *graphql.GraphQLResponse) *model.NexusGraphqlResponse {
	if response == nil {
		return nil
	}
	dataStr, err := json.Marshal(response.Data)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	return &model.NexusGraphqlResponse{
		Code:         intToPointer(int(response.Code)),
		TotalRecords: intToPointer(int(response.TotalRecords)),
		Message:      &response.Message,
		Last:         &response.Last,
		Data:         stringToPointer(string(dataStr)),
	}
}

func intToPointer(i int) *int {
	return &i
}

func stringToPointer(i string) *string {
	return &i
}
