type CInt uint8

type CrossPackageTester struct {
	Test gns.MyStr
}

type ReplicationSource struct {
	Kind SourceKind
}

type PolicyCfgAction struct {
	Action PolicyActionType `json:"action" mapstructure:"action"`
}

type ResourceGroupID struct {
	Name	string	`json:"name" mapstruction:"name"`
	Type	string	`json:"type" mapstruction:"type"`
}

type Config struct {
	GNS gns.Gns `nexus:"child"`
}

type EmptyStructTest struct{}

type queryFilters struct {
	StartTime		string
	EndTime			string
	Interval		string
	IsServiceDeployment	bool
	StartVal		int
}

type BArray []string

type TestValMarkers struct {
	//nexus-validation: MaxLength=8, MinLength=2, Pattern=ab
	MyStr	string	`json:"myStr" yaml:"myStr"`

	//nexus-validation: Maximum=8, Minimum=2
	//nexus-validation: ExclusiveMaximum=true
	MyInt	int	`json:"myInt" yaml:"myInt"`

	//nexus-validation: MaxItems=3, MinItems=2
	//nexus-validation: UniqueItems=true
	MySlice	[]string	`json:"mySlice" yaml:"mySlice"`
}

type RandomConst3 string

type MyConst string

type PolicyActionType string

type PolicyCfgActions []PolicyCfgAction

type StructWithEmbeddedField struct {
	SomeStruct
	gns.MyStr
}

type TempConst2 string

type RandomStatus struct {
	StatusX	int
	StatusY	int
}

type Host string

type Instance float32

type Answer struct {
	Name string
}

type ServiceSegmentRef struct {
	Field1	string
	Field2	string
}

type AMap map[string]string

type TempConst1 string

type TempConst3 string

type RandomConst1 string

type HostPort struct {
	Host	Host
	Port	Port
}

type AliasArr []int

type NonNexusType struct {
	Test int
}

type Cluster struct {
	Name	string
	MyID	int
}

type AdditionalStatus struct {
	StatusX	int
	StatusY	int
}

type ResourceGroupRef struct {
	Name	string
	Type	string
}

type MatchCondition struct {
	Name	string
	Type	gns.Host
}

type DFloat float32

type SomeStruct struct{}

type gnsQueryFilters struct {
	StartTime		string
	EndTime			string
	Interval		string
	IsServiceDeployment	bool
	StartVal		int
}

// This is Description struct.
type Description struct {
	Color		string
	Version		string
	ProjectId	string
	TestAns		[]Answer
	Instance	Instance
	HostPort	HostPort
}

type MyStr string

type metricsFilers struct {
	StartTime	string
	EndTime		string
	TimeInterval	string
	SomeUserArg1	string
	SomeUserArg2	int
	SomeUserArg3	bool
}

type GnsState struct {
	Working		bool
	Temperature	int
}

type ClusterNamespace struct {
	Cluster		MatchCondition
	Namespace	MatchCondition
}

type AdditionalDescription struct {
	DiscriptionA	string
	DiscriptionB	string
	DiscriptionC	string
	DiscriptionD	string
}

type RandomDescription struct {
	DiscriptionA	string
	DiscriptionB	string
	DiscriptionC	string
	DiscriptionD	string
}

type RandomConst2 string

type SourceKind string

type Port uint16

type ACPStatus struct {
	StatusABC	int
	StatusXYZ	int
}

type ACPSvcGroupLinkInfo struct {
	ServiceName	string
	ServiceType	string
}

type ResourceGroupIDs []ResourceGroupID

var nonNexusValue = 1

var nonValue int

const (
	Const3	TempConst3	= "Const3"
	Const2	TempConst2	= "Const2"
	Const1	TempConst1	= "Const1"
)

const (
	MyConst3	RandomConst3	= "Const3"
	MyConst2	RandomConst2	= "Const2"
	MyConst1	RandomConst1	= "Const1"
)

const (
	Object	SourceKind	= "Object"
	Type	SourceKind	= "Type"
	XYZ	MyConst		= "xyz"
)

const (
	PolicyActionType_Allow	PolicyActionType	= "ALLOW"
	PolicyActionType_Deny	PolicyActionType	= "DENY"
	PolicyActionType_Log	PolicyActionType	= "LOG"
	PolicyActionType_Mirror	PolicyActionType	= "MIRROR"
)

