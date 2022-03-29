package common

const (
	// indicates the key that is uniquly going to identify this node in a group
	// This information is copied form the linked node to the link properties
	//   link.properties["_NodeKeyName"] = name of the key in node property i.e. nodeId
	LinkFixedProp_NodeKeyName = "_NodeKeyName"
	// indicates what the unique id of the node that this link is attached to
	//   link.proeprty["_NodeKeyValue"] = unique id i.e. 123
	LinkFixedProp_NodeKeyValue = "_NodeKeyValue"
	// indicates the type of the node
	LinkFixedProp_NodeType = "_NodeType"
	// indicates if the link is a hard link.
	// when a node is deleted all the hard linked children are also deleted recursively
	LinkFixedProp_HardLink = "_HardLink"
	// if a soft link we need to know the destination node
	// path so we dont have to query the db many times.
	LinkFixedProp_SoftLinkDestinationPath  = "_SLDP"
	LinkFixedProp_RSoftLinkDestinationPath = "_RSLDP"
	LinkFixedProp_Revision                 = "revision"
	LinkFixedProp_createdBy                = "createdBy"
	LinkFixedProp_updatedBy                = "updatedBy"
	LinkFixedProp_creationTime             = "creationTime"
	LinkFixedProp_updateTime               = "updateTime"
	// introducing destination node id for etcd support.
	LinkFixedProp_destNodeId = "_DID"
	LinkType_Owner           = "OWNER"
	LinkType_ROwner          = "ROWNER"
	LinkType_Has             = "HAS"
)
