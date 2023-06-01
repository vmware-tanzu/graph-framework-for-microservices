package csptenant_test

import (
	csptenant "api-gw/internal/tenant/csp"
	"api-gw/pkg/common"
	"fmt"
	"net/http"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test Mock CSP APIs", func() {

	It("should test CSP tenant interaction methods", func() {
		cspTenant := csptenant.InitCSPTenant(
			"testToken",
			"44c80b09-3a36-453d-86db-527fbed2917d",
			"http://localhost",
		)

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", fmt.Sprintf("http://localhost/%s/%s", common.CSP_GATEWAY_ROOT, common.CSP_AUTHORIZE_URL), func(r *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{"access_token": "testToken", "token_type": "Token", "expires_in": 100})
			if err != nil {
				return &http.Response{}, err
			}
			return resp, nil
		})

		_, err := cspTenant.Token("test")
		Expect(err).To(BeNil())

		httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost/%s/%s", common.CSP_GATEWAY_ROOT, common.CSP_COMMERCE_API), func(r *http.Request) (*http.Response, error) {
			resultJson := csptenant.OrderSubscriptionResults{
				Results: []csptenant.OrderSubscription{{
					Offers: []csptenant.OrderOffer{
						{
							ProductID: "NSM-XX-CP",
							Status:    "ACTIVE",
						},
					},
				},
				},
			}
			resp, err := httpmock.NewJsonResponse(200, resultJson)
			if err != nil {
				fmt.Println(err)
				return &http.Response{}, err
			}
			return resp, nil
		})
		sku, err := cspTenant.ProductID("Test")
		fmt.Println(err)
		Expect(err).To(BeNil())
		Expect(sku).To(Equal("NSM-XX-CP"))

		httpmock.RegisterResponder("GET", "http://localhost/csp/gateway/commerce/api/v3/subscriptions", func(r *http.Request) (*http.Response, error) {
			resultJson := csptenant.OrderSubscriptionResults{
				Results: []csptenant.OrderSubscription{{
					Offers: []csptenant.OrderOffer{
						{
							ProductID: "NSM-XX-CP",
							Status:    "NOT_ACTIVe",
						},
					},
				},
				},
			}
			resp, err := httpmock.NewJsonResponse(200, resultJson)
			if err != nil {
				fmt.Println(err)
				return &http.Response{}, err
			}
			return resp, nil
		})
		sku, err = cspTenant.ProductID("Test")
		fmt.Println(err)
		Expect(err).To(BeNil())
		Expect(sku).To(Equal(""))

		httpmock.RegisterResponder("GET", "http://localhost/csp/gateway/commerce/api/v3/subscriptions", func(r *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(403, "Invalid token")
			return resp, fmt.Errorf("invalid token")
		})
		sku, err = cspTenant.ProductID("Test")
		fmt.Println(err)
		Expect(err).NotTo(BeNil())
		Expect(sku).To(Equal(""))

	})

})
