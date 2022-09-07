/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"time"

	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	log "github.com/sirupsen/logrus"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"authz-controller/pkg/utils"
	auth_nexus_org "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/authorization.nexus.org/v1"
)

// NexusRoleBindingReconciler reconciles a NexusRoleBinding object
type UserReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=authentication.nexus.org.authz-controller.com,resources=resourcerolebindings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=authentication.nexus.org.authz-controller.com,resources=resourcerolebindings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=authentication.nexus.org.authz-controller.com,resources=resourcerolebindings/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
const (
	CACertsPath            = "/etc/kubecerts"
	UserCertificateCRDName = "usercertificates.authorization.nexus.org"
)

func (r *UserReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var (
		nexusUser auth_nexus_org.User
		err       error
	)
	eventType := utils.Upsert
	if err = r.Get(ctx, req.NamespacedName, &nexusUser); err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Errorf("Error while trying to fetch User node with name %s", req.Name)
			return ctrl.Result{}, err
		}
		eventType = utils.Delete
	}

	userName := nexusUser.Labels[utils.DISPLAY_NAME_LABEL]
	log.Debugf("Creating Cert for user: %q", userName)

	certPrivateKey, cert, err := r.CreateUserCert(userName, CACertsPath)
	if err != nil {
		log.Errorf("Error while creating the certificate for user %s due to %s", userName, err)
		return ctrl.Result{}, err
	}

	certPrivateKeyString := base64.StdEncoding.EncodeToString(certPrivateKey.Bytes())

	certString := base64.StdEncoding.EncodeToString(cert.Bytes())
	certObjName := fmt.Sprintf("%s:%s", UserCertificateCRDName, userName)
	h := sha1.New()
	h.Write([]byte(certObjName))

	userCertificateName := hex.EncodeToString(h.Sum(nil))
	nexusUser.Labels["runtimes.runtime.nexus.org"] = "default"
	if err := r.CreateUserCertificate(
		ctx,
		certPrivateKeyString,
		certString,
		userCertificateName,
		nexusUser); err != nil {
		if err != nil {
			log.Errorf("Error while storing certificate for user %s", userName)
			return ctrl.Result{}, err
		}
	}

	log.Debugf("Received event %s for nexusUser node: Name %s", eventType, nexusUser.GetName())
	return ctrl.Result{}, nil
}

func (r *UserReconciler) CreateUserCertificate(ctx context.Context, certPrivateKeyString, certString, userCertificateName string,
	nexusUser auth_nexus_org.User) error {
	certCreateObject := auth_nexus_org.UserCertificate{
		Spec: auth_nexus_org.UserCertificateSpec{
			Key:  certPrivateKeyString,
			Cert: certString,
		},
		ObjectMeta: v1.ObjectMeta{
			Name:        userCertificateName,
			Labels:      nexusUser.Labels,
			Annotations: nexusUser.Annotations,
			OwnerReferences: []v1.OwnerReference{
				createOwnerReference(
					nexusUser.APIVersion,
					nexusUser.Kind,
					nexusUser.Name,
					nexusUser.GetUID(),
				),
			},
		},
	}
	return r.Client.Create(ctx, &certCreateObject, &client.CreateOptions{})
}

func (r *UserReconciler) CreateUserCert(name, path string) (*bytes.Buffer, *bytes.Buffer, error) {
	// convert ca crt file to cert
	caBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/ca.crt", path))
	if err != nil {
		log.Errorf("Error in reading the crt file: %s", err)
		return nil, nil, err
	}

	certObj, _ := pem.Decode(caBytes)
	if certObj == nil {
		return nil, nil, fmt.Errorf("could not load CA Cert file")
	}
	caCert, err := x509.ParseCertificate(certObj.Bytes)
	if err != nil {
		log.Errorf("Error in converting the crt file to cert: %s", err)
		return nil, nil, err
	}

	// convert key file to private key
	keyBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/ca.key", path))
	if err != nil {
		log.Errorf("Error in reading key file : %s", err)
		return nil, nil, err
	}
	keyObj, _ := pem.Decode(keyBytes)
	if keyObj == nil {
		return nil, nil, fmt.Errorf("could not load CA key file")
	}
	caPrivateKey, err := x509.ParsePKCS1PrivateKey(keyObj.Bytes)
	if err != nil {
		log.Errorf("Error in converting the key file to privateKey: %s", err)
		return nil, nil, err
	}

	cert := x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Country:       []string{"US"},
			Organization:  []string{"Default"},
			Province:      []string{""},
			Locality:      []string{""},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
			CommonName:    name,
		},
		NotBefore:             time.Now(),
		IsCA:                  false,
		BasicConstraintsValid: true,
		NotAfter:              time.Now().AddDate(365, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	clientPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Errorf("Error in generating client private key: %s", err)
		return nil, nil, err
	}
	clientCert, err := x509.CreateCertificate(rand.Reader, &cert, caCert, &clientPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		log.Errorf("Error in generating cert from client private key: %s", err)
		return nil, nil, err
	}

	certPrivateKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(clientPrivateKey),
	})

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: clientCert,
	})

	return certPrivateKeyPEM, certPEM, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&auth_nexus_org.User{}).
		Complete(r)
}
