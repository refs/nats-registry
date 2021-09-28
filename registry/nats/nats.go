package nats

import (
	"encoding/json"
	"fmt"
	"log"
	"natsreg/registry"
	"sync"

	"github.com/nats-io/nats.go"
)

const (
	RegisterSubj = "register_service"
)

type natsReg struct {
	state *state
	conn  *nats.Conn
	wg    sync.WaitGroup
}

// GetConn returns a connection to the configured NATS registry.
func (r *natsReg) GetConn() *nats.Conn {
	return r.conn
}

// Register adds a node to the registry by sending a RegisterSubj message to the bus.
func (r *natsReg) Register(n registry.Node) error {
	log.Println(fmt.Sprintf("node registered: %+v", n))
	b, err := json.Marshal(n)
	if err != nil {
		return err
	}

	if err := r.conn.Publish(RegisterSubj, b); err != nil {
		return err
	}

	return nil
}

// ListNodes lists all the registered nodes on the registry.
func (r *natsReg) ListNodes() []*registry.Node {
	return r.state.GetNodes()
}

// New returns a configured nats registry.
// It will fail if the NATS server is incorrectly configured.
func New() (registry.Registry, error) {
	c, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	r := natsReg{
		state: &state{
			nodes: []*registry.Node{},
		},
		conn: c,
		wg:   sync.WaitGroup{},
	}

	return &r, nil
}

// Start a registry service. This operation blocks, use in a go routine.
// listen for services registration on the "register_service" subject.
func (r *natsReg) Start() error {
	r.wg.Add(1)
	_, err := r.conn.Subscribe(RegisterSubj, func(m *nats.Msg) {
		n, err := formatNode(m.Data)
		if err != nil {
			log.Printf("invalid message: [%s]", err.Error())
		}

		r.state.Append(n)
	})
	if err != nil {
		return err
	}

	r.wg.Wait()
	return nil
}

// Shutdown stops the registry service.
func (r *natsReg) Shutdown() {
	{
		r.wg.Done()
		r.conn.Close()
	}
	log.Println("shutting down NATS registry...")
}

// formatNode creates a node from a []byte.
func formatNode(b []byte) (*registry.Node, error) {
	node := registry.Node{}
	if err := json.Unmarshal(b, &node); err != nil {
		return nil, err
	}

	return &node, nil
}
