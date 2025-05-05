package protocol

import (
	"net"
	"strings"
)

type client struct {
}

func NewClient() (*client, error) {
	return &client{}, nil
}

func (c *client) SendRequest(request Request) (Response, error) {
	parts := strings.Split(request.Path, "/")

	conn, err := net.Dial("udp", parts[0])
	if err != nil {
		return Response{}, err
	}
	defer conn.Close()

	request.Path = strings.Join(parts[1:], "/")

	conn.Write(request.Encode())

	wait := true

	return_response := Response{}

	received_responses := 0

	for wait {
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			return Response{}, err
		}
		response := Response{}
		response.Decode(buffer)

		if response.Status != STATUS_OK {
			return response, nil
		}

		received_responses++

		return_response.Body = append(return_response.Body, response.Body...)
		return_response.quantity = response.quantity
		return_response.number = response.number
		return_response.Status = response.Status

		if response.quantity == received_responses {
			wait = false
		}
	}

	return return_response, nil
}
