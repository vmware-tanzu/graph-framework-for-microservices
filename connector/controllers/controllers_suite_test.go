package controllers_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/stretchr/testify/mock"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controllers Suite")
}

func getHelperClient(ctx context.Context, crd *apiextensionsv1.CustomResourceDefinition) *Client {
	client := NewClient()
	client.On("Get",
		mock.IsType(ctx),
		mock.Anything,
		mock.Anything,
	).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(2).(*apiextensionsv1.CustomResourceDefinition)
		arg.ObjectMeta = crd.ObjectMeta
		arg.Spec = crd.Spec
		arg.Name = crd.Name
		arg.Namespace = crd.Namespace
		arg.Labels = crd.Labels
		arg.Annotations = crd.Annotations
		arg.APIVersion = crd.APIVersion
		arg.CreationTimestamp = crd.CreationTimestamp
		arg.DeletionGracePeriodSeconds = crd.DeletionGracePeriodSeconds
		arg.DeletionTimestamp = crd.DeletionTimestamp
		arg.GenerateName = crd.GenerateName
		arg.Kind = crd.Kind
		arg.Finalizers = crd.Finalizers
		arg.ResourceVersion = crd.ResourceVersion
		arg.Generation = crd.Generation
		arg.Status = crd.Status
	})
	return client
}

func getHelperClientForDeleteEvent(ctx context.Context) *Client {
	client := NewClient()
	client.On("Get",
		mock.IsType(ctx),
		mock.Anything,
		mock.Anything,
	).Return(errors.NewNotFound(resource("tests"), "3")).Run(func(args mock.Arguments) {
		arg := args.Get(2).(*apiextensionsv1.CustomResourceDefinition)
		arg.Status = crd.Status
	})
	return client
}

func resource(resource string) schema.GroupResource {
	return schema.GroupResource{Group: "", Resource: resource}
}

func getCRD() *apiextensionsv1.CustomResourceDefinition {
	return &apiextensionsv1.CustomResourceDefinition{
		ObjectMeta: v1.ObjectMeta{
			Annotations: map[string]string{
				"nexus": `{"name":"ApiCollaborationSpace.config","hierarchy":["roots.apix.mazinger.com","projects.apix.mazinger.com","configs.apix.mazinger.com"],"children":{"apidevspaces.config.mazinger.com":{"fieldName":"ApiDevSpaces","fieldNameGvk":"apiDevSpacesGvk","isNamed":true}}}`,
			},
			Name: "apicollaborationspaces.config.mazinger.com",
		},
		Spec: apiextensionsv1.CustomResourceDefinitionSpec{
			Group: "config.mazinger.com",
			Names: apiextensionsv1.CustomResourceDefinitionNames{
				Plural:     "apicollaborationspaces",
				Singular:   "apicollaborationspace",
				Kind:       "ApiCollaborationSpace",
				ListKind:   "ApiCollaborationSpaceList",
				ShortNames: []string{"apicollaborationspace"},
			},
			Scope: "Cluster",
			Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
				{
					Name: "v1",
				},
			},
		},
		Status: apiextensionsv1.CustomResourceDefinitionStatus{
			StoredVersions: []string{"v1"},
		},
	}
}
