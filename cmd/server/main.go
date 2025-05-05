package main

import (
	"fmt"

	"github.com/FelipeStillner/UdpSocket/lib/protocol"
)

func main() {
	server, err := protocol.NewServer()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	server.AddPath("hello")
	server.AddPath("unicode")
	server.AddPath("big")
	server.ListenRequests()
}
