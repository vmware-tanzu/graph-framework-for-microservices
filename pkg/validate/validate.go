package validate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"k8s.io/client-go/dynamic"

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

var StopCh chan struct{}

type CRDStates struct {
	ParentsMap     sync.Map
	IsSingletonMap sync.Map
}

func (c *CRDStates) ProcessNewCRDType(crd v1.CustomResourceDefinition) error {
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

	c.ParentsMap.Store(crd.Name, annotation.Hierarchy)
	log.Infof("Added %s to parents map (%v)", crd.Name, annotation.Hierarchy)

	c.IsSingletonMap.Store(crd.Name, annotation.IsSingleton)
	log.Infof("Added %s to IsSingleton map (%v)", crd.Name, annotation.IsSingleton)

	return nil
}

func (c *CRDStates) GetApiGroups() (apiGroups []string) {
	c.ParentsMap.Range(func(key, value any) bool {
		parts := strings.Split(key.(string), ".")
		apiGroups = append(apiGroups, strings.Join(parts[1:], "."))
		return true
	})
	return
}

func (c *CRDStates) GetParents(crdName string, client dynamic.Interface) ([]string, error) {
	parents, ok := c.ParentsMap.Load(crdName)
	if !ok {
		log.Infof("Couldn't determine parents for %s, relisting CRDs to make sure CRD type definition wasn't"+
			" added in the meantime", crdName)
		ProcessCRDs(client)
		parents, ok = c.ParentsMap.Load(crdName)
		if !ok {
			return nil, fmt.Errorf("parents info not present for crd %s", crdName)
		}
	}
	copiedParents := make([]string, len(parents.([]string)))

	copy(copiedParents, parents.([]string))
	return copiedParents, nil
}

func (c *CRDStates) IsSingleton(crdName string) bool {
	isSingleton, ok := c.IsSingletonMap.Load(crdName)
	if ok {
		return isSingleton.(bool)
	}
	return false
}

var CRDs = CRDStates{
	ParentsMap:     sync.Map{},
	IsSingletonMap: sync.Map{},
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
