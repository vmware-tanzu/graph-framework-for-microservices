package server

type Router interface {
	Start()
	RegisterRouter(urlPath string)
	RoutesNotification(stopCh chan struct{})
	StopServer()
}
