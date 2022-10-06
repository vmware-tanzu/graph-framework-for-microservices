package validate

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware-tanzu/graph-framework-for-microservices/src/kube-openapi/pkg/validation/spec"
	"github.com/vmware-tanzu/graph-framework-for-microservices/src/kube-openapi/pkg/validation/strfmt"
)

// valueValidator for string formats
func TestFormatValidator_EdgeCases(t *testing.T) {
	// Apply
	v := formatValidator{
		KnownFormats: strfmt.Default,
	}

	// formatValidator applies to: Items, Parameter,Schema

	s := spec.Schema{}
	s.Typed(stringType, "uuid")

	sources := []interface{}{&s}

	for _, source := range sources {
		// Default formats for strings
		assert.True(t, v.Applies(source, reflect.String))
		// Do not apply for number formats
		assert.False(t, v.Applies(source, reflect.Int))
	}

	assert.False(t, v.Applies("A string", reflect.String))
	assert.False(t, v.Applies(nil, reflect.String))
}
