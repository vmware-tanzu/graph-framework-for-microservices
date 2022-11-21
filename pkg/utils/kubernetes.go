package utils

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// CustomHTTPRoundTripper adds additional headers to the request
type CustomHTTPRoundTripper struct {
	headers map[string][]string
	rt      http.RoundTripper
}

func SetUpDynamicRemoteAPI(host, token, cert string, eObj *NexusEndpoint) (dynamic.Interface, error) {
	if eObj != nil {
		if eObj.Cloud == "AWS" {
			sess := session.Must(session.NewSession(&aws.Config{
				Region: aws.String(eObj.ClientRegion),
			}))
			eksSvc := eks.New(sess)
			input := &eks.DescribeClusterInput{
				Name: aws.String(eObj.ClientName),
			}
			result, err := eksSvc.DescribeCluster(input)
			if err != nil {
				return nil, fmt.Errorf("error calling DescribeCluster: %v", err)
			}
			dynamicRemoteAPI, err := NewClientset(result.Cluster)
			if err != nil {
				return nil, fmt.Errorf("error creating clientset: %v", err)
			}
			// Start refreshToken routine.
			tokenRefresher := NewTokenRefresher()
			go RefreshToken(context.Background(), tokenRefresher, result.Cluster.Name)
			go func() {
				for token := range tokenRefresher.TokenCh {
					if err := writeAccessTokenToFile(token.Token); err != nil {
						return
					}
				}
			}()
			return dynamicRemoteAPI, nil
		}
	}

	rawDecodedText, err := base64.StdEncoding.DecodeString(cert)
	if err != nil {
		return nil, fmt.Errorf("could not decode cert: %v", err)
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(rawDecodedText)

	conf := &rest.Config{
		Host:        host,
		BearerToken: token,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: pool,
			},
		},
	}

	// If cert is not provided, then skip cert verification.
	if cert == "" {
		conf.Transport = nil
		conf.TLSClientConfig.Insecure = true
	}
	CreateCustomHTTPRoundTripper(conf)
	dynamicRemoteAPI, err := dynamic.NewForConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("could not generate k8s dynamic remote client: %v", err)
	}
	return dynamicRemoteAPI, nil
}

func SetUpLocalAPI() (kubernetes.Interface, error) {
	conf, err := GetRestConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get config: %v", err)
	}
	localAPI, err := kubernetes.NewForConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("could not generate k8s local client: %v", err)
	}
	return localAPI, nil
}

func SetUpDynamicLocalAPI() (dynamic.Interface, error) {
	conf, err := GetRestConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get config: %v", err)
	}
	dynamicAPI, err := dynamic.NewForConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("could not generate K8s dynamic local client for config: %v", err)
	}
	return dynamicAPI, nil
}

func GetRestConfig() (*rest.Config, error) {
	filePath := os.Getenv("KUBECONFIG")
	if filePath != "" {
		config, err := clientcmd.BuildConfigFromFlags("", filePath)
		if err != nil {
			return nil, err
		}
		return config, nil
	}
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func SetUpAPIs() (kubernetes.Interface, dynamic.Interface, error) {
	localAPI, err := SetUpLocalAPI()
	if err != nil {
		return nil, nil, err
	}
	localDynamicAPI, err := SetUpDynamicLocalAPI()
	if err != nil {
		return nil, nil, err
	}
	return localAPI, localDynamicAPI, nil
}

func getTokenFromSecret(dynamicLocalClient dynamic.Interface) (string, error) {
	secretsResource := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "secrets"}

	ns := os.Getenv(secretNS)
	name := os.Getenv(secretName)
	unstructuredSecret, err := dynamicLocalClient.Resource(secretsResource).Namespace(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("error getting secret object: %v", err)
	}

	secret := &corev1.Secret{}
	if err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredSecret.Object, secret); err != nil {
		return "", fmt.Errorf("error converting an unstructured object to secret: %v", err)
	}

	tokenInByte, ok := secret.Data["token"]
	if !ok {
		return "", fmt.Errorf("error looking for token field in secret data.: %v", secret)
	}
	return string(tokenInByte), nil
}

func BuildRemoteClientAPI(remoteEndpointHost, remoteEndpointPort, remoteEndpointCert string, dynamicLocalClient dynamic.Interface) (dynamic.Interface, error) {
	host := fmt.Sprintf("%s:%s", remoteEndpointHost, remoteEndpointPort)

	accessToken, err := getTokenFromSecret(dynamicLocalClient)
	if err != nil {
		return nil, fmt.Errorf("error getting token %v", err)
	}
	return SetUpDynamicRemoteAPI(host, accessToken, remoteEndpointCert, nil)
}

func (h *CustomHTTPRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, vv := range h.headers {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}
	return h.rt.RoundTrip(req)
}

func CreateCustomHTTPRoundTripper(conf *rest.Config) {
	wt := conf.WrapTransport
	conf.WrapTransport = func(rt http.RoundTripper) http.RoundTripper {
		if wt != nil {
			rt = wt(rt)
		}
		return &CustomHTTPRoundTripper{
			headers: map[string][]string{"Content-Type": {"application/json"}},
			rt:      rt,
		}
	}
}
