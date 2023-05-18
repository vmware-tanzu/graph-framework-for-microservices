package model

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"gitlab.eng.vmware.com/nsx-allspark_users/go-protos/pkg/health"
	authnexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/authentication.nexus.vmware.com/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/labstack/gommon/log"
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
	reg_svc "gitlab.eng.vmware.com/nsx-allspark_users/go-protos/pkg/registration-service/global"
	middlewarenexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/domain.nexus.vmware.com/v1"
	tenantv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/tenantconfig.nexus.vmware.com/v1"
)

// adding this global variables for CORS to support multiple domain and header configuration
var CorsConfigOrigins = map[string][]string{}
var CorsConfigHeaders = map[string][]string{}

type EventType string

const (
	Upsert EventType = "Upsert"
	Delete EventType = "Delete"
)

type NexusAnnotation struct {
	Name                 string                     `json:"name,omitempty"`
	Hierarchy            []string                   `json:"hierarchy,omitempty"`
	Children             map[string]NodeHelperChild `json:"children,omitempty"`
	Links                map[string]NodeHelperChild `json:"links,omitempty"`
	NexusRestAPIGen      nexus.RestAPISpec          `json:"nexus-rest-api-gen,omitempty"`
	NexusRestAPIMappings map[string]string          `json:"nexus-rest-api-mappings,omitempty"`
	IsSingleton          bool                       `json:"is_singleton,omitempty"`
	Description          string                     `json:"description,omitempty"`
}

type NodeHelperChild struct {
	FieldName    string `json:"fieldName"`
	FieldNameGvk string `json:"fieldNameGvk"`
	IsNamed      bool   `json:"isNamed"`
}

type NodeInfo struct {
	Name            string
	ParentHierarchy []string
	Children        map[string]NodeHelperChild
	Links           map[string]NodeHelperChild
	IsSingleton     bool
	Description     string
}

type RestUriInfo struct {
	TypeOfURI URIType
}

type URIType int

const (
	DefaultURI URIType = iota
	SingleLinkURI
	NamedLinkURI
	StatusURI
)

func ConstructEchoPathParamURL(uri string) string {
	replacer := strings.NewReplacer("{", ":", "}", "")
	return replacer.Replace(uri)
}

type OidcNodeEvent struct {
	Oidc authnexusv1.OIDC
	Type EventType
}

type TenantNodeEvent struct {
	Tenant    tenantv1.Tenant
	Type      EventType
	RegClient reg_svc.GlobalRegistrationClient
}

type DatamodelInfo struct {
	Title string
}

type CorsNodeEvent struct {
	Cors middlewarenexusv1.CORSConfig
	Type EventType
}

// LinkGvk : This model used to carry fully qualified object <gvk> and
// hierarchy information.
type LinkGvk struct {
	Group     string   `json:"group,omitempty" yaml:"group,omitempty"`
	Kind      string   `json:"kind,omitempty" yaml:"kind,omitempty"`
	Name      string   `json:"name,omitempty" yaml:"name,omitempty"`
	Hierarchy []string `json:"hierarchy,omitempty" yaml:"hierarchy,omitempty"`
}

type ConnectorObject struct {
	Service    string
	Protocol   string
	Connection *grpc.ClientConn
}

func (v *ConnectorObject) GetVersion() (interface{}, error) {
	var result interface{}
	var err error
	if v.Protocol == "http" {
		resp, err := http.Get(v.Service)
		if err != nil {
			return nil, err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, err
		}
		return result, nil
	} else if v.Protocol == "grpc" {
		if v.Connection == nil {
			err = v.InitConnection()
			if err != nil {
				return nil, err
			}
		}
		client := health.NewPingClient(v.Connection)
		result, err = client.Echo(context.Background(), &health.PingInput{Sequence: 1, Msg: ""})
		if err != nil {
			log.Errorf("Error in grpc call to version , errMsg: %s", err)
			return nil, err
		}
		return result, nil
	}
	return result, nil
}

func (v *ConnectorObject) InitConnection() (err error) {
	var Connection *grpc.ClientConn
	if v.Protocol == "grpc" {
		var registration_retry int = 0
		for registration_retry < 10 {
			registration_retry = registration_retry + 1
			Connection, err = grpc.DialContext(context.TODO(), v.Service, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				if registration_retry == 10 {
					return err
				} else {
					time.Sleep(10)
					continue
				}
			}
			break
		}
	}
	v.Connection = Connection
	return nil
}
