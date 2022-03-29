package ifc

type Notification struct {
	// in case of node operation it's full path of node
	// in case of link operation its the fullpath of the parent node.
	ObjectPath NodePathList
	// the name of the service causing the change
	Source string
	// time when the change was initiated.
	Timestamp string
	// type of update.
	UpdateType string // UpdateType;
	// in case of NodeAdd/Del/Update it will have Node ID other wise it will have (UNSUED) LinkID
	UpdatedObjId string
	// in case of NodeAdd/Del/Update it will have Node Type
	// UNUSED: in case of LinkAdd/Del/Update it will have link Type --> Change to DestNodeType
	UpdatedObjType string
	// in case of NodeAdd/Del/Update it will have Node key
	// in case of LinkAdd/Del/Updat it will have destination node key Value.
	UpdatedObjKey string
	// in case of NodeAdd/Del/Update it will not be present
	// in case of linkAdd/Del/Update it will have the parent node id.
	UpdatedObjParentId string
	// transaction ID
	TraceId string
	// adding additional fields for k/v type consistent DB like etcd
	Value    map[string]interface{}
	Revision int64
}
