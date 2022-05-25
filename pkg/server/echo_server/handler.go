package echo_server

import (
	"api-gw/pkg/model"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"

	"api-gw/pkg/client"
)

type DefaultResponse struct {
	Message string `json:"message"`
}

// getHandler is used to process GET requests
func getHandler(c echo.Context) error {
	nc := c.(*NexusContext)
	crdName := model.UriToCRDType[nc.NexusURI]
	crdInfo := model.CrdTypeToNodeInfo[crdName]

	// List operation
	list := true

	// Get name from params
	var name string
	for _, param := range nc.ParamNames() {
		if param == crdInfo.Name {
			list = false
			name = nc.Param(param)
			if name == "" {
				if val, ok := nc.Codes[http.StatusBadRequest]; ok {
					return nc.JSON(http.StatusBadRequest, DefaultResponse{Message: val.Description})
				} else {
					log.Debugf("Could not find required param %s for request %s", crdInfo.Name, nc.Request().RequestURI)
					return nc.JSON(http.StatusBadRequest, DefaultResponse{Message: fmt.Sprintf("Could not find required param: %s", crdInfo.Name)})
				}
			}
		}
	}

	// Get name from query params
	if nc.QueryParams().Has(crdInfo.Name) {
		list = false
		name = nc.QueryParams().Get(crdInfo.Name)
	}

	// Mangle name
	labels := parseLabels(nc, crdInfo.ParentHierarchy)
	hashedName := nexus.GetHashedName(crdName, crdInfo.ParentHierarchy, labels, name)

	// Setup GroupVersionResource
	parts := strings.Split(crdName, ".")
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
			itemName := item.GetName()
			if val, ok := item.GetLabels()["nexus/display_name"]; ok {
				itemName = val
			}
			specs[itemName] = item.Object["spec"]
		}
		output = specs
	} else {
		obj, err := client.Client.Resource(gvr).Get(context.TODO(), hashedName, metav1.GetOptions{})
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
	crdName := model.UriToCRDType[nc.NexusURI]
	crdInfo := model.CrdTypeToNodeInfo[crdName]

	// Get name from the URI segment
	var name string
	for _, param := range nc.ParamNames() {
		if param == crdInfo.Name {
			name = nc.Param(param)
		}
	}

	// Get name from query params
	if val := nc.QueryParam(crdInfo.Name); val != "" {
		name = val
	}

	if name == "" {
		if val, ok := nc.Codes[http.StatusBadRequest]; ok {
			return nc.JSON(http.StatusBadRequest, DefaultResponse{Message: val.Description})
		} else {
			log.Debugf("Could not find required param %s for request %s", crdInfo.Name, nc.Request().RequestURI)
			return nc.JSON(http.StatusBadRequest, DefaultResponse{Message: fmt.Sprintf("Could not find required param: %s", crdInfo.Name)})
		}
	}

	// Parse body
	body := make(map[string]interface{})
	if err := nc.Bind(&body); err != nil {
		return err
	}

	// Setup GroupVersionResource
	parts := strings.Split(crdName, ".")
	gvr := schema.GroupVersionResource{
		Group:    strings.Join(parts[1:], "."),
		Version:  "v1",
		Resource: parts[0],
	}
	crdNameParts := strings.Split(crdInfo.Name, ".")
	labels := parseLabels(nc, crdInfo.ParentHierarchy)
	labels["nexus/is_name_hashed"] = "true"
	labels["nexus/display_name"] = name
	labels[crdInfo.Name] = name

	// Mangle name
	hashedName := nexus.GetHashedName(crdName, crdInfo.ParentHierarchy, labels, name)
	obj, err := client.Client.Resource(gvr).Get(context.TODO(), hashedName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			// Build object
			err = client.CreateObject(gvr,
				crdNameParts[0], hashedName, labels, body)
			if err != nil {
				return handleClientError(nc, err)
			}

			// Update parent
			var err error
			if len(crdInfo.ParentHierarchy) > 0 {
				parentCrdName := crdInfo.ParentHierarchy[len(crdInfo.ParentHierarchy)-1]
				parentCrd := model.CrdTypeToNodeInfo[parentCrdName]
				err = client.UpdateParentWithAddedChild(parentCrdName, parentCrd, labels, crdName, name, hashedName)
			}

			if err == nil {
				return c.JSON(http.StatusOK, DefaultResponse{Message: name})
			}
		}
		return handleClientError(nc, err)
	}

	obj.SetLabels(labels)
	obj.Object["spec"] = body

	// Update resource
	_, err = client.Client.Resource(gvr).Update(context.TODO(), obj, metav1.UpdateOptions{})
	if err != nil {
		return handleClientError(nc, err)
	}
	return c.JSON(http.StatusOK, DefaultResponse{Message: name})
}

