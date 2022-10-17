bin_name=server

build:
	@echo "building binary..."
# go generate gives missing go sum entry for module errors
# https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/issues/1483
# you will need to first do a go get -u gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git
	go run -mod=mod gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git generate .
	go build -o ${bin_name} server.go
