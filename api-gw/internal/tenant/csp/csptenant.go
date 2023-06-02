package csptenant

import (
	"api-gw/pkg/common"
	"encoding/json"
	"fmt"
	"time"

	"github.com/levigross/grequests"
	"github.com/pkg/errors"
)

type CSPTenant struct {
	serviceOwnerToken string
	serviceId         string // ServiceDefinitionID for this server

	server string

	token             *ApiToken
	lastAuthorizeTime time.Time
}

type ApiToken struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	AccessToken string `json:"access_token"`
}

const (
	OfferStatusActive              = "ACTIVE"
	OfferStatusPendingProvisioning = "PENDING_PROVISIONING"
)

type OrderOffer struct {
	ProductID string `json:"productId,omitempty"`
	Status    string `json:"status" validate:"required"`
}

type OrderSubscription struct {
	Offers []OrderOffer `json:"offers"`
}

type OrderSubscriptionResults struct {
	Results []OrderSubscription `json:"results"`
}

func Time(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

var cspTenant CSPTenant

func InitCSPTenant(serviceOwnerToken, serviceId, server string) *CSPTenant {
	if cspTenant == (CSPTenant{}) {
		cspTenant = CSPTenant{
			serviceOwnerToken: serviceOwnerToken,
			serviceId:         serviceId,
			server:            fmt.Sprintf("%s/%s", server, common.CSP_GATEWAY_ROOT),
		}
	}
	return &cspTenant
}

func (c *CSPTenant) url(url string) string {
	return fmt.Sprintf("%s/%s", c.server, url)
}

func (c *CSPTenant) authorize(tenantId string) error {
	options := grequests.RequestOptions{
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Basic %s", c.serviceOwnerToken),
		},
		Data: map[string]string{
			"grant_type": "client_credentials",
		},
	}

	resp, err := grequests.Post(c.url(common.CSP_AUTHORIZE_URL), &options)
	if err != nil {
		return errors.Wrap(err, "authorization error")
	}
	if !resp.Ok {
		return errors.New(fmt.Sprintf("authorization error %s", resp.String()))
	}
	if err := resp.JSON(&c.token); err != nil {
		return errors.Wrap(err, "authorization JSON error")
	}
	c.lastAuthorizeTime = time.Now()
	return nil
}

func (c *CSPTenant) expiration() int {
	return c.token.ExpiresIn - (c.token.ExpiresIn * 100 / 20)
}

func (c *CSPTenant) IsAuthorized() bool {
	return c.token != nil && c.lastAuthorizeTime.Add(time.Second*time.Duration(c.expiration())).After(time.Now())
}

func (c *CSPTenant) Token(tenantId string) (*ApiToken, error) {
	if c.IsAuthorized() {
		return c.token, nil
	}
	if err := c.authorize(tenantId); err != nil {
		return nil, err
	}
	return c.token, nil
}

func (c *CSPTenant) ProductID(tenantId string) (string, error) {

	// Get access token for the service owner.
	token, tokenErr := c.Token(tenantId)
	if tokenErr != nil {
		return "", errors.Wrap(tokenErr, fmt.Sprintf("unable to authorize for serviceDefinitionId %s and orgId %s",
			c.serviceId, tenantId))
	}

	options := &grequests.RequestOptions{}
	if options.Headers == nil {
		options.Headers = make(map[string]string)
	}
	options.Headers["csp-auth-token"] = token.AccessToken

	if options.Params == nil {
		options.Params = make(map[string]string)
	}
	options.Params["serviceDefinitionId"] = c.serviceId
	options.Params["orgId"] = tenantId

	resp, err := grequests.Get(c.url(common.CSP_COMMERCE_API), options)
	if err != nil || !resp.Ok {
		return "", errors.Wrap(err, fmt.Sprintf("getting subscriptions for serviceDefinitionId %s and orgId %s failed",
			c.serviceId, tenantId))
	}

	raw := resp.Bytes()
	var results OrderSubscriptionResults
	if err := json.Unmarshal(raw, &results); err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("unmarshalling subscriptions for serviceDefinitionId %s and orgId %s failed",
			c.serviceId, tenantId))
	}

	for _, subscription := range results.Results {
		for _, offer := range subscription.Offers {
			if offer.Status == OfferStatusActive {
				return offer.ProductID, nil
			}
		}
	}

	return "", errors.Wrap(err, fmt.Sprintf("unable to determine subscriptions for serviceDefinitionId %s and orgId %s",
		c.serviceId, tenantId))
}
