package models

import "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/integration/remote_api"

type Viewer struct {
	User *remote_api.User
}
