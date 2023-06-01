package declarative

import (
	"api-gw/pkg/config"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

type errorMessage struct {
	Message string `json:"message"`
}

func ApisHandler(c echo.Context) error {
	crdToSchemaMutex.Lock()
	defer crdToSchemaMutex.Unlock()

	crd := c.QueryParam("crd")
	if crd != "" {
		if val, ok := CrdToSchema[crd]; ok {
			return c.String(200, val)
		} else {
			return c.NoContent(http.StatusNotFound)
		}
	}

	return c.JSON(200, ApisList)
}

func ListHandler(c echo.Context) error {
	ec := c.(*EndpointContext)
	log.Debugf("ListHandler: %s <-> %s", c.Request().RequestURI, ec.SpecUri)

	url, err := BuildUrlFromParams(ec)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorMessage{Message: err.Error()})
	}

	log.Debugf("Making a request to: %s", url)
	resp, err := httpClient.Get(url)
	if err != nil {
		log.Warn(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	var respBody interface{}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		log.Warn(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(resp.StatusCode, respBody)
}

func GetHandler(c echo.Context) error {
	ec := c.(*EndpointContext)
	log.Debugf("GetHandler: %s <-> %s", c.Request().RequestURI, ec.SpecUri)

	url, err := BuildUrlFromParams(ec)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorMessage{Message: err.Error()})
	}

	log.Debugf("Making a request to: %s", url)
	resp, err := httpClient.Get(url)
	if err != nil {
		log.Warn(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	var respBody interface{}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		log.Warn(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(resp.StatusCode, respBody)
}

func PutHandler(c echo.Context) error {
	ec := c.(*EndpointContext)
	log.Debugf("PutHandler: %s <-> %s", c.Request().RequestURI, ec.SpecUri)

	body := make(map[string]interface{})
	if err := c.Bind(&body); err != nil {
		log.Warn(err)
		return c.JSON(http.StatusBadRequest, errorMessage{Message: "unable to parse body"})
	}

	if _, ok := body["metadata"]; !ok {
		return c.JSON(http.StatusBadRequest, errorMessage{Message: "metadata field not present"})
	}

	metadata := body["metadata"].(map[string]interface{})
	if _, ok := metadata["name"]; !ok {
		return c.JSON(http.StatusBadRequest, errorMessage{Message: "metadata.name field not present"})
	}

	url, err := BuildUrlFromBody(ec, metadata)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorMessage{Message: err.Error()})
	}

	// Build request
	req, _ := http.NewRequest(http.MethodPut, url, nil)
	if spec, ok := body["spec"]; ok {
		// Marshal spec from body
		jsonBody, err := json.Marshal(spec)
		if err != nil {
			log.Warn(err)
			return c.NoContent(http.StatusInternalServerError)
		}

		reqBody := bytes.NewBuffer(jsonBody)
		req, _ = http.NewRequest(http.MethodPut, url, reqBody)
		log.Debugf("Body: %s", reqBody.String())
	}
	req.Header.Set("Content-Type", "application/json")

	log.Debugf("Making a request to: %s", url)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Warn(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	var respBody interface{}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		log.Warn(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(resp.StatusCode, respBody)
}

func DeleteHandler(c echo.Context) error {
	ec := c.(*EndpointContext)
	log.Debugf("DeleteHandler: %s <-> %s", c.Request().RequestURI, ec.SpecUri)

	url, err := BuildUrlFromParams(ec)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorMessage{Message: err.Error()})
	}

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Warn(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	log.Debugf("Making a request to: %s", url)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Warn(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(resp.StatusCode)
}

func BuildUrlFromParams(ec *EndpointContext) (string, error) {
	url := config.Cfg.BackendService + ec.SpecUri
	labelSelector, _ := metav1.ParseToLabelSelector(ec.QueryParams().Get("labelSelector"))
	for _, param := range ec.Params {
		if param[1] == ec.Identifier {
			continue
		}

		labelVal := "default"
		if val, ok := labelSelector.MatchLabels[param[1]]; ok {
			labelVal = val
		}

		url = strings.Replace(url, param[0], labelVal, -1)
	}

	if ec.Single {
		url = strings.Replace(url, fmt.Sprintf("{%s}", ec.Identifier), ec.Param("name"), -1)
	}

	return url, nil
}

func BuildUrlFromBody(ec *EndpointContext, metadata map[string]interface{}) (string, error) {
	url := config.Cfg.BackendService + ec.SpecUri
	for _, param := range ec.Params {
		if param[1] == ec.Identifier {
			continue
		}

		labelVal := "default"

		if metadata["labels"] != nil {
			if val, ok := metadata["labels"].(map[string]interface{})[param[1]]; ok {
				labelVal = val.(string)
			}
		}
		url = strings.Replace(url, param[0], labelVal, -1)
	}
	url = strings.Replace(url, fmt.Sprintf("{%s}", ec.Identifier), metadata["name"].(string), -1)

	return url, nil
}
