package registry

type Node struct {
	Name string
	Addr string
}

// Registry provides an abstraction to find services.
type Registry interface {
	// Subscribe a node in the registry
	Subscribe()
	// ListNodes
	ListNodes() []Node
}
