package echo_server

import (
	"api-gw/pkg/config"
	"api-gw/pkg/openapi/declarative"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8sLabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"

	"api-gw/pkg/client"
	"api-gw/pkg/model"
	"api-gw/pkg/utils"

	"github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus"
)

type DefaultResponse struct {
	Message string `json:"message"`
}

// getHandler is used to process GET requests
func getHandler(c echo.Context) error {
	nc := c.(*NexusContext)
	crdName := model.UriToCRDType[nc.NexusURI]
	crdInfo := model.CrdTypeToNodeInfo[crdName]
	// Get name from params
	name := nexus.DEFAULT_KEY
	for _, param := range nc.ParamNames() {
		if param == crdInfo.Name {
			name = nc.Param(param)
			if name == "" {
				if val, ok := nc.Codes[http.StatusBadRequest]; ok {
					return nc.JSON(http.StatusBadRequest, DefaultResponse{Message: val.Description})
				} else {
					log.Debugf("Could not find required param %s for request %s", crdInfo.Name, nc.Request().RequestURI)
					log.Debugf("crdName: %s, nexusURI: %s, paramNames: %s, paramValues: %s", crdInfo.Name, nc.NexusURI, nc.ParamNames(), nc.ParamValues())
					return nc.JSON(http.StatusBadRequest, DefaultResponse{Message: fmt.Sprintf("Could not find required param: %s", crdInfo.Name)})
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

	obj, err := client.Client.Resource(gvr).Get(context.TODO(), hashedName, metav1.GetOptions{})
	if err != nil {
		return handleClientError(nc, err)
	}
	status := make(map[string]interface{})
	if _, ok := obj.Object["status"]; ok {
		status = obj.Object["status"].(map[string]interface{})
	}
	delete(status, "nexus")

	uriInfo, ok := model.GetUriInfo(nc.NexusURI)
	if ok {
		if uriInfo.TypeOfURI == model.SingleLinkURI || uriInfo.TypeOfURI == model.NamedLinkURI {
			return getLinkInfo(nc, uriInfo.TypeOfURI, crdInfo, obj)
		}
		if uriInfo.TypeOfURI == model.StatusURI {
			return nc.JSON(http.StatusOK, status)
		}
	}

	spec := make(map[string]interface{})
	if _, ok := obj.Object["spec"]; ok {
		spec = obj.Object["spec"].(map[string]interface{})
	}
	for _, v := range crdInfo.Children {
		delete(spec, v.FieldNameGvk)
	}
	for _, v := range crdInfo.Links {
		delete(spec, v.FieldNameGvk)
	}

	r := make(map[string]interface{})
	r["spec"] = spec
	r["status"] = status

	return nc.JSON(http.StatusOK, r)
}

// getLinkInfo returns the children/links of parent node based on the requested gvk.
func getLinkInfo(nc *NexusContext, uriType model.URIType, crdInfo model.NodeInfo, obj *unstructured.Unstructured) error {
	splittedUri := strings.Split(nc.NexusURI, "/")
	if len(splittedUri) < 2 {
		return nc.JSON(http.StatusBadRequest, DefaultResponse{Message: "Couldn't determine child object"})
	}

	linkFieldName := splittedUri[len(splittedUri)-1]
	var gvkField string
	for _, child := range crdInfo.Children {
		if child.FieldName == linkFieldName {
			gvkField = child.FieldNameGvk
		}
	}
	if gvkField == "" {
		for _, link := range crdInfo.Links {
			if link.FieldName == linkFieldName {
				gvkField = link.FieldNameGvk
			}
		}
	}
	if gvkField == "" {
		return nc.JSON(http.StatusInternalServerError, DefaultResponse{Message: "Couldn't determine gvk of link"})
	}
	spec, ok := obj.Object["spec"].(map[string]interface{})
	if !ok {
		return nc.JSON(http.StatusInternalServerError, DefaultResponse{Message: "Failed to parse spec of object"})
	}

	log.Debugf("URI %s, splitted URI %s, childFieldName %s, gvkField %s, spec %s, spec[gvkField] %s\n", nc.NexusURI,
		splittedUri, linkFieldName, gvkField, spec, spec[gvkField])

	if uriType == model.SingleLinkURI {
		l := &model.LinkGvk{}
		marshaled, err := json.Marshal(spec[gvkField])
		if err != nil {
			return nc.JSON(http.StatusInternalServerError, DefaultResponse{Message: "Couldn't marshal gvk of link"})
		}
		err = json.Unmarshal(marshaled, l)
		if err != nil {
			return nc.JSON(http.StatusInternalServerError, DefaultResponse{Message: "Couldn't unmarshal gvk of link"})
		}

		if len(l.Group) != 0 {
			resourceName := utils.GetGroupResourceName(l.Kind)
			item, err := getUnstructuredObject(l.Group, resourceName, l.Name)
			if err != nil {
				log.Errorf("Couldn't find object %q", l.Name)
				return nc.JSON(http.StatusNotFound, DefaultResponse{Message: "Couldn't find object"})
			}

			// set parent hierarchy
			crdType := utils.GetCrdType(l.Kind, l.Group)
			if crdNodeInfo, ok := model.GetCRDTypeToNodeInfo(crdType); ok {
				l.Hierarchy = utils.GetParentHierarchy(crdNodeInfo.ParentHierarchy, item.GetLabels())
			}

			// get display name of the object.
			if val, ok := item.GetLabels()[utils.DISPLAY_NAME_LABEL]; ok {
				l.Name = val
			}
			l.Group = l.Group + "/v1"
		}
		return nc.JSON(http.StatusOK, l)
	}

	if uriType == model.NamedLinkURI {
		m := make(map[string]model.LinkGvk)
		marshaled, err := json.Marshal(spec[gvkField])
		if err != nil {
			return nc.JSON(http.StatusInternalServerError, DefaultResponse{Message: "Couldn't marshal gvk of link"})
		}
		err = json.Unmarshal(marshaled, &m)
		if err != nil {
			return nc.JSON(http.StatusInternalServerError, DefaultResponse{Message: "Couldn't unmarshal gvk of link"})
		}

		list := make([]model.LinkGvk, len(m))
		i := 0
		hierarchy := []string{}
		for k, link := range m {
			// set parent hierarchy
			if i == 0 {
				resourceName := utils.GetGroupResourceName(link.Kind)
				item, err := getUnstructuredObject(link.Group, resourceName, link.Name)
				if err != nil {
					log.Errorf("Couldn't find object, skipping... %q", link.Name)
					continue
				}
				crdType := utils.GetCrdType(link.Kind, link.Group)
				if crdNodeInfo, ok := model.GetCRDTypeToNodeInfo(crdType); ok {
					hierarchy = utils.GetParentHierarchy(crdNodeInfo.ParentHierarchy, item.GetLabels())
				}
			}

			link.Hierarchy = hierarchy
			link.Name = k
			link.Group = link.Group + "/v1"
			list[i] = link
			i++
		}
		return nc.JSON(http.StatusOK, list)
	}
	return nc.JSON(http.StatusInternalServerError, DefaultResponse{Message: "Something went wrong during link processing"})
}

// listHandler is used to process GET list requests
func listHandler(c echo.Context) error {
	nc := c.(*NexusContext)
	crdName := model.UriToCRDType[nc.NexusURI]
	crdInfo := model.CrdTypeToNodeInfo[crdName]

	labels := make(k8sLabels.Set)
	for k, v := range parseLabels(nc, crdInfo.ParentHierarchy) {
		labels[k] = v
	}

	// Setup GroupVersionResource
	parts := strings.Split(crdName, ".")
	gvr := schema.GroupVersionResource{
		Group:    strings.Join(parts[1:], "."),
		Version:  "v1",
		Resource: parts[0],
	}
	opts := metav1.ListOptions{
		LabelSelector: labels.AsSelector().String(),
	}

	if c.QueryParams().Has("limit") {
		i, err := strconv.ParseInt(c.QueryParams().Get("limit"), 10, 64)
		if err != nil {
			return err
		}
		opts.Limit = i
	} else {
		opts.Limit = 500
	}

	if c.QueryParams().Has("continue") {
		opts.Continue = c.QueryParams().Get("continue")
	}

	resps := make([]map[string]interface{}, 0)
	objs, err := client.Client.Resource(gvr).List(context.TODO(), opts)
	if err != nil {
		return handleClientError(nc, err)
	}
	for _, item := range objs.Items {
		itemName := item.GetName()
		if val, ok := item.GetLabels()[utils.DISPLAY_NAME_LABEL]; ok {
			itemName = val
		}
		status := make(map[string]interface{})
		if _, ok := item.Object["status"]; ok {
			status = item.Object["status"].(map[string]interface{})
		}
		delete(status, "nexus")
		spec := make(map[string]interface{})
		if _, ok := item.Object["spec"]; ok {
			spec = item.Object["spec"].(map[string]interface{})
		}
		for _, v := range crdInfo.Children {
			delete(spec, v.FieldNameGvk)
		}
		for _, v := range crdInfo.Links {
			delete(spec, v.FieldNameGvk)
		}

		r := make(map[string]interface{})
		r["name"] = itemName
		r["spec"] = spec
		r["status"] = status
		resps = append(resps, r)
	}

	return nc.JSON(http.StatusOK, resps)
}

// getNameFromParam gets name from param if exists
func getNameFromParam(nc *NexusContext, crdInfo model.NodeInfo) (string, string) {
	var name string
	for _, param := range nc.ParamNames() {
		if param == crdInfo.Name {
			name = nc.Param(param)
		}
	}

	if crdInfo.IsSingleton {
		if name == "" {
			name = nexus.DEFAULT_KEY
		}
		if name != nexus.DEFAULT_KEY {
			msg := fmt.Sprintf("Wrong singleton node name %s: %s for request %s, only '%s' is allowed as name",
				crdInfo.Name, name, nc.Request().RequestURI, nexus.DEFAULT_KEY)
			log.Debugf(msg)
			log.Debugf("crdName: %s, nexusURI: %s, paramNames: %s, paramValues: %s", crdInfo.Name, nc.NexusURI, nc.ParamNames(), nc.ParamValues())
			return "", msg
		}
	}

	// Get name from query params
	if val := nc.QueryParam(crdInfo.Name); val != "" {
		name = val
	}

	if name == "" {
		if val, ok := nc.Codes[http.StatusBadRequest]; ok {
			return "", val.Description
		} else {
			log.Debugf("Could not find required param %s for request %s", crdInfo.Name, nc.Request().RequestURI)
			log.Debugf("crdName: %s, nexusURI: %s, paramNames: %s, paramValues: %s", crdInfo.Name, nc.NexusURI, nc.ParamNames(), nc.ParamValues())
			return "", fmt.Sprintf("Could not find required param: %s", crdInfo.Name)
		}
	}

	return name, ""
}

// putHandler is used to process PUT requests
func putHandler(c echo.Context) error {
	nc := c.(*NexusContext)
	crdName := model.UriToCRDType[nc.NexusURI]
	crdInfo := model.CrdTypeToNodeInfo[crdName]

	// Get name from the URI segment
	name, msg := getNameFromParam(nc, crdInfo)
	if len(name) == 0 {
		return nc.JSON(http.StatusBadRequest, DefaultResponse{Message: msg})
	}

	// Parse body
	body := make(map[string]interface{})
	if err := (&echo.DefaultBinder{}).BindBody(nc, &body); err != nil {
		return err
	}

	// Setup GroupVersionResource
	parts := strings.Split(crdName, ".")
	gvr := schema.GroupVersionResource{
		Group:    strings.Join(parts[1:], "."),
		Version:  "v1",
		Resource: parts[0],
	}
	//package.Struct
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
			if uriInfo, ok := model.GetUriInfo(nc.NexusURI); ok && uriInfo.TypeOfURI == model.StatusURI {
				return c.JSON(http.StatusNotFound, DefaultResponse{Message: "Can't put status subresource as nexus object not found"})
			}

			// Build object
			err = client.CreateObject(gvr,
				crdNameParts[1], hashedName, labels, body)
			if err != nil {
				return handleClientError(nc, err)
			}

			// Update parent
			if len(crdInfo.ParentHierarchy) > 0 {
				parentCrdName := crdInfo.ParentHierarchy[len(crdInfo.ParentHierarchy)-1]
				parentCrd := model.CrdTypeToNodeInfo[parentCrdName]
				err = client.UpdateParentWithAddedChild(parentCrdName, parentCrd, labels, crdInfo, crdName, name, hashedName)
			}

			if err == nil {
				return c.JSON(http.StatusOK, DefaultResponse{Message: name})
			}
		}
		return handleClientError(nc, err)
	}

	obj.SetLabels(labels)
	return updateResource(nc, gvr, obj, body, crdInfo)
}

// patchHandler is used to modify specific fields
func patchHandler(c echo.Context) error {
	nc := c.(*NexusContext)
	crdName := model.UriToCRDType[nc.NexusURI]
	crdInfo := model.CrdTypeToNodeInfo[crdName]

	// Get name from the URI segment
	name, msg := getNameFromParam(nc, crdInfo)
	if len(name) == 0 {
		return nc.JSON(http.StatusBadRequest, DefaultResponse{Message: msg})
	}

	// Construct GroupVersionResource
	gvr := utils.ConstructGVR(crdName)

	// Mangle name
	hashedName := nexus.GetHashedName(crdName, crdInfo.ParentHierarchy, parseLabels(nc, crdInfo.ParentHierarchy), name)

	// Parse body
	body := make(map[string]interface{})
	if err := (&echo.DefaultBinder{}).BindBody(nc, &body); err != nil {
		return err
	}

	// Handle PATCH request for status subresource
	uriInfo, ok := model.GetUriInfo(nc.NexusURI)
	if ok && uriInfo.TypeOfURI == model.StatusURI {
		// Do not patch "nexus" status subresource; only user defined status subresource can be patched.
		delete(body, "nexus")

		// Prepare status patch payload
		statusPayload := struct {
			Status map[string]interface{} `json:"status"`
		}{
			body,
		}
		patchBytes, err := json.Marshal(statusPayload)
		if err != nil {
			return nc.JSON(http.StatusBadRequest, DefaultResponse{Message: fmt.Sprintf("error while marshaling status payload: %s", err.Error())})
		}
		log.Debugf("user defined status subresource PatchBytes %+v for CR %q", string(patchBytes), name)
		_, err = client.Client.Resource(gvr).Patch(context.TODO(), hashedName, types.MergePatchType, patchBytes, metav1.PatchOptions{}, "status")
		if err != nil {
			return handleClientError(nc, err)
		}
		return nc.JSON(http.StatusOK, DefaultResponse{Message: "Status patch applied successfully"})
	}

	// Handle PATCH request for object spec
	// Do not patch child/link object; omit child/link from the request body
	for _, v := range crdInfo.Children {
		delete(body, v.FieldNameGvk)
	}

	for _, v := range crdInfo.Links {
		delete(body, v.FieldNameGvk)
	}

	// Prepare patch payload
	payload := struct {
		Spec map[string]interface{} `json:"spec"`
	}{
		body,
	}

	patchBytes, err := json.Marshal(payload)
	if err != nil {
		return nc.JSON(http.StatusBadRequest, DefaultResponse{Message: fmt.Sprintf("error while marshaling payload: %s", err.Error())})
	}

	log.Debugf("PatchBytes %+v for CR %q", string(patchBytes), name)
	_, err = client.Client.Resource(gvr).Patch(context.TODO(), hashedName, types.MergePatchType, patchBytes, metav1.PatchOptions{})
	if err != nil {
		return handleClientError(nc, err)
	}

	return nc.JSON(http.StatusOK, DefaultResponse{Message: "Patch applied successfully"})
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

	if crdInfo.IsSingleton {
		if name == "" {
			name = nexus.DEFAULT_KEY
		}
		if name != nexus.DEFAULT_KEY {
			msg := fmt.Sprintf("Wrong singleton node name %s: %s for request %s, only '%s' is allowed as name",
				crdInfo.Name, name, nc.Request().RequestURI, nexus.DEFAULT_KEY)
			log.Debugf(msg)
			return nc.JSON(http.StatusBadRequest, DefaultResponse{Message: msg})
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
	err := client.DeleteObject(gvr, crdName, crdInfo, hashedName)
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
				labels[parent] = nexus.DEFAULT_KEY
			}
		}
	}

	return labels
}

func getUnstructuredObject(apiGroup, resourceName, name string) (*unstructured.Unstructured, error) {
	gvr := schema.GroupVersionResource{
		Group:    apiGroup,
		Version:  "v1",
		Resource: resourceName,
	}

	item, err := client.Client.Resource(gvr).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return item, nil
}

type PatchOp struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

// updateResource is used to create or update object resource or status subresource of nexus object
func updateResource(nc *NexusContext, gvr schema.GroupVersionResource, obj *unstructured.Unstructured, body map[string]interface{},
	crdInfo model.NodeInfo) error {
	uriInfo, ok := model.GetUriInfo(nc.NexusURI)
	if ok && uriInfo.TypeOfURI == model.StatusURI {
		if _, ok := body["nexus"]; ok {
			return nc.JSON(http.StatusBadRequest, DefaultResponse{Message: "can't update nexus status subresource, only user defined status subresource update is allowed"})
		}

		// Make sure status field is present first
		var err error
		if _, ok := obj.Object["status"]; !ok {
			m := []byte("{\"status\":{}}")
			_, err = client.Client.Resource(gvr).Patch(context.TODO(), obj.GetName(), types.MergePatchType, m, metav1.PatchOptions{}, "status")
		}
		if err != nil {
			return handleClientError(nc, err)
		}

		patch := createStatusPatch(body)
		var patchBytes []byte
		patchBytes, err = json.Marshal(patch)
		if err != nil {
			return nc.JSON(http.StatusBadRequest, DefaultResponse{Message: fmt.Sprintf("error while marshaling json status subresource payload: %s", err.Error())})
		}

		// Update status subresource
		_, err = client.Client.Resource(gvr).Patch(context.TODO(), obj.GetName(), types.JSONPatchType, patchBytes, metav1.PatchOptions{}, "status")
		if err != nil {
			return handleClientError(nc, err)
		}
		return nc.JSON(http.StatusOK, DefaultResponse{Message: "Status Updated successfully"})
	}

	if nc.QueryParam("update_if_exists") == "false" {
		return nc.JSON(http.StatusForbidden, DefaultResponse{Message: "Already Exists."})
	}

	spec := obj.Object["spec"].(map[string]interface{})
	for _, v := range crdInfo.Children {
		if value, ok := spec[v.FieldNameGvk]; ok {
			body[v.FieldNameGvk] = value
		}
	}
	for _, v := range crdInfo.Links {
		if value, ok := spec[v.FieldNameGvk]; ok {
			body[v.FieldNameGvk] = value
		}
	}
	obj.Object["spec"] = body

	_, err := client.Client.Resource(gvr).Update(context.TODO(), obj, metav1.UpdateOptions{})
	if err != nil {
		return handleClientError(nc, err)
	}
	return nc.JSON(http.StatusOK, DefaultResponse{Message: "Updated successfully"})
}

func createStatusPatch(body map[string]interface{}) []PatchOp {
	patch := []PatchOp{}
	for k, v := range body {
		p := PatchOp{
			Op:    "replace",
			Path:  "/status/" + k,
			Value: v,
		}
		patch = append(patch, p)
	}
	return patch
}

func DebugAllHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"isNexusRuntimeEnabled":                   config.Cfg.EnableNexusRuntime,
		"backendService":                          config.Cfg.BackendService,
		"crdTypeToRestUris":                       model.CrdTypeToRestUris,
		"uriToUriInfo":                            model.UriToUriInfo,
		"crdTypeToNodeInfo":                       model.CrdTypeToNodeInfo,
		"datamodelToDatamodelInfo":                model.DatamodelToDatamodelInfo,
		"declarativePaths":                        declarative.ApisList,
		"totalHttpServerRestarts":                 TotalHttpServerRestartCounter,
		"httpServerRestartsFromOpenApiSpecUpdate": HttpServerRestartFromOpenApiSpecUpdateCounter,
	})
}
