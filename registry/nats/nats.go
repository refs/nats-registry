package nats

import (
	"encoding/json"
	"fmt"
	"natsreg/registry"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
)

type Registry struct {
	state *state
	conn  *nats.Conn
}

func (r *Registry) Conn() *nats.Conn {
	return r.conn
}

func (r *Registry) Subscribe() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		_, err := r.conn.Subscribe("register_service", func(m *nats.Msg) {
			node := registry.Node{}
			if err := json.Unmarshal(m.Data, &node); err != nil {
				panic(err)
			}

			r.state.nodes = append(r.state.nodes, node)

			fmt.Printf("Received a message: %+v\n", node)
		})
		if err != nil {
			panic(err)
		}
	}()

	<-stop
}

func (r *Registry) ListNodes() []registry.Node {
	return r.state.GetNodes()
}

// New returns a configured nats registry.
func New() registry.Registry {
	c, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	return &Registry{
		state: &state{
			nodes: []registry.Node{},
		},
		conn: c,
	}
}

// ugly. decouple state from this. we need more persistent storage anyway
type state struct {
	nodes []registry.Node
}

func (s *state) GetNodes() []registry.Node {
	return s.nodes
}
