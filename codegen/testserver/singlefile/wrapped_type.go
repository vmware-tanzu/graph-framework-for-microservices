package singlefile

import "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/codegen/testserver/singlefile/otherpkg"

type (
	WrappedScalar = otherpkg.Scalar
	WrappedStruct otherpkg.Struct
	WrappedMap    otherpkg.Map
	WrappedSlice  otherpkg.Slice
)
