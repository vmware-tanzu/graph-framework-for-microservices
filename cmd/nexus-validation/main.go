package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/validation.git/pkg/validate"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var client kubernetes.Interface
var dynamicClient dynamic.Interface

func init() {
	// Added this check for starting the deployment along with local apiserver config.
	kubeconfigFile := "/etc/kubeconfig/kubeconfig"
	var config *rest.Config
	_, err := os.Stat(kubeconfigFile)
	if errors.IsNotFound(err) {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigFile)
	}
	if err != nil {
		panic(err)
	}

	dynamicClient, err = dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	client, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
}

func main() {
	validate.ProcessCRDs(dynamicClient)
	validate.UpdateValidationWebhook(client)

	http.HandleFunc("/validate", ValidateHandler)
	http.HandleFunc("/validate-crd-type", ValidateCrdTypeHandler)
	http.HandleFunc("/healthz", healthz)

	cert := "/etc/nexus-validation/tls/tls.crt"
	key := "/etc/nexus-validation/tls/tls.key"
	log.Println("Started nexus-validation on port 443")
	log.Fatal(http.ListenAndServeTLS(":443", cert, key, nil))
}

func healthz(w http.ResponseWriter, r *http.Request) {

	_, err := w.Write([]byte("ok"))
	if err != nil {
		return
	}

}
func ValidateHandler(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	if len(body) == 0 {
		log.Warnln("empty body")
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	admReq := admissionv1.AdmissionReview{}
	if err := json.Unmarshal(body, &admReq); err != nil {
		log.Warnln("incorrect body")
		http.Error(w, "incorrect body", http.StatusBadRequest)
	}

	admRes, err := validate.Crd(dynamicClient, admReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(admRes)
	if err != nil {
		log.Warnf("could not encode response: %v", err)
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Warnf("could not write resopnse: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}

func ValidateCrdTypeHandler(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	if len(body) == 0 {
		log.Warnln("empty body")
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	admReq := admissionv1.AdmissionReview{}
	if err := json.Unmarshal(body, &admReq); err != nil {
		log.Warnln("incorrect body")
		http.Error(w, "incorrect body", http.StatusBadRequest)
	}

	admRes := validate.CrdType(client, admReq)
	resp, err := json.Marshal(admRes)
	if err != nil {
		log.Warnf("could not encode response: %v", err)
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Warnf("could not write resopnse: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}
