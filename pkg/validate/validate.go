package validate

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"

	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	DEFAULT_KEY          = "default"
	DISPLAY_NAME_LABEL   = "nexus/display_name"
	IS_NAME_HASHED_LABEL = "nexus/is_name_hashed"
)

type CRDStates struct {
	mtx            sync.Mutex
	ParentsMap     map[string][]string
	IsSingletonMap map[string]bool
}

func (c *CRDStates) ProcessNewCRDType(crd v1.CustomResourceDefinition) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	nexusStr, ok := crd.Annotations["nexus"]
	if !ok {
		return nil
	}

	annotation := &NexusAnnotation{}
	err := json.Unmarshal([]byte(nexusStr), &annotation)
	if err != nil {
		log.Errorf("could not unmarshal nexus annotation: %v", err)
		return errors.New("could not unmarshal nexus annotation")
	}

	c.ParentsMap[crd.Name] = annotation.Hierarchy
	log.Infof("Added %s to parents map (%v)", crd.Name, annotation.Hierarchy)

	c.IsSingletonMap[crd.Name] = annotation.IsSingleton
	log.Infof("Added %s to IsSingleton map (%v)", crd.Name, annotation.IsSingleton)

	return nil
}

func (c *CRDStates) GetApiGroups() (apiGroups []string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	for k := range c.ParentsMap {
		parts := strings.Split(k, ".")
		apiGroups = append(apiGroups, strings.Join(parts[1:], "."))
	}

	return
}

func (c *CRDStates) GetParents(crdName string) []string {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	parents := c.ParentsMap[crdName]
	copiedParents := make([]string, len(parents))

	copy(copiedParents, parents)
	return copiedParents
}

func (c *CRDStates) IsSingleton(crdName string) bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.IsSingletonMap[crdName]
}

var CRDs = CRDStates{
	ParentsMap:     make(map[string][]string),
	IsSingletonMap: make(map[string]bool),
}

type NexusAnnotation struct {
	Name        string   `json:"name,omitempty"`
	Hierarchy   []string `json:"hierarchy,omitempty"`
	IsSingleton bool     `json:"is_singleton"`
}

func UpdateValidationWebhook(client kubernetes.Interface) {

	webhookConf, err := client.AdmissionregistrationV1().ValidatingWebhookConfigurations().Get(context.TODO(), "nexus-validation.webhook.svc", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	apiGroups := CRDs.GetApiGroups()

	for i := 0; i < len(webhookConf.Webhooks); i++ {
		webhook := &webhookConf.Webhooks[i]
		if webhook.Name == "nexus-validation-crd.webhook.svc" {
			for n := 0; i < len(webhook.Rules); i++ {
				rule := &webhook.Rules[n]
				if len(apiGroups) > 0 {
					rule.APIGroups = apiGroups
				}
			}
		}
	}
	if len(apiGroups) > 0 {
		_, err = client.AdmissionregistrationV1().ValidatingWebhookConfigurations().Update(context.TODO(), webhookConf, metav1.UpdateOptions{})
		if err != nil {
			panic(err)
		}
	}

	log.Infof("Updated validating webhook configuration with new APIGroups: %v", apiGroups)
}
