package registry

import "natsreg/service"

type Node struct {
	Name string
	Addr string
}

// Registry provides an abstraction to find services.
type Registry interface {
	// Subscribe a node in the registry
	Register(Node) error
	// ListNodes
	ListNodes() []*Node

	service.Service
}
