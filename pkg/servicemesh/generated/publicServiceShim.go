package generated

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type GeneratedGnsPublicService struct{}

type PublicServiceConfigProperties struct {
	Fqdn                  string       `mapstructure:"fqdn" json:"fqdn"`
	Name                  string       `mapstructure:"name" json:"name"`
	ExternalPort          int          `mapstructure:"external_port" json:"external_port"`
	ExternalProtocol      string       `mapstructure:"external_protocol" json:"external_protocol"`
	TTL                   int          `mapstructure:"ttl" json:"ttl"`
	PublicDomain          PublicDomain `mapstructure:"public_domain" json:"public_domain"`
	HAPolicy              string       `mapstructure:"ha_policy" json:"ha_policy"`
	GSLB                  GSLB         `mapstructure:"gslb" json:"gslb" `
	WildcardCertificateId string       `mapstructure:"wildcard_certificate_id" json:"wildcard_certificate_id"`
	HealthcheckIds        []string     `mapstructure:"healthcheck_ids" json:"healthcheck_ids"`
}

type PublicDomain struct {
	ExternalDNSId string `mapstructure:"external_dns_id" json:"external_dns_id"`
	PrimaryDomain string `mapstructure:"primary_domain" json:"primary_domain"`
	SubDomain     string `mapstructure:"sub_domain" json:"sub_domain"`
	CertificateId string `mapstructure:"certificate_id" json:"certificate_id"`
}
type GSLB struct {
	Type string `mapstructure:"type" json:"type"`
}

func (g GeneratedGnsPublicService) GetUrl(input map[string]string) (string, error) {
	if val, ok := input["project"]; ok {
		base := fmt.Sprintf("/v1alpha2/projects/%s/global-namespaces", val)
		if val, ok := input["global-namespace"]; ok {
			base = fmt.Sprintf("%s/%s", base, val)
			if val, ok := input["public-service"]; ok {
				base = fmt.Sprintf("%s/public-service/%s", base, val)
			}

		}
		return base, nil
	} else {
		return "", errors.New("project not found in input args")
	}
}

func (g GeneratedGnsPublicService) GetBody(input interface{}) ([]byte, error) {
	fmt.Printf("GNS Prop Input: %+v \n", input)
	gnsPSProp := PublicServiceConfigProperties{}
	mapstructure.Decode(input, &gnsPSProp)
	fmt.Printf("GNS Prop Output: %+v \n", gnsPSProp)
	return json.Marshal(gnsPSProp)
}

func init() {
	GeneratedMap["project:global-namespace:public-service"] = GeneratedGnsPublicService{}
}
