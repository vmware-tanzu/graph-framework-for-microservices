package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/validation.git/pkg/validate"
	"io/ioutil"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var k8sClient dynamic.Interface

func init() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	k8sClient, err = dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
}

func main() {
	ProcessCRDs()

	http.HandleFunc("/validate", ValidateHandler)
	http.HandleFunc("/validate-crd-type", ValidateCrdTypeHandler)

	cert := "/etc/nexus-validation/tls/tls.crt"
	key := "/etc/nexus-validation/tls/tls.key"
	log.Println("Started nexus-validation on port 443")
	log.Fatal(http.ListenAndServeTLS(":443", cert, key, nil))
}

func ProcessCRDs() {
	gvr := schema.GroupVersionResource{
		Group:    "apiextensions.k8s.io",
		Version:  "v1",
		Resource: "customresourcedefinitions",
	}

	list, err := k8sClient.Resource(gvr).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, unstructuredCrd := range list.Items {
		var crd v1.CustomResourceDefinition
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredCrd.UnstructuredContent(), &crd)
		if err != nil {
			panic(err)
		}

		err = validate.ProcessCRDType(crd)
		if err != nil {
			panic(err)
		}
	}

	log.Infof("parents map: %v", validate.CrdParentsMap)
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

	admRes, err := validate.Crd(k8sClient, admReq)
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

	admRes := validate.CrdType(admReq)
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
