package main

import (
	"fmt"
	"net"

	"github.com/FelipeStillner/UdpSocket/lib/protocol"
)

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:1234")
	fmt.Printf("Running client on %v\n", conn.LocalAddr())
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	defer conn.Close()

	request := protocol.Request{
		Path: "",
		Body: []byte("Hello, server!"),
	}
	fmt.Printf("\nRequest to %v:\n\t%s\n", conn.RemoteAddr(), request.Body)
	conn.Write(request.Encode())

	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	response := protocol.Response{}
	response.Decode(buffer)
	fmt.Printf("Response from %v:\n\t%s\n", conn.RemoteAddr(), response.Body)
}
