package main

import (
	"fmt"
	"net"

	"github.com/FelipeStillner/UdpSocket/lib/protocol"
)

func main() {
	p := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: 1234,
		IP:   net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	defer ser.Close()
	fmt.Printf("Running server on %v\n", ser.LocalAddr())
	for {
		_, remoteaddr, err := ser.ReadFromUDP(p)

		message := protocol.Request{}
		message.Decode(p)

		fmt.Printf("\nRequest from %v:\n\t%s\n", remoteaddr, message.Body)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}

		response := protocol.Response{
			Body: []byte("Hello, client!"),
		}

		fmt.Printf("Response to %v:\n\t%s\n", remoteaddr, response.Body)
		_, err = ser.WriteToUDP(response.Encode(), remoteaddr)
		if err != nil {
			fmt.Printf("Couldn't send response %v", err)
		}
	}
}
