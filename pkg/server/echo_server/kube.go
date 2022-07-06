package echo_server

import (
	"api-gw/pkg/client"
	"api-gw/pkg/model"
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"net/http/httputil"
	"net/url"
)

// kubeSetupProxy is used to set up reverse proxy to an API server
func kubeSetupProxy(e *echo.Echo) {
	proxyUrl, err := url.Parse(client.Host)
	if err != nil {
		log.Warnf("Could not parse proxy URL: %v", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
	e.Any("*", echo.WrapHandler(proxy))
}

// kubeGetByNameHandler is used to process 'kubectl get <resource> <name>' requests
func kubeGetByNameHandler(c echo.Context) error {
	nc := c.(*NexusContext)

	gvr := schema.GroupVersionResource{
		Group:    nc.GroupName,
		Version:  "v1",
		Resource: nc.Resource,
	}
	obj, err := client.GetObject(gvr, c.Param("name"), metav1.GetOptions{})
	if err != nil {
		if status := kerrors.APIStatus(nil); errors.As(err, &status) {
			return c.JSON(int(status.Status().Code), status.Status())
		}
		c.Error(err)
	}

	return c.JSON(200, obj)
}

// kubeGetHandler is used to process `kubectl get <resource>' requests
func kubeGetHandler(c echo.Context) error {
	nc := c.(*NexusContext)

	opts := metav1.ListOptions{}
	if c.QueryParams().Has("labelSelector") {
		opts.LabelSelector = c.QueryParams().Get("labelSelector")
	}

	gvr := schema.GroupVersionResource{
		Group:    nc.GroupName,
		Version:  "v1",
		Resource: nc.Resource,
	}

	obj, err := client.Client.Resource(gvr).List(context.TODO(), opts)
	if err != nil {
		if status := kerrors.APIStatus(nil); errors.As(err, &status) {
			return c.JSON(int(status.Status().Code), status.Status())
		}
		c.Error(err)
	}
	return c.JSON(200, obj)
}

// kubePostHandler is used to process `kubectl apply` requests
func kubePostHandler(c echo.Context) error {
	nc := c.(*NexusContext)
	crdInfo := model.CrdTypeToNodeInfo[nc.CrdType]

	body := &unstructured.Unstructured{}
	if err := c.Bind(&body); err != nil {
		return err
	}
	labels := body.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	labels["nexus/is_name_hashed"] = "true"
	labels["nexus/display_name"] = body.GetName()

	orderedLabels := nexus.ParseCRDLabels(crdInfo.ParentHierarchy, labels)
	for _, key := range orderedLabels.Keys() {
		value, _ := orderedLabels.Get(key)
		labels[key.(string)] = value.(string)
	}

	hashedName := nexus.GetHashedName(nc.CrdType, crdInfo.ParentHierarchy, labels, body.GetName())
	body.SetLabels(labels)
	body.SetName(hashedName)

	gvr := schema.GroupVersionResource{
		Group:    nc.GroupName,
		Version:  "v1",
		Resource: nc.Resource,
	}

	// Get object to check if it exists
	obj, err := client.GetObject(gvr, hashedName, metav1.GetOptions{})
	if err != nil {

		// Create object if is not found
		if kerrors.IsNotFound(err) {
			if _, ok := body.UnstructuredContent()["spec"]; !ok {
				content := body.UnstructuredContent()
				content["spec"] = map[string]interface{}{}
				body.SetUnstructuredContent(content)
			}
			obj, err = client.Client.Resource(gvr).Create(context.TODO(), body, metav1.CreateOptions{})
			if err != nil {
				if status := kerrors.APIStatus(nil); errors.As(err, &status) {
					return c.JSON(int(status.Status().Code), status.Status())
				}
				c.Error(err)
			}

			var err error
			if len(crdInfo.ParentHierarchy) > 0 {
				parentCrdName := crdInfo.ParentHierarchy[len(crdInfo.ParentHierarchy)-1]
				parentCrd := model.CrdTypeToNodeInfo[parentCrdName]
				err = client.UpdateParentWithAddedChild(parentCrdName, parentCrd, labels, crdInfo, nc.CrdType, body.GetName(), hashedName)
			}

			if err != nil {
				if status := kerrors.APIStatus(nil); errors.As(err, &status) {
					return c.JSON(int(status.Status().Code), status.Status())
				}
				c.Error(err)
			}

			return c.JSON(201, obj)
		}

		if status := kerrors.APIStatus(nil); errors.As(err, &status) {
			return c.JSON(int(status.Status().Code), status.Status())
		}
		c.Error(err)
	}

	body.SetResourceVersion(obj.GetResourceVersion())
	obj, err = client.Client.Resource(gvr).Update(context.TODO(), body, metav1.UpdateOptions{})
	if err != nil {
		if status := kerrors.APIStatus(nil); errors.As(err, &status) {
			return c.JSON(int(status.Status().Code), status.Status())
		}
		c.Error(err)
	}

	return c.JSON(200, obj)
}

func kubeDeleteHandler(c echo.Context) error {
	nc := c.(*NexusContext)
	crdInfo := model.CrdTypeToNodeInfo[nc.CrdType]
	gvr := schema.GroupVersionResource{
		Group:    nc.GroupName,
		Version:  "v1",
		Resource: nc.Resource,
	}

	err := client.DeleteObject(gvr, nc.CrdType, crdInfo, c.Param("name"))
	if err != nil {
		if status := kerrors.APIStatus(nil); errors.As(err, &status) {
			return c.JSON(int(status.Status().Code), status.Status())
		}
		c.Error(err)
	}

	return c.JSON(200, map[string]interface{}{
		"kind":       "Status",
		"apiVersion": "v1",
		"metadata":   map[string]interface{}{},
		"status":     "Success",
		"details": map[string]interface{}{
			"name":  c.Param("name"),
			"group": nc.GroupName,
			"kind":  nc.Resource,
		},
	})
}
