package controllers_test

import (
	"context"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	apinxv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/api.nexus.vmware.com/v1"
	confignxv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/config.nexus.vmware.com/v1"
	nxv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/connect.nexus.vmware.com/v1"
	nexus_client "golang-appnet.eng.vmware.com/nexus-sdk/api/build/nexus-client"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

type SingletonNodes struct {
	APINexus *nexus_client.ApiNexus
	Config   *nexus_client.ConfigConfig
	Connect  *nexus_client.ConnectConnect
}

func initDatamodel(ctx context.Context, client *nexus_client.Clientset) *SingletonNodes {
	nexusapi, err := client.AddApiNexus(ctx, &apinxv1.Nexus{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nexus",
		},
		Spec: apinxv1.NexusSpec{},
	})
	Expect(err).NotTo(HaveOccurred())
	Expect(nexusapi.DisplayName()).To(Equal("nexus"))

	config, err := nexusapi.AddConfig(ctx, &confignxv1.Config{
		ObjectMeta: metav1.ObjectMeta{
			Name: "config",
		},
		Spec: confignxv1.ConfigSpec{},
	})
	Expect(err).NotTo(HaveOccurred())
	Expect(config.DisplayName()).To(Equal("config"))

	connect, err := config.AddConnect(ctx, &nxv1.Connect{
		ObjectMeta: metav1.ObjectMeta{
			Name: "connect",
		},
		Spec: nxv1.ConnectSpec{},
	})
	Expect(err).NotTo(HaveOccurred())
	Expect(connect.DisplayName()).To(Equal("connect"))
	return &SingletonNodes{
		APINexus: nexusapi,
		Config:   config,
		Connect:  connect,
	}
}

func initEnvVars() {
	err := os.Setenv("NEXUS_CONNECTOR_VERSION", "v1.0")
	Expect(err).NotTo(HaveOccurred())

	err = os.Setenv("NAMESPACE", "default")
	Expect(err).NotTo(HaveOccurred())
}

func cleanupEnv() {
	fmt.Println("cleaning up environment")
	_ = fakeClient.DeleteApiNexus(context.Background(), "nexus")
}

func newScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	_ = nxv1.AddToScheme(scheme)
	_ = clientgoscheme.AddToScheme(scheme)
	_ = coreV1.AddToScheme(scheme)
	return scheme
}

func getHelperClient(ctx context.Context, endpoint *nexus_client.ConnectNexusEndpoint) *Client {
	client := NewClient()
	client.On("Get",
		mock.IsType(ctx),
		mock.Anything,
		mock.Anything,
	).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(2).(*nxv1.NexusEndpoint)
		arg.ObjectMeta = endpoint.ObjectMeta
		arg.Spec = endpoint.Spec
		arg.Name = endpoint.Name
		arg.Namespace = endpoint.Namespace
		arg.Labels = endpoint.Labels
		arg.Annotations = endpoint.Annotations
		arg.APIVersion = endpoint.APIVersion
		arg.CreationTimestamp = endpoint.CreationTimestamp
		arg.DeletionGracePeriodSeconds = endpoint.DeletionGracePeriodSeconds
		arg.DeletionTimestamp = endpoint.DeletionTimestamp
		arg.GenerateName = endpoint.GenerateName
		arg.Kind = endpoint.Kind
		arg.Finalizers = endpoint.Finalizers
		arg.ResourceVersion = endpoint.ResourceVersion
		arg.Generation = endpoint.Generation
	})
	return client
}
func getHelperClientForDeleteEvent(ctx context.Context) *Client {
	client := NewClient()
	client.On("Get",
		mock.IsType(ctx),
		mock.Anything,
		mock.Anything,
	).Return(errors.NewNotFound(resource("tests"), "3"))
	return client
}

func resource(resource string) schema.GroupResource {
	return schema.GroupResource{Group: "", Resource: resource}
}
