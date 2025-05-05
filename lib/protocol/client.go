package protocol

import (
	"net"
)

type client struct {
	conn net.Conn
}

func NewClient() (*client, error) {
	conn, err := net.Dial("udp", "127.0.0.1:1234")
	if err != nil {
		return nil, err
	}
	return &client{
		conn: conn,
	}, nil
}

func (c *client) SendRequest(request Request) (Response, error) {
	c.conn.Write(request.Encode())

	wait := true

	return_response := Response{}

	received_responses := 0

	for wait {
		buffer := make([]byte, 1024)
		_, err := c.conn.Read(buffer)
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
