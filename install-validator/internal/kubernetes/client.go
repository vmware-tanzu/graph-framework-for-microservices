package kubernetes

import (
	"context"
	"fmt"

	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	ext "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"os"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
)

type Client struct {
	clientset  *ext.Clientset
	clientset2 *rest.RESTClient
	config     *rest.Config
	crds       []v1.CustomResourceDefinition
}

func NewClient() (Client, error) {
	c := Client{}
	err := c.setConfig()
	if err != nil {
		return c, err
	}
	err = c.setClientSet()
	if err != nil {
		return c, err
	}
	return c, nil
}

func (c *Client) setClientSet() error {
	clientset, err := ext.NewForConfig(c.config)

	if err != nil {
		return err
	}
	c.clientset = clientset

	return nil
}

func (c *Client) setConfig() error {
	config, err := rest.InClusterConfig()
	if err == nil {
		c.config = config
		return nil
	}

	home, exists := os.LookupEnv("HOME")
	if !exists {
		home = "/root"
	}
	configPath := filepath.Join(home, ".kube", "config")
	config, err = clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return err
	}
	c.config = config
	return nil
}

func (c *Client) GetCrd(name string) *v1.CustomResourceDefinition {
	for _, s := range c.crds {
		if s.Name == name {
			return &s
		}
	}
	return nil
}

func (c *Client) ApplyCrd(crd v1.CustomResourceDefinition) error {
	oldCrd := c.GetCrd(crd.Name)
	var err error
	if oldCrd == nil {
		_, err = c.clientset.ApiextensionsV1().CustomResourceDefinitions().Create(context.TODO(), &crd, metav1.CreateOptions{})
	} else {
		crd.ObjectMeta.ResourceVersion = oldCrd.ObjectMeta.ResourceVersion
		_, err = c.clientset.ApiextensionsV1().CustomResourceDefinitions().Update(context.TODO(), &crd, metav1.UpdateOptions{})
	}
	return err
}

func (c *Client) ListCrds() error {
	l, err := c.clientset.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{})
	c.crds = l.Items
	return err
}

func (c *Client) ListResources(crd v1.CustomResourceDefinition) ([]interface{}, error) {
	data, err := c.clientset.ApiextensionsV1beta1().RESTClient().Get().RequestURI(createURI(crd)).DoRaw(context.TODO())
	var obj map[string]interface{}
	if err := yaml.Unmarshal(data, &obj); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return obj["items"].([]interface{}), err
}

func createURI(crd v1.CustomResourceDefinition) string {
	return fmt.Sprintf("apis/%s/v1/%s", crd.Spec.Group, crd.Spec.Names.Plural)
}
