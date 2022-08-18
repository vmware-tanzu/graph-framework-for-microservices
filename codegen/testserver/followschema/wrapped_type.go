package followschema

import "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/codegen/testserver/followschema/otherpkg"

type (
	WrappedScalar = otherpkg.Scalar
	WrappedStruct otherpkg.Struct
	WrappedMap    otherpkg.Map
	WrappedSlice  otherpkg.Slice
)
