package handlers_test

import (
	"context"
	"fmt"
	"math"
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	fake_dynamic "k8s.io/client-go/dynamic/fake"

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

			endpointGvr := schema.GroupVersionResource{Group: "connect.nexus.org", Version: "v1", Resource: "nexusendpoints"}
			acGvr := schema.GroupVersionResource{Group: "config.mazinger.com", Version: "v1", Resource: "apicollaborationspaces"}
			adGvr := schema.GroupVersionResource{Group: "config.mazinger.com", Version: "v1", Resource: "apidevspaces"}
			gvrToListKind[acGvr] = "ApiCollaborationSpaceList"
			gvrToListKind[endpointGvr] = "NexusEndpointList"
			gvrToListKind[adGvr] = "ApiDevSpaceList"

			scheme := runtime.NewScheme()
			client = fake_dynamic.NewSimpleDynamicClientWithCustomListKinds(scheme, gvrToListKind)

			nexusEndpoint := &unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name": "default",
					},
					"spec": map[string]interface{}{
						"host": fmt.Sprintf("http://%s", parts[0]),
						"port": parts[1],
						"cert": "abc",
					},
				},
			}
			_, err := client.Resource(endpointGvr).Create(context.TODO(), nexusEndpoint, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())

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
			_, err = client.Resource(endpointGvr).Create(context.TODO(), inValidNexusEndpoint, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())

			wrongTypeNexusEndpoint := &unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name": "wrongType",
					},
					"spec": map[string]interface{}{
						"port": math.Inf(1),
					},
				},
			}
			_, err = client.Resource(endpointGvr).Create(context.TODO(), wrongTypeNexusEndpoint, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())

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
			_, err = client.Resource(acGvr).Create(context.TODO(), ac1, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())

			gvr = utils.GetGVRFromCrdType(utils.ReplicationConfigCRD)
			conf, err = config.LoadConfig("./../config/test_utils/correct.yaml")
			Expect(err).NotTo(HaveOccurred())

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
			replicationConfig = &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"source": map[string]interface{}{
							"kind": "Type",
							"type": map[string]interface{}{
								"group":   "config.mazinger.com",
								"version": "v1",
								"kind":    "ApiCollaborationSpace",
							},
						},
						"destination": map[string]interface{}{
							"hierarchical": false,
						},
						"remoteEndpointGvk": map[string]interface{}{
							"group": "connect.nexus.org",
							"kind":  "NexusEndpoint",
							"name":  "default",
						},
					},
				},
			}
			err := handler.Create(replicationConfig)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should fail object creation when unable to connect to destination.", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/apicollaborationspaces/ac1"),
					ghttp.RespondWith(200, "turbo: true"),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/apis/config.mazinger.com/v1/apicollaborationspaces"),
					ghttp.RespondWith(400, "NotOK"),
				),
			)
			replicationConfig = &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"source": map[string]interface{}{
							"kind": "Type",
							"type": map[string]interface{}{
								"group":   "config.mazinger.com",
								"version": "v1",
								"kind":    "ApiCollaborationSpace",
							},
						},
						"destination": map[string]interface{}{
							"hierarchical": false,
						},
						"remoteEndpointGvk": map[string]interface{}{
							"group": "connect.nexus.org",
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
							"group": "connect.nexus.org",
							"kind":  "NexusEndpoint",
							"name":  "default",
						},
					},
				},
			}
			err := handler.Create(replicationConfig)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should skip replicationconfigs configured for different endpoints", func() {
			os.Setenv(utils.RemoteEndpointPort, "")
			replicationConfig = &unstructured.Unstructured{
				Object: map[string]interface{}{
					"spec": map[string]interface{}{
						"remoteEndpointGvk": map[string]interface{}{
							"group": "connect.nexus.org",
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
							"group": "connect.nexus.org",
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
							"group": "connect.nexus.org",
							"kind":  "NexusEndpoint",
							"name":  "wrongType",
						},
					},
				},
			}
			err := handler.Create(replicationConfig)
			Expect(err.Error()).To(ContainSubstring("failed to marshal endpoint spec "))
		})
	})
})
