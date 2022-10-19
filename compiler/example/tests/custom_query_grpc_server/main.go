package main

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/tests/custom_query_grpc_server/server"
)

func main() {
	server.StartQueryGRPC(1122)
}
