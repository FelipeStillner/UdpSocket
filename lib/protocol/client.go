package protocol

import (
	"net"
	"sort"
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

	received_responses := []Response{}

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

		received_responses = append(received_responses, response)

		if response.quantity == len(received_responses) {
			wait = false
		}
	}

	sort.Slice(received_responses, func(i, j int) bool {
		return received_responses[i].number < received_responses[j].number
	})

	return_response.Status = received_responses[len(received_responses)-1].Status
	for _, response := range received_responses {
		return_response.Body = append(return_response.Body, response.Body...)
	}

	return return_response, nil
}
