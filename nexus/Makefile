generate:
	protoc --go_out=generated --go-grpc_out=generated proto/**/*
	mockgen -source=generated/graphql/query_grpc.pb.go > generated/graphql/mock/mock.go

