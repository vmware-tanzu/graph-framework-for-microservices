package common

// Important: Please keep this enum in sync with graphDB.go:41 constructore
const (
	// indicates the key that is uniquly going to identify this node in a group
	//    node.properties["_NodeKeyName"] = name of the key i.e. nodeId
	//        that will mean node.properties["nodeId"] will need to be defined and will hold unique id for the node
	RootNode_Key              = "root"
	NodeFixedProp_NodeKeyName = "_NodeKeyName"
	// indicates if this is a root of the tree i.e. does not have a parent.
	NodeFixedProp_IsRoot = "_isRoot"
	// indicates what is the default key name to use in case of Singeltons or if app does not specify.
	NodeFixedProp_NodeDefaultKeyName = "_name"
	// the default value to use in case of singleton for node.properties["_name"]
	NodeFixedProp_NodeSingletonKeyValue = "default"
	NodeFixedProp_Revision              = "revision"
	NodeFixedProp_createdBy             = "createdBy"
	NodeFixedProp_updatedBy             = "updatedBy"
	NodeFixedProp_creationTime          = "creationTime"
	NodeFixedProp_updateTime            = "updateTime"
	NodeFixedProp_changeId              = "changeId"
	NodeFixedProp_ToBeDeleted           = "toBeDeleted"
)
