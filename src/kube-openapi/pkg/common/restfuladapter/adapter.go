package restfuladapter

import (
	"github.com/emicklei/go-restful"
	"github.com/vmware-tanzu/graph-framework-for-microservices/src/kube-openapi/pkg/common"
)

// AdaptWebServices adapts a slice of restful.WebService into the common interfaces.
func AdaptWebServices(webServices []*restful.WebService) []common.RouteContainer {
	var containers []common.RouteContainer
	for _, ws := range webServices {
		containers = append(containers, &WebServiceAdapter{ws})
	}
	return containers
}
