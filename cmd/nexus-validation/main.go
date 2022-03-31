package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	datamodel "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/validation.git/pkg/nexus/generated/client/clientset/versioned"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/validation.git/pkg/validate"
	"io/ioutil"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/rest"
	"net/http"
)

var dmClient *datamodel.Clientset

func init() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	dmClient, err = datamodel.NewForConfig(config)
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/validate", ValidateHandler)

	cert := "/etc/nexus-validation/tls/tls.crt"
	key := "/etc/nexus-validation/tls/tls.key"
	log.Println("Started nexus-validation on port 443")
	log.Fatal(http.ListenAndServeTLS(":443", cert, key, nil))
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

	raw := admReq.Request.Object.Raw
	obj := struct {
		metav1.ObjectMeta `json:"metadata,omitempty"`
	}{}
	err := json.Unmarshal(raw, &obj)
	if err != nil {
		log.Warnf("could not unmarshal object meta: %v", err)
		http.Error(w, "could not unmarshal object meta", http.StatusInternalServerError)
		return
	}

	resourceName := fmt.Sprintf("%s.%s", admReq.Request.Resource.Resource, admReq.Request.Resource.Group)

	result, msg, err := validate.Validate(dmClient, resourceName, obj.ObjectMeta.GetLabels())
	if err != nil {
		log.Warnln("could not validate object")
		http.Error(w, "could not validate object", http.StatusInternalServerError)
		return
	}

	admRes := admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:     admReq.Request.UID,
			Allowed: result,
			Result: &metav1.Status{
				Message: msg,
			},
		},
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
