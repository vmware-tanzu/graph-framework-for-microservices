package k8s_utils

import (
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	ext "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type ClientInt interface {
	ApplyCrd(crd v1.CustomResourceDefinition) error
	DeleteCrd(name string) error
	FetchCrds() error
	GetCrds() []v1.CustomResourceDefinition
	GetGroup() string
	GetCrd(name string) *v1.CustomResourceDefinition
	ListResources(crd v1.CustomResourceDefinition) ([]interface{}, error)
}

type Client struct {
	Clientset ext.Interface
	crds      []v1.CustomResourceDefinition
	group     string
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
	c.crds = l.Items
	if len(l.Items) > 0 {
		parts := strings.Split(l.Items[0].Spec.Group, ".")
		c.group = strings.Join(parts[1:], ".")
	}
	return err
}

func (c *Client) ListResources(crd v1.CustomResourceDefinition) ([]interface{}, error) {
	data, err := c.Clientset.ApiextensionsV1beta1().RESTClient().Get().RequestURI(createURI(crd)).DoRaw(context.TODO())
	if err != nil {
		return nil, err
	}
	var obj map[string]interface{}
	if err = yaml.Unmarshal(data, &obj); err != nil {
		return nil, err
	}
	return obj["items"].([]interface{}), err
}

func createURI(crd v1.CustomResourceDefinition) string {
	return fmt.Sprintf("apis/%s/v1/%s", crd.Spec.Group, crd.Spec.Names.Plural)
}
