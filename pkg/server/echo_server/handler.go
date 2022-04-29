package echo_server

import (
	"api-gw/controllers"
	"api-gw/pkg/client"
	"context"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"net/http"
	"strings"
)

type DefaultResponse struct {
	Message string `json:"message"`
}

// getHandler is used to process GET requests
func getHandler(c echo.Context) error {
	nc := c.(*NexusContext)
	crdType := controllers.GlobalURIToCRDType[nc.NexusURI]
	crd := controllers.GlobalCRDTypeToNode[crdType]

	// List operation
	list := true

	// Get name from params
	var name string
	for _, param := range nc.ParamNames() {
		if param == crd.Name {
			list = false
			name = nc.Param(param)
			if name == "" {
				if val, ok := nc.Codes[http.StatusBadRequest]; ok {
					return nc.JSON(http.StatusBadRequest, DefaultResponse{Message: val.Description})
				} else {
					return nc.NoContent(http.StatusInternalServerError)
				}
			}
		}
	}

	// Get name from query params
	if nc.QueryParams().Has(crd.Name) {
		list = false
		name = nc.QueryParams().Get(crd.Name)
	}

	// Setup GroupVersionResource
	parts := strings.Split(crdType, ".")
	gvr := schema.GroupVersionResource{
		Group:    strings.Join(parts[1:], "."),
		Version:  "v1",
		Resource: parts[0],
	}

	var output interface{}
	if list {
		specs := make(map[string]interface{})
		objs, err := client.Client.Resource(gvr).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return handleClientError(nc, err)
		}
		for _, item := range objs.Items {
			specs[item.GetName()] = item.Object["spec"]
		}
		output = specs
	} else {
		obj, err := client.Client.Resource(gvr).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return handleClientError(nc, err)
		}
		output = obj.Object["spec"]
	}

	return nc.JSON(http.StatusOK, output)
}

// putHandler is used to process PUT requests
func putHandler(c echo.Context) error {
	nc := c.(*NexusContext)
	crdType := controllers.GlobalURIToCRDType[nc.NexusURI]
	crd := controllers.GlobalCRDTypeToNode[crdType]

	// Get name from the URI segment
	var name string
	for _, param := range nc.ParamNames() {
		if param == crd.Name {
			name = nc.Param(param)
		}
	}

	// Get name from query params
	if val := nc.QueryParam(crd.Name); val != "" {
		name = val
	}

	if name == "" {
		if val, ok := nc.Codes[http.StatusBadRequest]; ok {
			return nc.JSON(http.StatusBadRequest, DefaultResponse{Message: val.Description})
		} else {
			return nc.NoContent(http.StatusInternalServerError)
		}
	}

	// Parse labels
	labels := make(map[string]string)
	for _, parent := range crd.ParentHierarchy {
		if c, ok := controllers.GlobalCRDTypeToNode[parent]; ok {
			if v := nc.Param(c.Name); v != "" {
				labels[parent] = v
			} else {
				labels[parent] = "default"
			}
		}
	}

	// Parse body
	var body map[string]interface{}
	if err := nc.Bind(&body); err != nil {
		return err
	}

	// Setup GroupVersionResource
	parts := strings.Split(crdType, ".")
	gvr := schema.GroupVersionResource{
		Group:    strings.Join(parts[1:], "."),
		Version:  "v1",
		Resource: parts[0],
	}
	crdNameParts := strings.Split(crd.Name, ".")

	// Build object
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": gvr.GroupVersion().String(),
			"kind":       crdNameParts[0],
			"metadata": map[string]interface{}{
				"name":   nexus.GetHashedName(crdType, crd.ParentHierarchy, labels, name),
				"labels": labels,
			},
			"spec": body,
		},
	}

	// Create resource
	res, err := client.Client.Resource(gvr).Create(context.TODO(), obj, metav1.CreateOptions{})
	if err != nil {
		return handleClientError(nc, err)
	}
	return c.JSON(http.StatusOK, DefaultResponse{Message: res.GetName()})
}

// deleteHandler is used to process DELETE requests
func deleteHandler(c echo.Context) error {
	nc := c.(*NexusContext)
	crdType := controllers.GlobalURIToCRDType[nc.NexusURI]
	crd := controllers.GlobalCRDTypeToNode[crdType]

	// Get name from params
	var name string
	for _, param := range nc.ParamNames() {
		if param == crd.Name {
			name = nc.Param(param)
			if name == "" {
				if val, ok := nc.Codes[http.StatusBadRequest]; ok {
					return nc.JSON(http.StatusBadRequest, DefaultResponse{Message: val.Description})
				} else {
					return nc.NoContent(http.StatusInternalServerError)
				}
			}
		}
	}

	// Get name from query params
	if nc.QueryParams().Has(crd.Name) {
		name = nc.QueryParams().Get(crd.Name)
	}

	// Setup GroupVersionResource
	parts := strings.Split(crdType, ".")
	gvr := schema.GroupVersionResource{
		Group:    strings.Join(parts[1:], "."),
		Version:  "v1",
		Resource: parts[0],
	}

	// Get object from kubernetes
	err := client.Client.Resource(gvr).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return handleClientError(nc, err)
	}

	return nc.NoContent(http.StatusOK)
}

// handleClientError is used to parse client errors and map them to the corresponding statuses from HTTPCodesResponses
func handleClientError(c echo.Context, err error) error {
	nc := c.(*NexusContext)
	log.Warn(err)

	if errors.IsNotFound(err) {
		if val, ok := nc.Codes[http.StatusNotFound]; ok {
			return c.JSON(http.StatusNotFound, DefaultResponse{Message: val.Description})
		}
	}

	return c.JSON(http.StatusInternalServerError, DefaultResponse{Message: err.Error()})
}
