package model_test

import (
	"api-gw/pkg/model"
	"net/http"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ = Describe("Construct tests", func() {

	It("should construct new datamodel (vmware-test.org)", func() {
		unstructuredObj := unstructured.Unstructured{
			Object: map[string]interface{}{
				"spec": map[string]interface{}{
					"title": "VMWare Datamodel",
				},
			},
		}

		model.ConstructDatamodel(model.Upsert, "vmware-test.org", &unstructuredObj)
		Expect(model.DatamodelToDatamodelInfo).To(HaveKey("vmware-test.org"))
	})

	It("should delete datamodel vmware-test.org", func() {
		unstructuredObj := unstructured.Unstructured{
			Object: map[string]interface{}{
				"spec": map[string]interface{}{
					"title": "VMWare Datamodel",
				},
			},
		}

		model.ConstructDatamodel(model.Delete, "vmware-test.org", &unstructuredObj)
		Expect(model.DatamodelToDatamodelInfo).ToNot(HaveKey("vmware-test.org"))
	})
	It("should verify initConnection works", func() {
		connector := model.ConnectorObject{
			Service:  "http://localhost:80/version",
			Protocol: "http",
		}
		Expect(connector.InitConnection()).To(BeNil())

		grpcConnector := model.ConnectorObject{
			Service:  "localhost:3000",
			Protocol: "grpc",
		}
		Expect(grpcConnector.InitConnection()).To(BeNil())

		_, err := connector.GetVersion()
		Expect(err).NotTo(BeNil())

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		//Mock test for http endpoint
		httpmock.RegisterResponder("GET", "http://localhost:80/version", func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{"version": "test"})
			if err != nil {
				return &http.Response{}, err
			}
			return resp, nil
		})

		version, err := connector.GetVersion()
		Expect(err).To(BeNil())
		Expect(version).To(BeEquivalentTo(map[string]interface{}{"version": "test"}))

		_, err = grpcConnector.GetVersion()
		Expect(err).NotTo(BeNil())

	})
})
