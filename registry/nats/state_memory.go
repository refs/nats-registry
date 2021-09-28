package nats

import "natsreg/registry"

// ugly. decouple state from this. we need more persistent storage anyway
type state struct {
	nodes []*registry.Node
}

func (s *state) GetNodes() []*registry.Node {
	return s.nodes
}

func (s *state) Append(n *registry.Node) {
	s.nodes = append(s.nodes, n)
}
