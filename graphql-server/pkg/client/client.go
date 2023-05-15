package client

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

var Client dynamic.Interface

func New(config *rest.Config) (err error) {
	Client, err = dynamic.NewForConfig(config)
	if err != nil {
		return err
	}
	return nil
}
