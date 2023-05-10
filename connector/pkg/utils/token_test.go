package utils_test

import (
	"bytes"
	"context"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/service/eks"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"

	"connector/pkg/utils"
)

var _ = Describe("Token tests", func() {
	var (
		cluster   *eks.Cluster
		logBuffer bytes.Buffer
		t         utils.TokenRefresher
	)
	BeforeEach(func() {
		log.SetLevel(log.DebugLevel)
		log.SetOutput(&logBuffer)

		clusterName := "foo"
		clusterEndpoint := "https://foo.eks.amazonaws.com"
		cert := "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1ekNDQWMrZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJeU1Ea3hNekV3TVRjeU1Wb1hEVE15TURreE1ERXdNVGN5TVZvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTnZrCjlndlZXUWxFYjhFaCtqK2RDSnpnYXE5MzFaT1RWWGZ5cldZUm93M1JhOFFCZms4OXJMRGVKQXU3aWQrendyZkcKWjFBUTVqMXFPWXUrVE5MaWducU40dURoRzREVENYZXpoN0V3TzRJTFZHTGZpYUVvdlRkZ2xjVUZYS0MvQUVvWQpBeDk4Qm43UzVmUTE1blFiUUMxWk9SczU3VVZhYWphZVBlaUJPQm05d0ljdDJrY1hWZTZMQndLVHpFTXJ4UnZaCmFJZ05hd1dqWHhVMWxINUNUazd2b0ppZW1makRtTkpZQ0dIVTVwM2NNUDF1YmVnM3RXOEFSY0hBeDNPbXJPZmEKblRvbUdlaC9MWnhlNndUUUpoRGY1WmdGMFFVaVQ1citwOG9vaGRZVTllYmZINTYxNHBuV2dJcTkrMFBiN1FYeApqYW9wUjhCcjN3YnVDcCtuNmdzQ0F3RUFBYU5DTUVBd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0hRWURWUjBPQkJZRUZKbFZMeHdRQlN6STJFSlFLK3RacitvM1poNzdNQTBHQ1NxR1NJYjMKRFFFQkN3VUFBNElCQVFDWkVNZ09KSUlFSEUxSVlkRzUreHFKK25meW1HMEpMNFNVeEpxNTZBL1hzaVcyQ2prdwpnTWNxVkVFR2tBTTVjc2FHN2h5UXFpeUk2YVl2ajNvZWVROFpDMWdnQS9LR3hXUkZYT2pPaTNOaG0ySzFBQkxqCkxMbEs4UGJidmVnVFZsQitIc3AzdUl2SFk1V2tlcU5XWGg0UFZBa1gzVzBNSmlndnlXdHE5MzEwOU1kWGkzRkQKVEdwTWFwaG9CeGU4MC9pSHg2d213RVdLMmVtakVMS1hja0J5dHYrdkpMYjkvb0VqbnpmdXZldDIrSW5Wb0tnOQp5b0NPeUFabVpLdk5uWmRtM0t1aWk2M2Z0K25TaVltZytleWZ3NkttQ0tOMlpYVDJ4WWVPMk0xd2dlL1JrZ2lUCng5dzRqbWRXWFF6TDc1QkFVQ1VmV3JnbXlQMjdmNlJ2a3grRgotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="

		t = utils.NewTokenRefresher()
		cluster = &eks.Cluster{
			Name:     &clusterName,
			Endpoint: &clusterEndpoint,
			CertificateAuthority: &eks.Certificate{
				Data: &cert,
			},
		}
	})

	It("Should create remote-client with token successfully", func() {
		client, err := utils.NewClientset(cluster)
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())

		token, err := os.ReadFile(utils.TokenFileName)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(token)).ToNot(BeEmpty())
	})

	It("Should fail for invalid cert.", func() {
		cert := "foo"
		newCluster := cluster
		newCluster.CertificateAuthority.Data = &cert

		client, err := utils.NewClientset(newCluster)
		Expect(err).To(HaveOccurred())
		Expect(client).To(BeNil())
		Expect(err).To(MatchError("failed to decode cert of remote cluster foo: illegal base64 data at input byte 0"))
	})

	It("Should refresh token successfully", func() {
		t = utils.TokenRefresher{
			TokenCh:              make(chan token.Token),
			RefreshDuration:      3 * time.Second,
			RetryRefreshDuration: 1 * time.Second,
		}
		go utils.RefreshToken(context.TODO(), t, cluster.Name)
		Eventually(t.TokenCh, 5*time.Second).Should(Receive())
	})

	It("Should retry refreshing the token if fails", func() {
		t = utils.TokenRefresher{
			TokenCh:              make(chan token.Token),
			RefreshDuration:      3 * time.Second,
			RetryRefreshDuration: 1 * time.Second,
		}
		name := ""
		go utils.RefreshToken(context.TODO(), t, &name)
		Eventually(func() string { return logBuffer.String() }, 4*time.Second).Should(ContainSubstring("Refresh failed with an error: ClusterID is required, retrying in 1s"))
	})
})
