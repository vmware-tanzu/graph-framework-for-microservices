package common

const (
	UpdateType_NodeAdd = iota
	UpdateType_NodeDelete
	UpdateType_NodeUpdate
	UpdateType_LinkAdd
	UpdateType_LinkDelete
	UpdateType_LinkUpdate
)
const (
	UpdateType_StrNodeAdd    = "NodeAdd"
	UpdateType_StrNodeDelete = "NodeDelete"
	UpdateType_StrNodeUpdate = "NodeUpdate"
	UpdateType_StrLinkAdd    = "LinkAdd"
	UpdateType_StrLinkDelete = "LinkDelete"
	UpdateType_StrLinkUpdate = "LinkUpdate"
	Rlink_EnableFeatureFlag  = "rlink-enable"
	DM_Debug_Server_Enable   = "dm-debug-server-enable"
)

func UpdateTypeConvFromString(ut string) int {
	switch ut {
	case UpdateType_StrNodeAdd:
		return UpdateType_NodeAdd
	case UpdateType_StrNodeDelete:
		return UpdateType_NodeDelete
	case UpdateType_StrNodeUpdate:
		return UpdateType_NodeUpdate
	case UpdateType_StrLinkAdd:
		return UpdateType_LinkAdd
	case UpdateType_StrLinkDelete:
		return UpdateType_LinkDelete
	default:
		return UpdateType_LinkUpdate
	}
}

type UpsertLinkOpNodeObj struct {
	NodeId, NodePath, NodeType, NodeKeyValue string
	NodeProp                                 map[string]interface{}
}

type UpsertLinkOpObj struct {
	LinkType    string
	SrcNodeObj  UpsertLinkOpNodeObj
	DestNodeObj UpsertLinkOpNodeObj
	LinkPropIn  map[string]interface{}
	IsSnglton   bool
}
