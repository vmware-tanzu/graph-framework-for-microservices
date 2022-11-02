package client

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

var Client dynamic.Interface
var Host string

func New(config *rest.Config) (err error) {
	Host = config.Host
	Client, err = dynamic.NewForConfig(config)
	if err != nil {
		return err
	}
	return nil
}
