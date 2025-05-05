package main

import (
	"fmt"

	"github.com/FelipeStillner/UdpSocket/lib/protocol"
)

func main() {
	request := protocol.Request{
		Path: "127.0.0.1:1234/hello",
	}
	fmt.Printf("\nRequest:\n\t%s\n", request.Path)

	client, err := protocol.NewClient()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	response, err := client.SendRequest(request)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Response %s:\n\t%s\n", protocol.TranslateStatus(response.Status), response.Body)
}
