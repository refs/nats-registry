package echoer

import (
	"encoding/json"
	"log"
	"natsreg/registry"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
)

// echoer is a service that echoes whatever it is in the URL

// Run starts an echoer service that registers itself.
func Run(r registry.Registry) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		nc, err := nats.Connect(nats.DefaultURL)
		if err != nil {
			panic(err)
		}

		if err := nc.Publish("register_service", []byte("{\"Name\": \"echoer-45495056-1f95-11ec-9f9f-3758f4a85c33\",\"Addr\": \"127.0.0.1:7766\"}")); err != nil {
			panic(err)
		}

		if err := nc.Publish("register_service", []byte("{\"Name\": \"echoer-2-45495056-1f95-11ec-9f9f-3758f4a85c33\",\"Addr\": \"127.0.0.1:7766\"}")); err != nil {
			panic(err)
		}
	}()

	go serve(r)
	<-stop
}

func serve(r registry.Registry) {
	http.HandleFunc("/list", listNodesHandler(r))
	http.HandleFunc("/register", registerNodeHandler(r))
	log.Fatal(http.ListenAndServe(":7766", nil))
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
			Name: name,
			Addr: address,
		}

		_, err := json.Marshal(n)
		if err != nil {
			println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
