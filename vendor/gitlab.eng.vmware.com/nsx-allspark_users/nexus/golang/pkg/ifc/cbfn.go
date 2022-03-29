package ifc

type PropertyType map[string]interface{}
type CallbackFuncNode func(node BaseNodeInterface, ut int, od, nd PropertyType)
type CallbackFuncLink func(node BaseNodeInterface, ut int, nkey string, od, nd PropertyType)

type NodeCallbackNotification struct {
	Node       BaseNodeInterface
	UpdateType int
	OldData    PropertyType
	NewData    PropertyType
}
type LinkCallbackNotification struct {
	Node       BaseNodeInterface
	UpdateType int
	NodeKey    string
	OldData    PropertyType
	NewData    PropertyType
}
