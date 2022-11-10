package ignored

import "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel/config/gns"

type Config struct {
	GNS gns.Gns `nexus:"child"`
}
