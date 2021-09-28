package echoer

import (
	"encoding/json"
	"fmt"
	"log"
	"natsreg/registry"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
)

// Run starts an echoer service that registers itself.
func Run(r registry.Registry) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	// https://stackoverflow.com/questions/43424787/how-to-use-next-available-port-in-http-listenandserve/43425461
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	log.Println("echoer listening on", listener.Addr().String())

	if err := r.Register(registry.Node{
		Name: fmt.Sprintf("%s-%s", "echoer", uuid.New().String()),
		Addr: listener.Addr().String(),
	}); err != nil {
		panic(err)
	}

	go serve(r, listener)
	<-stop
}

func serve(r registry.Registry, listener net.Listener) {
	http.HandleFunc("/list", listNodesHandler(r))
	http.HandleFunc("/register", registerNodeHandler(r))
	log.Fatal(http.Serve(listener, nil))
}

func listNodesHandler(reg registry.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := reg.ListNodes()
		b, err := json.Marshal(resp)
		if err != nil {
			println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		if _, err := w.Write(b); err != nil {
			println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}
}

// node info comes in the request body
// name: 	string
// address: string
func registerNodeHandler(reg registry.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		name := r.Form.Get("name")
		address := r.Form.Get("address")

		n := &registry.Node{
			Name: fmt.Sprintf("%s-%s", name, uuid.New().String()),
			Addr: address,
		}

		_, err := json.Marshal(n)
		if err != nil {
			println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		if err := reg.Register(*n); err != nil {
			log.Println("could not register node", n)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
