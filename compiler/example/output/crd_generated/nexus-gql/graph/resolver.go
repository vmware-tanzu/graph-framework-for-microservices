package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/golang/protobuf/proto"
	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/generated/graphql"
	qm "golang-appnet.eng.vmware.com/nexus-sdk/nexus/generated/query-manager"
	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/crd_generated/nexus-gql/graph/model"
)

type Resolver struct{}

type GrpcClient interface {
	Request(query proto.Message) (interface{}, error)
}

type GrpcClients struct {
	mtx     sync.Mutex
	Clients map[string]GrpcClient
}

func (s *GrpcClients) Request(endpoint string, apiType nexus.GraphQlApiType, query proto.Message) (interface{}, error) {
	cl, err := s.getClient(endpoint, apiType)
	if err != nil {
		return nil, err
	}
	return cl.Request(query)
}

func (s *GrpcClients) addClient(endpoint string, apiType nexus.GraphQlApiType) (GrpcClient, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	switch apiType {
	case nexus.GraphQLQueryApi:
		cl := &NexusQueryClient{graphql.NewServerClient(conn)}
		s.Clients[fmt.Sprint(apiType)+"/"+endpoint] = cl
		return cl, nil
	case nexus.GetMetricsApi:
		cl := &QmClient{qm.NewServerClient(conn)}
		s.Clients[fmt.Sprint(apiType)+"/"+endpoint] = cl
		return cl, nil
	default:
		return nil, fmt.Errorf("unsupported GraphQlApiType %v", apiType)
	}
}

func (s *GrpcClients) getClient(endpoint string, apiType nexus.GraphQlApiType) (GrpcClient, error) {
	s.mtx.Lock()
	cl, ok := s.Clients[endpoint]
	s.mtx.Unlock()
	if ok {
		return cl, nil
	}
	return s.addClient(endpoint, apiType)
}

type NexusQueryClient struct {
	graphql.ServerClient
}

func (c *NexusQueryClient) Request(query proto.Message) (interface{}, error) {
	q, ok := query.(*graphql.GraphQLQuery)
	if !ok {
		return nil, fmt.Errorf("wrong format of query used for nexus query")
	}
	resp, err := c.Query(context.TODO(), q)
	if err != nil {
		return nil, err
	}
	return grpcResToGraphQl(resp)
}

func grpcResToGraphQl(response *graphql.GraphQLResponse) (*model.NexusGraphqlResponse, error) {
	if response == nil {
		return nil, nil
	}
	dataStr, err := json.Marshal(response.Data)
	if err != nil {
		return nil, err
	}
	return &model.NexusGraphqlResponse{
		Code:         intToPointer(int(response.Code)),
		TotalRecords: intToPointer(int(response.TotalRecords)),
		Message:      &response.Message,
		Last:         &response.Last,
		Data:         stringToPointer(string(dataStr)),
	}, nil
}

type QmClient struct {
	qm.ServerClient
}

func (c *QmClient) Request(query proto.Message) (interface{}, error) {
	q, ok := query.(*qm.MetricArg)
	if !ok {
		return nil, fmt.Errorf("wrong format of query used for metrics query")
	}
	resp, err := c.GetMetrics(context.TODO(), q)
	if err != nil {
		return nil, err
	}
	return timeSeriesResToTimeSeriesData(resp)
}

func timeSeriesResToTimeSeriesData(response *qm.TimeSeriesResponse) (*model.TimeSeriesData, error) {
	if response == nil {
		return nil, nil
	}
	dataStr, err := json.Marshal(response.Data)
	if err != nil {
		return nil, err
	}
	return &model.TimeSeriesData{
		Code:         intToPointer(int(response.Code)),
		TotalRecords: intToPointer(int(response.TotalRecords)),
		Message:      &response.Message,
		Last:         &response.Last,
		Data:         stringToPointer(string(dataStr)),
	}, nil
}

func intToPointer(i int) *int {
	return &i
}

func stringToPointer(i string) *string {
	return &i
}

func pointerToString[T any, Ptr *T](any Ptr) string {
	if any == nil {
		return ""
	}
	return fmt.Sprintf("%v", *any)
}
