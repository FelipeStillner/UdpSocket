package protocol

import (
	"fmt"
	"strconv"
	"strings"
)

type Request struct {
	Path string
}

type Response struct {
	Status   int
	quantity int
	number   int
	Body     []byte
}

func (m *Request) Encode() []byte {
	return []byte(fmt.Sprintf("%s", m.Path))
}

func (m *Request) Decode(data []byte) {
	m.Path = string(data)
}

func (m *Response) Encode() []byte {
	return []byte(fmt.Sprintf("%d\n%d\n%d\n%s", m.Status, m.quantity, m.number, string(m.Body)))
}

func (m *Response) Decode(data []byte) {
	parts := strings.Split(string(data), "\n")
	m.Status, _ = strconv.Atoi(parts[0])
	m.quantity, _ = strconv.Atoi(parts[1])
	m.number, _ = strconv.Atoi(parts[2])
	m.Body = []byte(strings.Join(parts[3:], "\n"))
}
