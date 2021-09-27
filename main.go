package main

import (
	"natsreg/registry/nats"
	"natsreg/services/echoer"
)

func main() {
	r := nats.New()
	go r.Subscribe()
	echoer.Run(r)
}
