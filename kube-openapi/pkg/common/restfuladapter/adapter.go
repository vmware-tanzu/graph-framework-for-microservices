package restfuladapter

import (
	"github.com/emicklei/go-restful"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kube-openapi.git/pkg/common"
)

// AdaptWebServices adapts a slice of restful.WebService into the common interfaces.
func AdaptWebServices(webServices []*restful.WebService) []common.RouteContainer {
	var containers []common.RouteContainer
	for _, ws := range webServices {
		containers = append(containers, &WebServiceAdapter{ws})
	}
	return containers
}
