package followschema

import (
	"fmt"
	"io"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/graphql"
)

func MarshalBytes(b []byte) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = fmt.Fprintf(w, "%q", string(b))
	})
}

func UnmarshalBytes(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case string:
		return []byte(v), nil
	case *string:
		return []byte(*v), nil
	case []byte:
		return v, nil
	default:
		return nil, fmt.Errorf("%T is not []byte", v)
	}
}
