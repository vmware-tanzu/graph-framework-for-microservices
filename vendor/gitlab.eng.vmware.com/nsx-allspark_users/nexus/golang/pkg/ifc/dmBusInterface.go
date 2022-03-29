package ifc

type NotificationStats struct {
	NodeUpdateTx  int
	NodeUpdateRx  int
	LinkUpdateTx  int
	LinkUpdateRx  int
	RLinkUpdateTx int
	RLinkUpdateRx int
}

type DMBusInterface interface {
	Init(g interface{})
	Shutdown()
	RegisterCB(cbfn func(msg *Notification))
	GetStats() NotificationStats
	AddMessageDelay(n uint32)
}
