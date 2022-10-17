package models

import "github.com/vmware-tanzu/graph-framework-for-microservices/src/gqlgen/integration/remote_api"

type Viewer struct {
	User *remote_api.User
}
