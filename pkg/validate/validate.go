package validate

import (
	"context"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

var CrdParentsMap = map[string][]string{}

type NexusAnnotation struct {
	Hierarchy []string `json:"hierarchy"`
}

func UpdateValidationWebhook(client *kubernetes.Clientset) {
	webhookConf, err := client.AdmissionregistrationV1().ValidatingWebhookConfigurations().Get(context.TODO(), "nexus-validation.webhook.svc", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	var apiGroups []string
	for k := range CrdParentsMap {
		parts := strings.Split(k, ".")
		apiGroups = append(apiGroups, strings.Join(parts[1:], "."))
	}

	for i := 0; i < len(webhookConf.Webhooks); i++ {
		webhook := &webhookConf.Webhooks[i]
		if webhook.Name == "nexus-validation-crd.webhook.svc" {
			for n := 0; i < len(webhook.Rules); i++ {
				rule := &webhook.Rules[n]
				rule.APIGroups = apiGroups
			}
		}
	}

	_, err = client.AdmissionregistrationV1().ValidatingWebhookConfigurations().Update(context.TODO(), webhookConf, metav1.UpdateOptions{})
	if err != nil {
		panic(err)
	}

	log.Info("Updated validating webhook configuration with new APIGroups: %v", apiGroups)
}
