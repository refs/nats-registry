package main

import (
	"natsreg/registry/nats"
	"natsreg/services/echoer"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	r, _ := nats.New()
	go r.Start()
	echoer.Run(r)

	<-stop
	r.Shutdown()
}
