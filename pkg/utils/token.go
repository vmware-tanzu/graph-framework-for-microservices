package utils

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

const TokenFileName = "token"

type TokenRefresher struct {
	TokenCh                               chan token.Token
	RefreshDuration, RetryRefreshDuration time.Duration
}

func NewTokenRefresher() TokenRefresher {
	return TokenRefresher{
		TokenCh:              make(chan token.Token),
		RefreshDuration:      12 * time.Minute,
		RetryRefreshDuration: 5 * time.Second,
	}
}

// NewClientset is called when connector connects to destination via AWS roles.
func NewClientset(cluster *eks.Cluster) (dynamic.Interface, error) {
	ca, err := base64.StdEncoding.DecodeString(aws.StringValue(cluster.CertificateAuthority.Data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode cert of remote cluster %v: %v", *cluster.Name, err)
	}

	t, err := getToken(cluster.Name)
	if err != nil {
		return nil, fmt.Errorf("token fetch failed with an error: %v", err)
	}

	if err := writeAccessTokenToFile(t.Token); err != nil {
		return nil, fmt.Errorf("failed to write token to the file: %v", err)
	}
	clientset, err := dynamic.NewForConfig(
		&rest.Config{
			Host:            aws.StringValue(cluster.Endpoint),
			BearerTokenFile: TokenFileName,
			TLSClientConfig: rest.TLSClientConfig{
				CAData: ca,
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("could not generate k8s dynamic remote client: %v", err)
	}
	return clientset, nil
}

func RefreshToken(ctx context.Context, t TokenRefresher, name *string) {
	for {
		time.Sleep(t.RefreshDuration)
		for {
			select {
			case <-ctx.Done():
				log.Warnf("Stopping refreshing of token")
				close(t.TokenCh)
				return
			default:
				token, err := getToken(name)
				if err != nil {
					log.Errorf("Refresh failed with an error: %v, retrying in %v", err, t.RetryRefreshDuration)
					time.Sleep(t.RetryRefreshDuration)
					continue
				}
				t.TokenCh <- token
			}
			break
		}
		log.Infof("Access token refreshed successfully, refreshing again in %s", t.RefreshDuration)
	}
}

func writeAccessTokenToFile(token string) error {
	f, err := os.OpenFile(TokenFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("opening file %q for write failed with error: %v", TokenFileName, err)
	}
	defer f.Close()

	if _, err = f.WriteString(token); err != nil {
		return fmt.Errorf("writing token to file %s failed with error: %v", TokenFileName, err)
	}
	return nil
}

func getToken(name *string) (token.Token, error) {
	gen, err := token.NewGenerator(true, false)
	if err != nil {
		return token.Token{}, err
	}
	opts := &token.GetTokenOptions{
		ClusterID: aws.StringValue(name),
	}
	t, err := gen.GetWithOptions(opts)
	if err != nil {
		return token.Token{}, err
	}
	return t, nil
}
