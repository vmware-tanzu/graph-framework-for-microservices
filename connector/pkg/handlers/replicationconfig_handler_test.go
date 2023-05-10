package handlers_test

import (
	"context"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	fake_dynamic "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/testing"

	"connector/pkg/config"
	"connector/pkg/handlers"
	"connector/pkg/utils"
)

var _ = Describe("ReplicationConfig Tests", func() {
	When("An ReplicationConfig CR is created", func() {
		var (
			gvr               schema.GroupVersionResource
			conf              *config.Config
			handler           *handlers.ReplicationConfigHandler
			replicationConfig *unstructured.Unstructured
			gvrToListKind     = make(map[schema.GroupVersionResource]string)
			server            *ghttp.Server
			client            *fake_dynamic.FakeDynamicClient
		)

		BeforeEach(func() {
			server = ghttp.NewServer()
			parts := strings.Split(server.Addr(), ":")
			os.Setenv(utils.RemoteEndpointHost, fmt.Sprintf("http://%s", parts[0]))
			os.Setenv(utils.RemoteEndpointPort, parts[1])

			endpointGvr := schema.GroupVersionResource{Group: "connect.nexus.vmware.com", Version: "v1", Resource: "nexusendpoints"}
			acGvr := schema.GroupVersionResource{Group: "config.mazinger.com", Version: "v1", Resource: "apicollaborationspaces"}
			adGvr := schema.GroupVersionResource{Group: "config.mazinger.com", Version: "v1", Resource: "apidevspaces"}
			gvrToListKind[acGvr] = "ApiCollaborationSpaceList"
			gvrToListKind[endpointGvr] = "NexusEndpointList"
			gvrToListKind[adGvr] = "ApiDevSpaceList"

			scheme := runtime.NewScheme()
			client = fake_dynamic.NewSimpleDynamicClientWithCustomListKinds(scheme, gvrToListKind)

			// Valid Nexus-Endpoint CR.
			nexusEndpoint := getNexusEndpointObject("default", fmt.Sprintf("http://%s", parts[0]), parts[1], "")

			// Invalid Nexus-Endpoint CR.
			inValidNexusEndpoint := &unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name": "defaultNew",
					},
					"spec": map[string]interface{}{
						"port": map[string]interface{}{
							"invalid": "invalid",
						},
					},
				},
			}

			// Wrong Type Nexus-Endpoint CR.
			wrongTypeNexusEndpoint := getNexusEndpointObject("wrongType", fmt.Sprintf("http://%s", parts[0]), math.Inf(1), "")

			// Wrong Cert Nexus-Endpoint CR.
			wrongCertNexusEndpoint := getNexusEndpointObject("wrongCert", fmt.Sprintf("http://%s", parts[0]), parts[1], "abc")

			for _, val := range []*unstructured.Unstructured{nexusEndpoint, inValidNexusEndpoint, wrongTypeNexusEndpoint, wrongCertNexusEndpoint} {
				_, err := client.Resource(endpointGvr).Create(context.TODO(), val, metav1.CreateOptions{})
				Expect(err).NotTo(HaveOccurred())
			}

			ac1 := &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "config.mazinger.com/v1",
					"kind":       "ApiCollaborationSpace",
					"metadata": map[string]interface{}{
						"name": "ac1",
						"labels": map[string]interface{}{
							Root:    "root",
							Project: "project",
							Config:  "config",
						},
					},
					"spec": map[string]interface{}{
						"example": "example",
					},
				},
			}
			_, err := client.Resource(acGvr).Create(context.TODO(), ac1, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())

			gvr = utils.GetGVRFromCrdType(utils.ReplicationConfigCRD, utils.V1Version)
			conf, err = config.LoadConfig("./../config/test_utils/correct.yaml")
			Expect(err).NotTo(HaveOccurred())

			conf.PeriodicSyncInterval = 3 * time.Second
			handler = handlers.NewReplicationConfigHandler(gvr, conf, client)
		})

		It("Should connect to the destination endpoint correctly and"+
			"replicate the existing objects.", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/apicollaborationspaces/ac1"),
					ghttp.RespondWith(200, "turbo: true"),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/apis/config.mazinger.com/v1/apicollaborationspaces"),
					ghttp.RespondWith(200, "{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"ApiCollaborationSpace\",\"metadata\":{\"labels\":{},\"name\":\"ac1\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}}"),
				),
			)
			replicationConfig = getReplicationConfigObject()
			err := handler.Create(replicationConfig)
			Expect(err).NotTo(HaveOccurred())

			Expect(server.ReceivedRequests()).Should(HaveLen(2))
		})

		It("Should fail object creation when CRD Type not found on the destination.", func() {
			client.Fake.PrependReactor("list", "apicollaborationspaces",
				func(action testing.Action) (bool, runtime.Object, error) {
					return true, nil, fmt.Errorf("nope")
				})

			replicationConfig = &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"source": map[string]interface{}{
							"name": "INVALID_OBJECT",
							"kind": "Object",
							"object": map[string]interface{}{
								"objectType": map[string]interface{}{
									"group":   Group,
									"version": "v1",
									"kind":    AcKind,
								},
							},
						},
						"destination": map[string]interface{}{
							"hierarchical": false,
						},
						"remoteEndpointGvk": map[string]interface{}{
							"group": "connect.nexus.vmware.com",
							"kind":  "NexusEndpoint",
							"name":  "default",
						},
					},
				},
			}
			err := handler.Create(replicationConfig)
			Expect(err.Error()).To(ContainSubstring("error replicating desired nodes"))
		})

		It("Should replicate the existing objects when source is hierarchical", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/apicollaborationspaces/ac1"),
					ghttp.RespondWith(200, "turbo: true"),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/apis/config.mazinger.com/v1/apicollaborationspaces"),
					ghttp.RespondWith(200, "{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"ApiCollaborationSpace\",\"metadata\":{\"labels\":{\"configs.config.mazinger.com\":\"config\",\"projects.config.mazinger.com\":\"project\",\"roots.config.mazinger.com\":\"root\"},\"name\":\"ac1\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}}"),
				),
			)
			replicationConfig := &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"source": map[string]interface{}{
							"name": "ac1",
							"kind": "Object",
							"object": map[string]interface{}{
								"objectType": map[string]interface{}{
									"group":   Group,
									"version": "v1",
									"kind":    AcKind,
								},
								"hierarchical": true,
								"hierarchy": map[string]interface{}{
									"labels": []map[string]interface{}{
										{
											"key":   Root,
											"value": "root",
										},
										{
											"key":   Project,
											"value": "project",
										},
										{
											"key":   Config,
											"value": "config",
										},
									},
								},
							},
						},
						"destination": map[string]interface{}{
							"hierarchical": false,
						},
						"remoteEndpointGvk": map[string]interface{}{
							"group": "connect.nexus.vmware.com",
							"kind":  "NexusEndpoint",
							"name":  "default",
						},
					},
				},
			}
			err := handler.Create(replicationConfig)
			Expect(err).NotTo(HaveOccurred())
			Expect(server.ReceivedRequests()).Should(HaveLen(2))
		})

		It("Should skip replicationconfigs configured for different endpoints", func() {
			os.Setenv(utils.RemoteEndpointPort, "")
			replicationConfig = &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"remoteEndpointGvk": map[string]interface{}{
							"group": "connect.nexus.vmware.com",
							"kind":  "NexusEndpoint",
							"name":  "default",
						},
					},
				},
			}

			conf, err := config.LoadConfig("./../config/test_utils/correct.yaml")
			Expect(err).NotTo(HaveOccurred())
			handler = handlers.NewReplicationConfigHandler(gvr, conf, client)

			// Server will not receive this create request.
			err = handler.Create(replicationConfig)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should fail when configured with invalid fields", func() {
			replicationConfig = &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"source": "inValidValue",
					},
				},
			}
			err := handler.Create(replicationConfig)
			Expect(err.Error()).To(ContainSubstring("failed to unmarshal replicationconfig spec"))
		})

		It("Should fail when endpoint object not found", func() {
			replicationConfig = &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{},
				},
			}
			err := handler.Create(replicationConfig)
			Expect(err.Error()).To(ContainSubstring("failed to get endpoint object"))
		})

		It("Should fail when invalid endpoint object is configured", func() {
			replicationConfig = &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"remoteEndpointGvk": map[string]interface{}{
							"group": "connect.nexus.vmware.com",
							"kind":  "NexusEndpoint",
							"name":  "defaultNew",
						},
					},
				},
			}

			err := handler.Create(replicationConfig)
			Expect(err.Error()).To(ContainSubstring("failed to unmarshal endpoint spec of defaultNew"))
		})

		It("Should fail when replicationconfig of wrong type is configured", func() {
			replicationConfig = &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": math.Inf(1),
				},
			}

			err := handler.Create(replicationConfig)
			Expect(err.Error()).To(ContainSubstring("failed to marshal replicationconfig spec"))
		})

		It("Should fail when nexusendpoint of wrong type is configured", func() {
			replicationConfig = &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"remoteEndpointGvk": map[string]interface{}{
							"group": "connect.nexus.vmware.com",
							"kind":  "NexusEndpoint",
							"name":  "wrongType",
						},
					},
				},
			}
			err := handler.Create(replicationConfig)
			Expect(err.Error()).To(ContainSubstring("failed to marshal endpoint spec "))
		})

		It("Should fail when nexusendpoint with wrong cert is configured", func() {
			replicationConfig = &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"remoteEndpointGvk": map[string]interface{}{
							"group": "connect.nexus.vmware.com",
							"kind":  "NexusEndpoint",
							"name":  "wrongCert",
						},
					},
				},
			}
			err := handler.Create(replicationConfig)
			Expect(err.Error()).To(ContainSubstring("could not decode cert: illegal base64 data at input byte 0"))
		})
	})
})
