package protocol

import (
	"fmt"
	"math"
	"net"
	"os"
	"path/filepath"
)

type server struct {
	conn  net.UDPConn
	paths []string
}

func NewServer() (*server, error) {
	addr := net.UDPAddr{
		Port: 1234,
		IP:   net.ParseIP("127.0.0.1"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return nil, err
	}
	return &server{
		conn: *conn,
	}, nil
}

func (s *server) AddPath(path string) {
	s.paths = append(s.paths, path)
}

func (s *server) ListenRequests() {
	fmt.Printf("Running server on %v\n", s.conn.LocalAddr())
	buffer := make([]byte, 2048)
	for {
		n, remoteaddr, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		data := make([]byte, n)
		copy(data, buffer[:n])
		go s.handleRequest(data, remoteaddr)
	}
}

func (s *server) handleRequest(p []byte, remoteaddr *net.UDPAddr) {
	message := Request{}
	message.Decode(p)

	fmt.Printf("\nRequest from %v:\n\t%s\n", remoteaddr, message.Path)

	responses := []Response{}
	pathFound := false
	for _, path := range s.paths {
		if path == message.Path {
			pathFound = true
			break
		}
	}

	if pathFound {
		body, err := getBody(message.Path)
		if err != nil {
			responses = append(responses, Response{
				Status:   STATUS_INTERNAL_SERVER_ERROR,
				Body:     []byte("Error reading file"),
				quantity: 1,
				number:   0,
			})
		} else {
			size := len(body)
			quantity := int(math.Ceil(float64(size) / 1024))
			for i := 0; i < quantity; i++ {
				begin := i * 1024
				end := begin + 1024
				if end > size {
					end = size
				}
				responses = append(responses, Response{
					Status:   STATUS_OK,
					Body:     body[begin:end],
					quantity: quantity,
					number:   i,
				})
			}
		}
	} else {
		responses = append(responses, Response{
			Status:   STATUS_NOT_FOUND,
			Body:     []byte("Path not found"),
			quantity: 1,
			number:   0,
		})
	}

	for _, response := range responses {
		fmt.Printf("Response to %v:\n\t%s\n", remoteaddr, response.Body)
		s.conn.WriteToUDP(response.Encode(), remoteaddr)
	}
}

func getBody(path string) ([]byte, error) {
	filePath := filepath.Join("paths", path)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return content, nil
}