// deleteHandler is used to process DELETE requests
func deleteHandler(c echo.Context) error {
	nc := c.(*NexusContext)
	crdName := model.UriToCRDType[nc.NexusURI]
	crdInfo := model.CrdTypeToNodeInfo[crdName]

	// Get name from params
	var name string
	for _, param := range nc.ParamNames() {
		if param == crdInfo.Name {
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
	if nc.QueryParams().Has(crdInfo.Name) {
		name = nc.QueryParams().Get(crdInfo.Name)
	}

	// Mangle name
	labels := parseLabels(nc, crdInfo.ParentHierarchy)
	hashedName := nexus.GetHashedName(crdName, crdInfo.ParentHierarchy, labels, name)

	// Setup GroupVersionResource
	parts := strings.Split(crdName, ".")
	gvr := schema.GroupVersionResource{
		Group:    strings.Join(parts[1:], "."),
		Version:  "v1",
		Resource: parts[0],
	}

	// Get object from kubernetes
	err := client.DeleteObject(gvr, crdName, crdInfo, hashedName, name)
	if err != nil {
		return handleClientError(nc, err)
	}

	return nc.NoContent(http.StatusOK)
}

// handleClientError is used to parse client errors and map them to the corresponding statuses from HTTPCodesResponses
func handleClientError(c echo.Context, err error) error {
	nc := c.(*NexusContext)
	log.Warn(err)

	switch {
	case errors.IsNotFound(err):
		if val, ok := nc.Codes[http.StatusNotFound]; ok {
			return c.JSON(http.StatusNotFound, DefaultResponse{Message: val.Description})
		}
	case errors.IsAlreadyExists(err), errors.IsConflict(err):
		if val, ok := nc.Codes[http.StatusConflict]; ok {
			return c.JSON(http.StatusConflict, DefaultResponse{Message: val.Description})
		}
	case errors.IsInternalError(err):
		if val, ok := nc.Codes[http.StatusInternalServerError]; ok {
			return c.JSON(http.StatusInternalServerError, DefaultResponse{Message: val.Description})
		}
	case errors.IsBadRequest(err):
		if val, ok := nc.Codes[http.StatusBadRequest]; ok {
			return c.JSON(http.StatusBadRequest, DefaultResponse{Message: val.Description})
		}
	case errors.IsForbidden(err):
		if val, ok := nc.Codes[http.StatusForbidden]; ok {
			return c.JSON(http.StatusForbidden, DefaultResponse{Message: val.Description})
		}
	case errors.IsGone(err):
		if val, ok := nc.Codes[http.StatusGone]; ok {
			return c.JSON(http.StatusGone, DefaultResponse{Message: val.Description})
		}
	case errors.IsInvalid(err):
		if val, ok := nc.Codes[http.StatusUnprocessableEntity]; ok {
			return c.JSON(http.StatusUnprocessableEntity, DefaultResponse{Message: val.Description})
		}
	}

	return c.JSON(http.StatusInternalServerError, DefaultResponse{Message: err.Error()})
}

func parseLabels(c echo.Context, parents []string) map[string]string {
	nc := c.(*NexusContext)
	// Parse labels
	labels := make(map[string]string)
	for _, parent := range parents {
		if c, ok := model.CrdTypeToNodeInfo[parent]; ok {
			if v := nc.Param(c.Name); v != "" {
				labels[parent] = v
			} else if nc.QueryParams().Has(c.Name) {
				labels[parent] = nc.QueryParams().Get(c.Name)
			} else {
				labels[parent] = "default"
			}
		}
	}

	return labels
}
