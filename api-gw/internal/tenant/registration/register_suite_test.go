package registration_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRegister(t *testing.T) {

	RegisterFailHandler(Fail)
	RunSpecs(t, "TenantRegistration Suite")

}
