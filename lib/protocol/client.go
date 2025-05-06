package protocol

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"slices"
	"sort"
	"strings"
	"time"
)

var (
	MAX_RETRIES = 5
	TIMEOUT     = 5 * time.Second
	LOSS_RATE   = 0
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
	request.Numbers = []int{}

	encoded_request, err := request.Encode()
	if err != nil {
		return Response{}, err
	}

	conn.Write(encoded_request)

	received_responses, err := receiveResponse(conn)
	if err != nil {
		return Response{}, err
	}

	shouldRetry, retries := verifyRetries(received_responses)
	try := 0
	for shouldRetry && try < MAX_RETRIES {
		request.Numbers = retries
		fmt.Printf("Retrying: %v\n", request.Numbers)
		encoded_request, err := request.Encode()
		if err != nil {
			return Response{}, err
		}
		conn.Write(encoded_request)
		received_reties_responses, err := receiveResponse(conn)
		if err != nil {
			return Response{}, err
		}
		received_responses = append(received_responses, received_reties_responses...)
		shouldRetry, retries = verifyRetries(received_responses)
		try++
	}

	return_response := joinResponses(received_responses)

	return return_response, nil
}

func joinResponses(received_responses []Response) Response {
	return_response := Response{}
	sort.Slice(received_responses, func(i, j int) bool {
		return received_responses[i].Number < received_responses[j].Number
	})

	return_response.Status = received_responses[len(received_responses)-1].Status
	for _, response := range received_responses {
		return_response.Body = append(return_response.Body, response.Body...)
	}

	return return_response
}

func receiveResponse(conn net.Conn) ([]Response, error) {
	responses_quantity := 0

	wait := true

	received_responses := []Response{}

	err := conn.SetReadDeadline(time.Now().Add(TIMEOUT))
	if err != nil {
		return []Response{}, err
	}

	for wait {
		buffer := make([]byte, 2048)
		_, err := conn.Read(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Printf("Timeout waiting for response\n")
				break
			}
			return []Response{}, err
		}

		response := Response{}
		err = response.Decode(buffer)
		if err != nil {
			return []Response{}, err
		}

		if rand.Intn(100) < LOSS_RATE {
			fmt.Printf("Simulating loss of response: %d\n", response.Number)
			continue
		}

		if response.Status != STATUS_OK {
			return []Response{response}, nil
		}

		fmt.Printf("Received OK response: %d\n", response.Number)

		received_responses = append(received_responses, response)

		responses_quantity++
		if response.Quantity == responses_quantity {
			wait = false
		}
	}

	return received_responses, nil
}

func verifyRetries(received_responses []Response) (bool, []int) {
	notRetry := []int{}
	quantity := 0
	for _, response := range received_responses {
		quantity = response.Quantity
		if response.Status != STATUS_OK {
			return false, []int{}
		}
		hash, err := response.getHash()
		if err != nil {
			return true, []int{}
		}
		if bytes.Equal(response.Hash, hash) {
			notRetry = append(notRetry, response.Number)
		}
	}
	if quantity == 0 {
		return true, []int{}
	}

	retries := []int{}
	for i := 0; i < quantity; i++ {
		if !slices.Contains(notRetry, i) {
			retries = append(retries, i)
		}
	}

	return len(retries) > 0, retries
}
