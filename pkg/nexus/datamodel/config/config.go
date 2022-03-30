package config

import (
	"helloworld/nexus"
)

type Config struct {
	nexus.Node
	ExampleStr string
}
