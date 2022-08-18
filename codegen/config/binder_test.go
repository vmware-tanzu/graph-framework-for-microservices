package config

import (
	"go/types"
	"testing"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/internal/code"

	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

func TestBindingToInvalid(t *testing.T) {
	binder, schema := createBinder(Config{})
	_, err := binder.TypeReference(schema.Query.Fields.ForName("messages").Type, &types.Basic{})
	require.EqualError(t, err, "Message has an invalid type")
}

func TestSlicePointerBinding(t *testing.T) {
	t.Run("without OmitSliceElementPointers", func(t *testing.T) {
		binder, schema := createBinder(Config{
			OmitSliceElementPointers: false,
		})

		ta, err := binder.TypeReference(schema.Query.Fields.ForName("messages").Type, nil)
		if err != nil {
			panic(err)
		}

		require.Equal(t, ta.GO.String(), "[]*gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/codegen/config/testdata/autobinding/chat.Message")
	})

	t.Run("with OmitSliceElementPointers", func(t *testing.T) {
		binder, schema := createBinder(Config{
			OmitSliceElementPointers: true,
		})

		ta, err := binder.TypeReference(schema.Query.Fields.ForName("messages").Type, nil)
		if err != nil {
			panic(err)
		}

		require.Equal(t, ta.GO.String(), "[]gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/codegen/config/testdata/autobinding/chat.Message")
	})
}

func createBinder(cfg Config) (*Binder, *ast.Schema) {
	cfg.Models = TypeMap{
		"Message": TypeMapEntry{
			Model: []string{"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/codegen/config/testdata/autobinding/chat.Message"},
		},
	}
	cfg.Packages = &code.Packages{}

	cfg.Schema = gqlparser.MustLoadSchema(&ast.Source{Name: "TestAutobinding.schema", Input: `
		type Message { id: ID }

		type Query {
			messages: [Message!]!
		}
	`})

	b := cfg.NewBinder()

	return b, cfg.Schema
}
