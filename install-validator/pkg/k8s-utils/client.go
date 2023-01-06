package k8s_utils

import (
	"context"
	"fmt"
	"os"
	"strings"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	ext "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClientInt interface {
	ApplyCrd(crd v1.CustomResourceDefinition) error
	DeleteCrd(name string) error
	FetchCrds() error
	GetCrds() []v1.CustomResourceDefinition
	GetGroup() string
	GetCrd(name string) *v1.CustomResourceDefinition
	ListResources(crd v1.CustomResourceDefinition) ([]unstructured.Unstructured, error)
	FetchGroup(groupPath string) error
}

type Client struct {
	Clientset     ext.Interface
	DynamicClient dynamic.Interface
	crds          []v1.CustomResourceDefinition
	group         string
}

func (c *Client) GetCrd(name string) *v1.CustomResourceDefinition {
	for _, s := range c.crds {
		if s.Name == name {
			return &s
		}
	}
	return nil
}

func (c *Client) GetCrds() []v1.CustomResourceDefinition {
	return c.crds
}

func (c *Client) GetGroup() string {
	return c.group
}

func (c *Client) DeleteCrd(name string) error {
	return c.Clientset.ApiextensionsV1().CustomResourceDefinitions().Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (c *Client) ApplyCrd(crd v1.CustomResourceDefinition) error {
	oldCrd := c.GetCrd(crd.Name)
	var err error
	if oldCrd == nil {
		_, err = c.Clientset.ApiextensionsV1().CustomResourceDefinitions().Create(context.TODO(), &crd, metav1.CreateOptions{})
	} else {
		crd.ObjectMeta.ResourceVersion = oldCrd.ObjectMeta.ResourceVersion
		_, err = c.Clientset.ApiextensionsV1().CustomResourceDefinitions().Update(context.TODO(), &crd, metav1.UpdateOptions{})
	}
	return err
}

func (c *Client) FetchCrds() error {
	l, err := c.Clientset.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	c.crds = l.Items
	return nil
}

func (c *Client) ListResources(crd v1.CustomResourceDefinition) ([]unstructured.Unstructured, error) {
	data, err := c.DynamicClient.Resource(
		schema.GroupVersionResource{
			Group:    crd.Spec.Group,
			Version:  crd.Spec.Versions[0].Name,
			Resource: crd.Spec.Names.Plural,
		}).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return data.Items, err
}

func (c *Client) FetchGroup(groupPath string) error {
	file, err := os.ReadFile(fmt.Sprintf("%s", groupPath))
	if err != nil {
		return err
	}
	c.group = strings.TrimSpace(string(file))
	return nil
}

func GetRestConfig() (*rest.Config, error) {
	filePath := os.Getenv("KUBECONFIG")
	if filePath != "" {
		config, err := clientcmd.BuildConfigFromFlags("", filePath)
		if err != nil {
			return nil, err
		}
		return config, nil
	}
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}
