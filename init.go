package main

import (
	"github.com/saxocellphone/multilift/server"
)

func main() {
	work := make(chan server.Request)
	go server.GetRequests(work)
	server.NewBalancer().Balance(work)
}
