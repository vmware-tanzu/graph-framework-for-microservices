package main

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/tests/custom_query_grpc_server/server"
)

func main() {
	server.StartQueryGRPC(1122)
}
